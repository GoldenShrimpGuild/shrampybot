package event

import (
	"shrampybot/router"
)

// Note: Auth is disabled for this route so individual endpoints
// must implement their own auth or be public!

func EventController(route *router.Route) *router.Response {
	resp := &router.Response{
		Body:       route.Router.ErrorBody(5, ""),
		StatusCode: "500",
		Headers:    &router.DefaultResponseHeaders,
	}
	if len(route.Path) < 2 {
		resp.Body = route.Router.ErrorBody(10, "")
		resp.StatusCode = "400"
		return resp
	}
	switch route.Path[1] {
	case "webhook":
		c := NewWebhookView()
		resp = c.CallMethod(route)
	}

	return resp
}
