package gsg

import (
	"shrampybot/router"
	"shrampybot/utility"
)

func GSGController(route *router.Route) *router.Response {
	resp := &router.Response{
		Body:       route.Router.ErrorBody(7),
		StatusCode: "403",
		Headers:    &router.DefaultResponseHeaders,
	}
	if len(route.Path) < 2 {
		resp.Body = route.Router.ErrorBody(10)
		resp.StatusCode = "400"
		return resp
	}

	// Scopes from JWT
	scopes := route.Router.Event.Scopes

	switch route.Path[1] {
	case "streamer":
		if utility.MatchScope(scopes, "gsg:streamer") {
			c := NewStreamerView()
			return c.CallMethod(route)
		}
	}

	return resp
}
