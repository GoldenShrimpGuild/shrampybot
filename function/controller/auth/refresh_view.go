package auth

import (
	"encoding/json"
	"log"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RefreshView struct {
	router.View `tstype:",extends,required"`
}

type RefreshRequestBody struct {
	RefreshToken string `json:"refresh"`
}

type RefreshResponseBody struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
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
	case "OPTIONS":
		return v.Options(route)
	}

	return router.NewResponse(router.GenericBodyDataFlat{}, "500")
}

func (v *RefreshView) Post(route *router.Route) *router.Response {
	log.Println("Entered route: Auth.Refresh.Post")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	reqBody := RefreshRequestBody{}
	json.Unmarshal([]byte(route.Body), &reqBody)

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		response.StatusCode = "500"
		return response
	}

	token := validateRefreshToken(reqBody.RefreshToken)
	if token == nil || !token.Valid {
		response.StatusCode = "401"
		return response
	}
	claims, res := token.Claims.(jwt.MapClaims)
	if !res {
		response.StatusCode = "500"
		return response
	}

	oAuth, err := n.GetOAuth(claims["kid"].(string))
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

	accessToken, err := generateAccessToken(oAuth)
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
		UserID:       oAuth.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	bodyBytes, _ := json.Marshal(body)
	response.Body = string(bodyBytes)

	response.StatusCode = "200"
	log.Println("Exiting route: Auth.Validate.Post")
	return response
}
