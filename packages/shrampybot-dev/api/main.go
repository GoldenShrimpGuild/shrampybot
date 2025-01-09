package main

import (
	"context"
	"shrampybot/router"

	"shrampybot/controller/event"
	"shrampybot/controller/user"
)

func Main(ctx context.Context, e router.Event) router.Response {
	router := router.NewRouter(&ctx, &e)
	router.AddRoute("event", event.EventController)
	router.AddRoute("user", user.UserController)

	return router.Route()
}
