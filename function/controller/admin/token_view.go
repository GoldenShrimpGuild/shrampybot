package admin

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"shrampybot/router"
	"shrampybot/utility"
	"shrampybot/utility/nosqldb"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TokenView struct {
	router.View `tstype:",extends,required"`
}

type NewTokenRequestBody struct {
	ExpiresAt time.Time `json:"expires_at"`
	Purpose   string    `json:"purpose"`
	Scopes    []string  `json:"scopes"`
}

type NewTokenResponseBody struct {
	TokenId string `json:"token_id"`
	Token   string `json:"token"`
}

func NewTokenView() *TokenView {
	c := TokenView{}
	return &c
}

func (v *TokenView) CallMethod(route *router.Route) *router.Response {
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

// // Get complete list of tokens or individual token info by ID
// // This will not return actual tokens as we aren't storing them server-side
// func (v *TokenView) Get(route *router.Route) *router.Response {
// 	var err error

// 	log.Println("Entered route: Admin.Token.Get")
// 	response := &router.Response{}
// 	response.Headers = &router.DefaultResponseHeaders

// 	response.Body = string(outBytes)
// 	response.StatusCode = "200"
// 	log.Println("Exited route: Admin.Token.Get")
// 	return response
// }

func (v *TokenView) Post(route *router.Route) *router.Response {
	var err error

	log.Println("Entered route: Admin.Token.Post")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	claims := route.Router.Event.Claims

	// Disallow static tokens from provisioning static tokens
	if claims["aud"] == "static" {
		log.Println("Cannot use a static token to provision as static token.")
		response.StatusCode = "403"
		return response
	}

	who := claims["sub"].(string)

	// Parse submitted category data
	requestBody := NewTokenRequestBody{}
	err = json.Unmarshal([]byte(route.Body), &requestBody)
	if err != nil {
		log.Printf("Could not unmarshal body json: %v\n", err)
		response.StatusCode = "500"
		return response
	}

	// Validate scopes and assemble new list
	validScopes := []string{}
	for _, scope := range requestBody.Scopes {
		if slices.Contains(utility.ValidStaticTokenScopes, scope) {
			validScopes = append(validScopes, scope)
		}
	}

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	// Prep database item
	static := nosqldb.StaticTokenDatum{}
	static.Id = uuid.NewString()
	static.CreatorId = who
	static.CreatedAt = time.Now()
	static.ExpiresAt = requestBody.ExpiresAt
	static.Purpose = requestBody.Purpose
	static.Revoked = false
	static.SecretKey = utility.GenerateRandomHex(sha256.BlockSize)
	static.Scopes = strings.Join(validScopes, " ")

	err = n.PutStaticToken(&static)
	if err != nil {
		log.Printf("Could not write token to table: %v\n", err)
		response.StatusCode = "500"
		return response
	}

	// Create JWT
	jwt, err := generateStaticToken(&static)
	if err != nil {
		log.Printf("Generate static token failed: %v\n", err)
		response.StatusCode = "500"
		return response
	}

	output := NewTokenResponseBody{
		TokenId: static.Id,
		Token:   jwt,
	}
	outBytes, err := json.Marshal(output)
	if err != nil {
		log.Printf("Could not marshal response output: %v\n", err)
		response.StatusCode = "500"
		return response
	}

	response.Body = string(outBytes)
	response.StatusCode = "200"
	log.Println("Exited route: Admin.Token.Post")
	return response
}
