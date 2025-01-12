package event

import (
	"encoding/json"
	"log"
	"shrampybot/connector/twitch"
	"shrampybot/router"
)

var (
	eventMap = map[string]func(sub *twitch.Subscription, event *map[string]string) error{
		"stream.online":  streamOnlineCallback,
		"stream.offline": streamOfflineCallback,
	}
)

type WebhookView struct {
	router.View
}

func NewWebhookView() *WebhookView {
	c := WebhookView{}
	return &c
}

func (v *WebhookView) CallMethod(route *router.Route) *router.Response {
	switch route.Method {
	case "GET":
		return v.Get(route)
	case "POST":
		return v.Post(route)
	case "PUT":
		return v.Put(route)
	case "PATCH":
		return v.Patch(route)
	case "DELETE":
		return v.Delete(route)
	}

	return router.NewResponse(router.GenericBody{}, "500")
}

// Handler for Twitch event webhooks.
func (c *WebhookView) Post(route *router.Route) *router.Response {
	response := router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	if !route.Router.Event.CheckTwitchAuthorization() {
		response.Body = route.Router.ErrorBody(14, "")
		response.StatusCode = "403"
	}

	log.Printf("Request body: %v\n", route.Router.Event.Body)

	// Prepare to return plain text as per Twitch's spec rather than json
	response.Body = ""
	response.StatusCode = "204"
	response.Headers = &router.ResponseHeaders{ContentType: "text/plain"}

	sub := &twitch.Subscription{}
	switch route.Router.Event.Headers.TwitchEventsubMessageType {

	case "webhook_callback_verification":
		requestBody := twitch.ChallengeWebhook{}
		json.Unmarshal([]byte(route.Router.Event.Body), &requestBody)
		// sub = requestBody.Subscription

		log.Println("Received webhook_callback_verification request.")
		response.Body = requestBody.Challenge
		response.StatusCode = "200"

	case "revocation":
		requestBody := twitch.RevocationWebhook{}
		json.Unmarshal([]byte(route.Router.Event.Body), &requestBody)
		sub = requestBody.Subscription

		log.Printf("Received revocation request: %v\n", sub)

	case "notification":
		requestBody := twitch.NotificationWebhook{}
		json.Unmarshal([]byte(route.Router.Event.Body), &requestBody)
		sub = requestBody.Subscription

		log.Printf("Received notification: %v\n", sub.Type)

		for subType, callback := range eventMap {
			if subType == sub.Type {
				callback(sub, requestBody.Event)
				break
			}
		}
	}

	return &response
}

func streamOnlineCallback(sub *twitch.Subscription, event *map[string]string) error {
	return nil
}

func streamOfflineCallback(sub *twitch.Subscription, event *map[string]string) error {
	return nil
}
