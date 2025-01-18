package main

import (
	"context"
	"shrampybot/controller/admin"
	"shrampybot/controller/auth"
	"shrampybot/controller/event"
	"shrampybot/controller/public"
	"shrampybot/router"

	"github.com/aws/aws-lambda-go/lambda"
)

func Main(ctx context.Context, ev router.Event) (router.AWSResponse, error) {
	router := router.NewRouter(&ctx, &ev)
	router.AddRoute("admin", admin.AdminController, true)
	router.AddRoute("auth", auth.AuthController, false)
	router.AddRoute("event", event.EventController, false)
	router.AddRoute("public", public.PublicController, false)

	routeResp := router.Route()
	return routeResp.FormatAWS(), nil
}

func main() {
	lambda.Start(Main)
}
