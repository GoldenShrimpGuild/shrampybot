package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"shrampybot/connector/discord"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RefreshView struct {
	router.View `tstype:",extends,required"`
}

// type RefreshRequestBody struct {
// 	RefreshToken string `json:"refresh"`
// }

type RefreshResponseBody struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access"`
}

func NewRefreshView() *RefreshView {
	c := RefreshView{}
	return &c
}

func (v *RefreshView) CallMethod(route *router.Route) *router.Response {
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

func (v *RefreshView) Post(route *router.Route) *router.Response {
	log.Println("Entered route: Auth.Refresh.Post")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	cookies, err := http.ParseCookie(route.Router.Event.Headers.Cookie)
	if err != nil {
		log.Printf("Issue parsing cookies from header: %v\n", err)
		response.StatusCode = "500"
		return response
	}
	var oldRefreshToken string
	for _, c := range cookies {
		if c.Name == "RefreshToken" {
			oldRefreshToken = c.Value
		}
	}
	if oldRefreshToken == "" {
		log.Println("No RefreshToken provided.")
		response.StatusCode = "500"
		return response
	}

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		response.StatusCode = "500"
		return response
	}

	token := validateRefreshToken(oldRefreshToken)
	if token == nil || !token.Valid {
		response.StatusCode = "401"
		return response
	}
	claims, res := token.Claims.(jwt.MapClaims)
	if !res {
		response.StatusCode = "500"
		return response
	}

	oAuth, err := n.GetOAuth(claims["sub"].(string))
	if err != nil {
		response.StatusCode = "500"
		return response
	}
	// Generate new refresh UUID and save
	oAuth.RefreshUID = uuid.NewString()
	err = n.PutOAuth(oAuth)
	if err != nil {
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
	scopes, err := dc.LocalScopesFromMembership(claims["sub"].(string))
	if err != nil {
		log.Printf("No scopes could be built for user %v: %v\n", claims["sub"].(string), err)
		response.StatusCode = "403"
		return response
	}

	accessToken, err := generateAccessToken(oAuth, scopes)
	if err != nil {
		log.Printf("Could not generate access token: %v\n", err)
		response.StatusCode = "500"
		return response
	}
	refreshToken, err := generateRefreshToken(oAuth)
	if err != nil {
		log.Printf("Could not generate refresh token: %v\n", err)
		response.StatusCode = "500"
		return response
	}

	// We already stored the RefreshUID but we won't be storing any detail
	// about the tokens themselves. Shit's going to be handled dynamically yo.

	body := RefreshResponseBody{
		UserID:      oAuth.Id,
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
