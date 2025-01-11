package admin

import (
	"shrampybot/connector/twitch"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
)

type Collection struct {
	router.View
}

func NewCollection() *Collection {
	c := Collection{}
	return &c
}

// Get should return summary stats for our collection of data.
func (c *Collection) Get(route *router.Route) *router.Response {
	body := router.GenericBody{
		Count: 0,
		Data:  []any{},
	}

	response := router.NewResponse(body, "200")
	return response
}

// A PATCH call will gather, assemble, and update all the user data required
// to do other Shrampy tasks. This is the linchpin of Shrampybot.
func (c *Collection) Patch(route *router.Route) *router.Response {
	body := router.GenericBody{
		Status: &router.Status{
			Msg: "Collection Patch!",
		},
		Count: 0,
		Data:  []any{},
	}

	// Fetch team member logins from Twitch
	th, err := twitch.NewClient()
	if err != nil {
		body.Status.Msg = "Failed to connect to Twitch."
		response := router.NewResponse(body, "500")
		return response
	}
	teamMembers, err := th.GetTeamMembers()
	if err != nil {
		body.Status.Msg = "Failed to retrieve team members from Twitch."
		response := router.NewResponse(body, "500")
		return response
	}
	loginList := []string{}
	for _, tm := range *teamMembers {
		loginList = append(loginList, tm.UserLogin)
	}

	// Fetch user records for logins on our list
	users, err := th.GetUsers(&loginList)
	if err != nil {
		body.Status.Msg = "Failed to retrieve users from Twitch."
		response := router.NewResponse(body, "500")
		return response
	}

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		body.Status.Msg = "Failed to connect to db."
		response := router.NewResponse(body, "500")
		return response
	}
	_ = n
	// TODO: complete working out batch write operations for DynamoDB

	// Update our user records
	// users, err := n.GetTwitchUsers()
	// if err != nil {
	// 	body.Status.Msg = "Could not retrieve data."
	// 	response := router.NewResponse(body, "500")
	// 	return response
	// }

	// Munge users into usable format
	for _, u := range *users {
		body.Data = append(body.Data, u)
	}
	body.Count = int64(len(*users))

	response := router.NewResponse(body, "200")
	return response
}

func (v *Collection) CallMethod(route *router.Route) *router.Response {
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

	return router.NewResponse(router.GenericBody{}, "500")
}
