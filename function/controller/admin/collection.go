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

// Get should return a summary of our data.
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
	var response *router.Response

	body := router.GenericBody{
		Count: 0,
		Data:  []any{},
	}

	users, err := getTwitchUsers()
	if err != nil || len(*users) == 0 {
		route.Router.ErrorBody(2, "")
		response = router.NewResponse(body, "500")
		return response
	}

	err = saveTwitchUsers(users)
	if err != nil {
		route.Router.ErrorBody(3, "")
		response = router.NewResponse(body, "500")
		return response
	}

	// Munge users into displayable format
	for _, u := range *users {
		body.Data = append(body.Data, u["login"])
	}
	body.Count = int64(len(*users))
	response = router.NewResponse(body, "200")
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

func getTwitchUsers() (*[]map[string]string, error) {
	var err error
	var users *[]map[string]string
	var loginList []string

	// connect to Twitch
	th, _ := twitch.NewClient()

	// Preparing to parallelize assembling the logins list
	// from multiple sources
	ch := make(chan string)
	go th.GetTeamMemberLoginsThreaded(ch)
	// TODO: Add threaded mastodon queries
	for login := range ch {
		loginList = append(loginList, login)
	}

	// Fetch user records for logins on our list
	users, err = th.GetUsers(&loginList)
	if err != nil {
		return users, err
	}

	return users, nil
}

func saveTwitchUsers(users *[]map[string]string) error {
	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		return err
	}

	err = n.PutTwitchUsers(users)
	if err != nil {
		return err
	}

	return nil
}
