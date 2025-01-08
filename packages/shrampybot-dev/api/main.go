package main

import (
	"context"
	"shrampybot/router"

	"shrampybot/controller/admin"
	"shrampybot/controller/event"
	"shrampybot/controller/user"
)

// func showEvent(event *router.Event) router.Response {
// 	data, _ := json.MarshalIndent(event, "", "  ")

// 	resp := router.Response{
// 		Body:       string(data),
// 		StatusCode: "200",
// 		Headers: router.ResponseHeaders{
// 			ContentType: "application/json",
// 		},
// 	}

// 	return resp
// }

func Main(ctx context.Context, e router.Event) router.Response {
	router := router.NewRouter(&ctx, &e)
	router.AddRoute("admin", admin.AdminController)
	router.AddRoute("event", event.EventController)
	router.AddRoute("user", user.UserController)

	return router.Route()
}
