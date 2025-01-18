package auth

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"shrampybot/connector/discord"
	"shrampybot/router"
	"shrampybot/utility"
	"shrampybot/utility/nosqldb"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

type ValidateView struct {
	router.View
}

type ValidateBody struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

func NewValidateView() *ValidateView {
	c := ValidateView{}
	return &c
}

func (v *ValidateView) CallMethod(route *router.Route) *router.Response {
	switch route.Method {
	case "GET":
		return v.Get(route)
	case "POST":
		return v.Post(route)
	case "PUT":
		return v.Put(route)
	case "PATCH":
		return v.Patch(route)
	case "DELETE":
		return v.Delete(route)
	}

	return router.NewResponse(router.GenericBodyDataFlat{}, "500")
}

func (v *ValidateView) Post(route *router.Route) *router.Response {
	log.Println("Entered route: Auth.Validate.Post")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		response.StatusCode = "500"
		return response
	}

	code := route.Query.Get("code")
	if code == "" {
		log.Println("No code provided in post query.")
		response.StatusCode = "403"
		return response
	}
	referer := route.Router.Event.Headers.Referer
	if referer == "" {
		log.Println("No referer header.")
		response.StatusCode = "403"
		return response
	}

	dOAuth, err := discordTokenExchange(code, referer)
	if err != nil {
		log.Println("Could not complete Discord token exchange.")
		response.StatusCode = "403"
		return response
	}

	d, err := discord.NewOAuthClient(dOAuth.AccessToken)
	if err != nil {
		log.Println("Could not create new Discord oauth client.")
		response.StatusCode = "403"
		return response
	}

	// Twofold purpose here:
	// 1. Test whether our newly generated credentials work
	// 2. Retrieve user ID / username for indexing our stored credentials
	user, err := d.GetSelf()
	if err != nil || user.ID == "" {
		log.Println("Could not retrieve Discord user with new OAuth credentials.")
		response.StatusCode = "500"
		return response
	}

	// Add ID and username fields to our Discord OAuth object and store to DB
	dOAuth.Id = user.ID
	dOAuth.Username = user.Username
	err = n.PutDiscordOAuth(dOAuth)
	if err != nil {
		log.Println("Could not store Discord OAuth record for %v to db.", dOAuth.Username)
		response.StatusCode = "500"
		return response
	}

	// Retrieve/produce signing key for shrampybot JWT for the user
	sbOAuth, err := n.GetOAuth(user.ID)
	if err != nil || sbOAuth.SecretKey == "" {
		sbOAuth.Id = user.ID
		sbOAuth.SecretKey = utility.GenerateRandomHex(sha256.BlockSize)
	}
	// Update refresh UID and store since we'll be generating new tokens
	sbOAuth.RefreshUID = uuid.NewString()
	err = n.PutOAuth(sbOAuth)
	if err != nil {
		log.Println("Could not generate ")
		response.StatusCode = "500"
		return response
	}

	// TODO: gen access token
	accessToken, err := generateAccessToken(sbOAuth)
	if err != nil {
		log.Println("Could not generate access token.")
		response.StatusCode = "500"
		return response
	}
	refreshToken, err := generateRefreshToken(sbOAuth)
	if err != nil {
		log.Println("Could not generate refresh token.")
		response.StatusCode = "500"
		return response
	}

	// We already stored the RefreshUID but we won't be storing any detail
	// about the tokens themselves. Shit's going to be handled dynamically yo.

	body := ValidateBody{
		UserID:       dOAuth.Id,
		Username:     dOAuth.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	bodyBytes, _ := json.Marshal(body)
	response.Body = string(bodyBytes)

	response.StatusCode = "200"
	log.Println("Exiting route: Auth.Validate.Post")
	return response
}

func discordTokenExchange(code string, redirect_url string) (*nosqldb.DiscordOAuthDatum, error) {
	oAuthResponse := nosqldb.DiscordOAuthDatum{}

	query_data := url.Values{}
	query_data.Set("code", code)
	query_data.Set("grant_type", "authorization_code")
	query_data.Set("redirect_uri", redirect_url)

	client := &http.Client{}
	request, err := http.NewRequest("POST", fmt.Sprintf("%voauth2/token", discordgo.EndpointAPI), strings.NewReader(query_data.Encode()))
	if err != nil {
		log.Println("Could not create new Discord token request.")
		return &oAuthResponse, err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("Unsuccessful request to Discord for new token: %v\n", err)
		return &oAuthResponse, err
	}
	if resp.StatusCode > 399 {
		log.Println("Error when requesting new Discord token.")
		return &oAuthResponse, nil
	}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &oAuthResponse)

	return &oAuthResponse, nil
}
