package admin

import (
	"shrampybot/router"
)

func AdminController(route *router.Route) router.Response {
	return router.Response{
		Body:       "{\"msg\": \"Called AdminController\"}",
		StatusCode: "200",
		Headers: router.ResponseHeaders{
			ContentType: "application/json",
		},
	}
}
