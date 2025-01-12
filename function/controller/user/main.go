package user

import (
	"encoding/json"
	"shrampybot/router"
)

func UserController(route *router.Route) *router.Response {
	bodyMap := map[string]any{
		"status": map[string]any{
			"msg": "Routed to UserController",
		},
		"count": 0,
		"data":  []map[string]any{},
	}

	body, _ := json.Marshal(bodyMap)
	return &router.Response{
		Body:       string(body),
		StatusCode: "200",
		Headers:    &router.DefaultResponseHeaders,
	}
}
