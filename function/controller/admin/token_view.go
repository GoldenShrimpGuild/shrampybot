package admin

import (
	"encoding/json"
	"log"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
	"time"
)

type TokenView struct {
	router.View `tstype:",extends,required"`
}

type NewTokenRequestBody struct {
	ExpiresAt time.Time `json:"expires_at"`
	Purpose   string    `json:"purpose"`
}

type TokenBody struct {
	router.GenericBodyDataFlat `tstype:",extends,required"`
	Data                       *[]nosqldb.StaticTokenDatum `json:"data" tstype:"nosqldb.StaticTokenDatum[]"`
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

func (v *TokenView) Post(route *router.Route) *router.Response {
	var err error

	log.Println("Entered route: Admin.Token.Get")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	claims := route.Router.Event.Claims

	// Disallow static tokens from provisioning static tokens
	if claims["aud"] == "static" {
		log.Println("Cannot use a static token to provision as static token.")
		response.StatusCode = "403"
		return response
	}

	// Parse submitted category data
	requestBody := NewTokenRequestBody{}
	err = json.Unmarshal([]byte(route.Body), &requestBody)
	if err != nil {
		log.Printf("Could not unmarshal body json: %v\n", err)
		response.StatusCode = "500"
		return response
	}

	// // Instantiate DynamoDB
	// n, err := nosqldb.NewClient()
	// if err != nil {
	// 	log.Println("Could not instantiate dynamodb.")
	// 	response.StatusCode = "500"
	// 	return response
	// }

	log.Println("Exited route: Admin.Token.Get")
	return response
}
