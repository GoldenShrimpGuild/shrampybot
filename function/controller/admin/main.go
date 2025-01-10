package admin

import (
	"shrampybot/router"
)

func AdminController(route *router.Route) router.Response {
	resp := router.Response{
		Body:       route.Router.ErrorBody(5, ""),
		StatusCode: "500",
		Headers:    &router.DefaultResponseHeaders,
	}
	if len(route.Path) < 2 {
		return resp
	}
	switch route.Path[1] {
	case "collection":
		c := NewCollection()
		return *c.CallMethod(route)
	}

	return resp
}
