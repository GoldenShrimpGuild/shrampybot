package auth

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"shrampybot/config"
	"shrampybot/connector/discord"
	"shrampybot/router"
	"shrampybot/utility"
	"shrampybot/utility/nosqldb"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ValidateView struct {
	router.View `tstype:",extends,required"`
}

type ValidateRequestBody struct {
	Code      string `json:"code"`
	GrantType string `json:"grant_type"`
}

type ValidateResponseBody struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	AccessToken string `json:"access"`
	// Commenting out RefreshToken as it will be handled with httponly cookies
	// RefreshToken string `json:"refresh"`
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
	case "OPTIONS":
		return v.Options(route)
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

	// code := route.Query.Get("code")
	// if code == "" {
	// 	log.Println("No code provided in post query.")
	// 	response.StatusCode = "403"
	// 	return response
	// }
	reqBody := ValidateRequestBody{}
	json.Unmarshal([]byte(route.Body), &reqBody)

	referer := route.Router.Event.Headers.Referer
	if referer == "" {
		log.Println("No referer header.")
		response.StatusCode = "403"
		return response
	}

	dOAuth, err := discordTokenExchange(reqBody.Code, referer)
	if err != nil {
		log.Println("Could not complete Discord token exchange.")
		response.StatusCode = "403"
		return response
	}

	d, err := discord.NewOAuthClient(dOAuth)
	if err != nil {
		log.Println("Could not create new Discord oauth client.")
		response.StatusCode = "403"
		return response
	}
	// No need to immediately save credentials in this instance, since we just created
	// fresh ones. Save can safely occur after a data fetch.

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
		log.Printf("Could not store Discord OAuth record for %v to db.\n", dOAuth.Username)
		response.StatusCode = "500"
		return response
	}

	err = mapDiscordConnections(user.ID, user.Username, n, d)
	if err != nil {
		log.Println("Could not map Discord connections.")
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
		log.Println("Could not store OAuth record.")
		response.StatusCode = "500"
		return response
	}

	// Connect to discord with GSG bot credentials
	dc, err := discord.NewBotClient()
	if err != nil {
		response.StatusCode = "500"
		return response
	}

	// Determine JWT scopes
	scopes, err := dc.LocalScopesFromMembership(user.ID)
	if err != nil {
		log.Printf("No scopes could be built for user %v: %v\n", user.ID, err)
		response.StatusCode = "403"
		return response
	}

	accessToken, err := generateAccessToken(sbOAuth, scopes)
	if err != nil {
		log.Printf("Could not generate access token: %v\n", err)
		response.StatusCode = "500"
		return response
	}
	refreshToken, err := generateRefreshToken(sbOAuth)
	if err != nil {
		log.Printf("Could not generate refresh token: %v\n", err)
		response.StatusCode = "500"
		return response
	}

	// We already stored the RefreshUID but we won't be storing any detail
	// about the tokens themselves. Shit's going to be handled dynamically yo.

	body := ValidateResponseBody{
		UserID:      sbOAuth.Id,
		AccessToken: accessToken,
	}
	bodyBytes, _ := json.Marshal(body)

	// RefreshToken in httponly cookie
	cookie := http.Cookie{
		Name:        "RefreshToken",
		Value:       refreshToken,
		HttpOnly:    true,
		SameSite:    http.SameSiteNoneMode,
		Secure:      true,
		Partitioned: true,
		Expires:     time.Now().Add(336 * time.Hour),
	}
	response.Headers.SetCookie = cookie.String()

	response.Body = string(bodyBytes)
	response.StatusCode = "200"
	log.Println("Exiting route: Auth.Validate.Post")
	return response
}

func discordTokenExchange(code string, redirect_base string) (*nosqldb.DiscordOAuthDatum, error) {
	oAuthResponse := nosqldb.DiscordOAuthDatum{}

	log.Printf("Code: %v\n", code)

	query_data := url.Values{}
	query_data.Set("grant_type", "authorization_code")
	query_data.Set("code", code)
	query_data.Set("redirect_uri", fmt.Sprintf("%vshrampybot/auth/validate_oauth", redirect_base))
	query_data.Set("client_id", config.DiscordClientId)
	query_data.Set("client_secret", config.DiscordClientSecret)

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, "https://discord.com/api/v10/oauth2/token", strings.NewReader(query_data.Encode()))
	if err != nil {
		log.Println("Could not create new Discord token request.")
		return &oAuthResponse, err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// request.SetBasicAuth(url.QueryEscape(string(client_id)), url.QueryEscape(client_secret))
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("Unsuccessful request to Discord for new token: %v\n", err)
		return &oAuthResponse, err
	}
	if resp.StatusCode > 399 {
		log.Printf("Error when requesting new Discord token: %v\n", resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("New discord token response body: %v\n", string(bodyBytes))
		return &oAuthResponse, fmt.Errorf("error when requesting new discord token")
	}
	responseMap := map[string]any{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &responseMap)

	oAuthResponse.TokenType = responseMap["token_type"].(string)
	oAuthResponse.AccessToken = responseMap["access_token"].(string)
	oAuthResponse.RefreshToken = responseMap["refresh_token"].(string)
	oAuthResponse.ExpiresIn = responseMap["expires_in"].(float64)
	oAuthResponse.ExpiresAt = time.Now().Add(time.Duration(oAuthResponse.ExpiresIn) * time.Second)
	oAuthResponse.Scope = responseMap["scope"].(string)

	return &oAuthResponse, nil
}
