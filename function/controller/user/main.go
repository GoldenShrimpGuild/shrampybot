package user

import (
	"shrampybot/router"
)

func UserController(route *router.Route) router.Response {
	return router.Response{
		Body: map[string]any{
			"status": map[string]any{
				"msg": "Routed to UserController",
			},
			"count": 0,
			"data":  []map[string]any{},
		},
		StatusCode: "200",
		Headers:    &router.DefaultResponseHeaders,
	}
}
