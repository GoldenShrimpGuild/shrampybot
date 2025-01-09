package event

import (
	"shrampybot/router"
)

func EventController(route *router.Route) router.Response {

	// Unhandled exception
	return router.Response{
		Body:       route.Router.ErrorBody(5, ""),
		StatusCode: "500",
		Headers:    &router.DefaultResponseHeaders,
	}
}
