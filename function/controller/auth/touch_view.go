package auth

import (
	"encoding/json"
	"log"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
	"strings"
)

type TouchView struct {
	router.View `tstype:",extends,required"`
}

type TouchResponseBody struct {
	UserId string `json:"user_id,omitempty"`
	Status string `json:"status,omitempty"`
}

func NewTouchView() *TouchView {
	c := TouchView{}
	return &c
}

func (v *TouchView) CallMethod(route *router.Route) *router.Response {
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

func (v *TouchView) Get(route *router.Route) *router.Response {
	log.Println("Entered route: Auth.Touch.Get")
	var err error
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders
	body := TouchResponseBody{}

	if !route.Router.Event.CheckAuthorizationJWT() {
		log.Println("Failed JWT Auth check.")
		body.Status = "expired"
		bodyBytes, _ := json.Marshal(body)
		response.Body = string(bodyBytes)
		response.StatusCode = "401"
		return response
	}
	// Get token object, defined when CheckingAuthorizationJWT above
	claims := route.Router.Event.Claims

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		response.StatusCode = "500"
		return response
	}

	oAuth, err := n.GetOAuth(claims["sub"].(string))
	if err != nil {
		response.StatusCode = "500"
		return response
	}
	body.UserId = claims["sub"].(string)

	if strings.Contains(oAuth.RefreshUID, "REVOKED") {
		response.StatusCode = "401"
		body.Status = "logged out"
		bodyBytes, _ := json.Marshal(body)
		response.Body = string(bodyBytes)
		return response
	}

	body.Status = "ok"
	bodyBytes, _ := json.Marshal(body)
	response.Body = string(bodyBytes)

	response.StatusCode = "200"
	log.Println("Exiting route: Auth.Touch.Get")
	return response
}
