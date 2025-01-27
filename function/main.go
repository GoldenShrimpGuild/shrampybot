package main

import (
	"context"
	"encoding/json"
	"shrampybot/controller/admin"
	"shrampybot/controller/auth"
	"shrampybot/controller/event"
	"shrampybot/controller/gsg"
	"shrampybot/controller/public"
	"shrampybot/router"

	"github.com/aws/aws-lambda-go/lambda"
)

func Main(ctx context.Context, ev map[string]any) (router.AWSResponse, error) {
	var evnt router.Event
	evBytes, _ := json.Marshal(ev)
	// Uncomment if there's a need to log headers
	// log.Println(string(evBytes))
	json.Unmarshal(evBytes, &evnt)

	router := router.NewRouter(&ctx, &evnt)
	router.AddRoute("admin", admin.AdminController, true)
	router.AddRoute("gsg", gsg.GSGController, true)

	// These don't necessarily lack auth, but they handle auth
	// themselves in various ways
	router.AddRoute("auth", auth.AuthController, false)
	router.AddRoute("event", event.EventController, false)
	router.AddRoute("public", public.PublicController, false)

	routeResp := router.Route()
	return routeResp.FormatAWS(), nil
}

func main() {
	lambda.Start(Main)
}
