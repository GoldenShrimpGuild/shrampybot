package public

import (
	"shrampybot/router"
)

func PublicController(route *router.Route) *router.Response {
	resp := &router.Response{
		Body:       route.Router.ErrorBody(5),
		StatusCode: "500",
		Headers:    &router.DefaultResponseHeaders,
	}
	if len(route.Path) < 2 {
		resp.Body = route.Router.ErrorBody(10)
		resp.StatusCode = "400"
		return resp
	}

	switch route.Path[1] {
	case "multi":
		c := NewMultiView()
		return c.CallMethod(route)
	case "stream":
		c := NewStreamView()
		return c.CallMethod(route)
	}

	return resp
}
