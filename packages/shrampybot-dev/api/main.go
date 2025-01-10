package main

import (
	"context"
	"encoding/json"
	"shrampybot/controller/event"
	"shrampybot/controller/user"
	"shrampybot/router"

	"github.com/aws/aws-lambda-go/lambda"
)

func Main(ctx context.Context, ev router.Event) (json.RawMessage, error) {
	// e := router.Event{}
	// err := json.Unmarshal(ev, &e)
	// if err != nil {
	// 	return json.RawMessage{}, err
	// }

	router := router.NewRouter(&ctx, &ev)
	router.AddRoute("event", event.EventController)
	router.AddRoute("user", user.UserController)

	routeResp := router.Route()
	bodyBytes, err := json.Marshal(routeResp.Body)
	if err != nil {
		return json.RawMessage{}, err
	}

	output := map[string]any{
		"body":            string(bodyBytes),
		"statusCode":      routeResp.StatusCode,
		"headers":         routeResp.Headers,
		"isBase64Encoded": "false",
	}

	response, err := json.Marshal(output)
	if err != nil {
		return json.RawMessage{}, err
	}
	return response, nil
}

func main() {
	lambda.Start(Main)
}
