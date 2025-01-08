package user

import (
	"shrampybot/router"
)

func UserController(route *router.Route) router.Response {
	return router.Response{
		Body:       "{\"msg\": \"Called UserController\"}",
		StatusCode: "200",
		Headers: router.ResponseHeaders{
			ContentType: "application/json",
		},
	}
}
