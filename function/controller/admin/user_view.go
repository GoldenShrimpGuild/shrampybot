package admin

import (
	"encoding/json"
	"log"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
)

type UserView struct {
	router.View `tstype:",extends,required"`
}

func NewUserView() *UserView {
	c := UserView{}
	return &c
}

func (v *UserView) CallMethod(route *router.Route) *router.Response {
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

func (c *UserView) Get(route *router.Route) *router.Response {
	log.Println("Entered route: Admin.User.Get")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders
	var logins []*nosqldb.TwitchUserDatum

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	// Fetch login names from our stored Twitch users
	logins, err = n.GetTwitchUsers()
	if err != nil {
		log.Println("Could not get saved Twitch logins.")
		response.StatusCode = "500"
		return response
	}

	body := map[string]any{}
	body["count"] = len(logins)
	body["data"] = logins
	bodyBytes, _ := json.Marshal(body)

	response.StatusCode = "200"
	response.Body = string(bodyBytes)

	log.Println("Exited route: Admin.User.Get")
	return response
}
