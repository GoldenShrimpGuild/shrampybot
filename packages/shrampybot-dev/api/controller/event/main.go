package event

import (
	"shrampybot/router"
)

func EventController(route *router.Route) router.Response {
	return router.Response{
		Body:       "{\"msg\": \"Called EventController\"}",
		StatusCode: "200",
		Headers: router.ResponseHeaders{
			ContentType: "application/json",
		},
	}
}
