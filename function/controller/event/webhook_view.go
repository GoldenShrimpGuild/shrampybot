package event

import (
	"encoding/json"
	"errors"
	"log"
	"shrampybot/connector/twitch"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
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

	return router.NewResponse(router.GenericBodyDataFlat{}, "500")
}

// Handler for Twitch event webhooks.
func (c *WebhookView) Post(route *router.Route) *router.Response {
	response := router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	if !route.Router.Event.CheckTwitchAuthorization() {
		response.Body = route.Router.ErrorBody(14)
		response.StatusCode = "403"
		return &response
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
		json.Unmarshal([]byte(route.Body), &requestBody)
		// sub = requestBody.Subscription

		log.Println("Received webhook_callback_verification request.")
		response.Body = requestBody.Challenge
		response.StatusCode = "200"

		// TODO: Record

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
	log.Println("Entered streamOnlineCallback")

	if (*event)["type"] != "live" {
		return errors.New("event stream type is not 'live'")
	}

	userId := (*event)["broadcaster_user_id"]

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		return err
	}

	// Get user from twitch_users
	user, err := n.GetTwitchUser(userId)
	if err != nil {
		log.Printf("Could not find user record for %v\n", userId)
	}
	_ = user

	// TODO: Check for duplicate event/notification IDs.

	// TODO: Debounce duplicate event IDs

	// TODO: Check if stream category is in category_map

	// TODO: Get live stream information from Twitch

	// TODO: Debounce duplicate stream IDs

	// TODO: Record stream information to stream_history

	// TODO: Prepare preview image and other needed data

	// TODO: Use goroutines to post to 3 platforms simultaneously

	// TODO: Update eventsub_notice_history with platform post data

	return nil
}

func streamOfflineCallback(sub *twitch.Subscription, event *map[string]string) error {
	log.Println("Entered streamOnlineCallback")

	// TODO: Log stream offline time in eventsub_notice_history

	return nil
}
