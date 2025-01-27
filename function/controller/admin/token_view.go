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

var (
	validStaticTokenScopes = []string{
		"login",
		"dev",
		"admin",
		"admin:users",
		"admin:categories",
		"admin:tokens",
	}
)

type TokenView struct {
	router.View `tstype:",extends,required"`
}

type NewTokenRequestBody struct {
	ExpiresAt time.Time `json:"expires_at"`
	Purpose   string    `json:"purpose"`
	Scopes    []string  `json:"scopes"`
}

type OutputStaticTokenInfo struct {
	Id        string    `json:"id"`
	CreatorId string    `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
	Scopes    string    `json:"scopes,omitempty"`
	Purpose   string    `json:"purpose"`
}

type NewTokenResponseBody struct {
	OutputStaticTokenInfo `tstype:",extends,required"`
	Token                 string `json:"token,omitempty"`
}

type ExtTokenResponseBody struct {
	Count int                      `json:"count"`
	Data  []*OutputStaticTokenInfo `json:"data"`
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
func (v *TokenView) Get(route *router.Route) *router.Response {
	var err error

	log.Println("Entered route: Admin.Token.Get")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	tokens, err := n.GetStaticTokensNoDecrypt()
	if err != nil {
		log.Println("Could not retrieve tokens from db.")
		response.StatusCode = "500"
		return response
	}

	tokensBytes, _ := json.Marshal(tokens)
	respBody := ExtTokenResponseBody{}
	json.Unmarshal(tokensBytes, &respBody.Data)
	respBody.Count = len(respBody.Data)

	outBytes, _ := json.Marshal(respBody)

	response.Body = string(outBytes)
	response.StatusCode = "200"
	log.Println("Exited route: Admin.Token.Get")
	return response
}

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
		if slices.Contains(validStaticTokenScopes, scope) {
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

	output := NewTokenResponseBody{}

	staticBytes, _ := json.Marshal(static)
	json.Unmarshal(staticBytes, &output)
	output.Token = jwt

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

func (v *TokenView) Delete(route *router.Route) *router.Response {
	var err error
	log.Println("Entered route: Admin.Token.Delete")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	if len(route.Path) == 3 {
		token, err := n.GetStaticToken(route.Path[2])
		if err != nil {
			log.Println("Retrieve token failed.")
			response.StatusCode = "500"
			return response
		}

		log.Printf("Revoking static token for ID: %v\n", route.Path[2])
		token.Revoked = true

		err = n.PutStaticToken(token)
		if err != nil {
			log.Println("Save token failed.")
			response.StatusCode = "500"
			return response
		}

	} else {
		log.Println("No ID specified.")
		response.StatusCode = "400"
		return response
	}

	response.StatusCode = "200"
	log.Println("Exited route: Admin.Token.Delete")
	return response
}
