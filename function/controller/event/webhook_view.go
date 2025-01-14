package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"shrampybot/config"
	"shrampybot/connector/bluesky"
	"shrampybot/connector/discord"
	"shrampybot/connector/mastodon"
	"shrampybot/connector/twitch"
	"shrampybot/router"
	"shrampybot/utility"
	"shrampybot/utility/nosqldb"
	"strconv"
	"strings"

	"github.com/litui/helix/v3"
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

// Handler for Twitch event webhooks which all come in as POST
func (c *WebhookView) Post(route *router.Route) *router.Response {
	response := router.Response{}

	// Flag to determine if event processing should happen.
	// Used with duplicate checking.
	doNotProcess := false

	// Prepare to return plain text as per Twitch's spec rather than json
	response.Headers = &router.ResponseHeaders{ContentType: "text/plain"}

	if !route.Router.Event.CheckTwitchAuthorization() {
		response.Body = "Authentication failed."
		response.StatusCode = "403"
		return &response
	}

	log.Printf("Request body: %v\n", route.Router.Event.Body)

	response.Body = ""
	response.StatusCode = "204"

	log.Println("Checking for duplicate Twitch message ID.")
	// Before going further, check if we've received this message before
	if messageIsDuplicate(route.Router.Event.Headers.TwitchEventsubMessageId) {
		doNotProcess = true
	} else {
		// Record new eventsub message for duplicate checking
		n, _ := nosqldb.NewClient()
		n.PutEventsubMessage(&nosqldb.EventsubMessageDatum{
			Id:    route.Router.Event.Headers.TwitchEventsubMessageId,
			Time:  route.Router.Event.Headers.TwitchEventsubMessageTimestamp,
			Type:  route.Router.Event.Headers.TwitchEventsubMessageType,
			Retry: route.Router.Event.Headers.TwitchEventsubMessageRetry,
		})
	}
	// Continue running so our responses align, but doNotProcess should
	// bypass any further logic.

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

		if !doNotProcess {
			log.Println("Processing event notification.")
			for subType, callback := range eventMap {
				if subType == sub.Type {
					callback(sub, requestBody.Event)
					break
				}
			}
		} else {
			log.Println("Not processing notification due to duplicate notice.")
		}
	}

	return &response
}

func streamOnlineCallback(sub *twitch.Subscription, eventMap *map[string]string) error {
	log.Println("Entered streamOnlineCallback")

	// Unmarshal event data into helix struct
	event := helix.EventSubStreamOnlineEvent{}
	evBytes, _ := json.Marshal(eventMap)
	json.Unmarshal(evBytes, &event)

	// Connect to the systems we'll need for lookups/storage
	n, _ := nosqldb.NewClient()
	t, err := twitch.NewClient()
	if err != nil {
		log.Println("Could not connect to Twitch API. Can't continue.")
		return err
	}

	if event.Type != "live" {
		return errors.New("event stream type is not 'live'")
	}
	userId := event.BroadcasterUserID

	// Get user from twitch_users table
	user, err := n.GetTwitchUser(userId)
	if err != nil {
		log.Printf("Could not find user record for %v\n", userId)
		return err
	}
	log.Printf("UserID %v matches user %v\n", userId, user.Login)

	// Contact Twitch for actual stream info
	// (Done in this order because the twitch cli mock client sucks at setting event ID)
	log.Printf("Fetching stream from twitch for user %v\n", user.Login)
	tStream, err := t.GetStreamByUserId(userId)
	if err != nil {
		log.Printf("Error fetching stream from twitch.")
		return err
	}

	// Lookup stream in our history using Twitch stream ID
	stream, err := n.GetStream(tStream.ID)
	if err != nil {
		log.Printf("Error getting history record for stream ID %v\n", event.ID)
	}
	if stream == nil || stream.ID == "" {
		// If no such stream in our history, grab it from Twitch
		log.Println("No stream found in our history, loading from Twitch.")
		// Convert to our stream history type
		tsBytes, _ := json.Marshal(tStream)
		json.Unmarshal(tsBytes, &stream)
	} else {
		// If stream is already in our history then we've received a notice
		// for it already. Stop processing.
		log.Println("Found duplicate stream in our history. Stopping processing.")
		return nil
	}

	category, err := n.GetCategoryByName(stream.GameName)
	if err != nil {
		log.Printf("Error looking for category %v in table: %v\n", stream.GameName, err)
		return err
	}
	if category == nil || category.Id == "" {
		log.Printf("Category %v is not in our map. Stopping processing.\n", stream.GameName)
		return nil
	}

	// TODO: Add description exclusion filter logic here

	// Add/update stream information in table
	// We do this ASAP so that we can debounce if duplicate notices come in
	err = n.PutStream(stream)
	if err != nil {
		log.Println("Could not save stream information to table. Stopping processing.")
		return err
	}

	// Fetch image data to use in each social media post
	previewImage := &utility.Image{}
	altText := fmt.Sprintf("Preview of %v's stream on Twitch.", user.DisplayName)
	dimensions := strings.Split(config.StreamThumbResolution, "x")
	width, _ := strconv.Atoi(dimensions[0])
	height, _ := strconv.Atoi(dimensions[1])

	previewImage, _ = utility.NewFromThumbnailURL(
		stream.ThumbnailURL,
		width,
		height,
		altText,
	)

	log.Printf("Starting message post goroutines.")

	postChan := make(chan utility.PostResponse)
	go discordPostRoutine(stream, previewImage, postChan)
	go mastodonPostRoutine(user, stream, category, previewImage, postChan)
	go blueskyPostRoutine(stream, category, previewImage, postChan)

	postRoutines := 3 // increase based on number of goroutines above
	for i := 0; i < postRoutines; i++ {
		resp := <-postChan
		switch resp.Platform {

		case discord.PlatformName:
			stream.DiscordPostId = resp.Id
			stream.DiscordPostUrl = resp.Url

		case bluesky.PlatformName:
			stream.BlueskyPostId = resp.Id
			stream.BlueskyPostUrl = resp.Url

		case mastodon.PlatformName:
			stream.MastodonPostId = resp.Id
			stream.MastodonPostUrl = resp.Url
		}
	}

	err = n.PutStream(stream)
	if err != nil {
		log.Printf("Failed to write stream updates after posting.")
		return err
	}

	return nil
}

func streamOfflineCallback(sub *twitch.Subscription, eventMap *map[string]string) error {
	log.Println("Entered streamOnlineCallback")

	// Unmarshal event data into helix struct
	event := helix.EventSubStreamOfflineEvent{}
	evBytes, _ := json.Marshal(eventMap)
	json.Unmarshal(evBytes, &event)

	// TODO: Log stream offline time in eventsub_notice_history

	return nil
}

func messageIsDuplicate(messageId string) bool {
	// Instantiate DynamoDB
	n, _ := nosqldb.NewClient()
	eventsub, _ := n.GetEventsubMessage(messageId)

	return eventsub.Id != ""
}

func discordPostRoutine(stream *nosqldb.StreamHistoryDatum, image *utility.Image, c chan utility.PostResponse) {
	dc, _ := discord.NewClient()
	streamUrl := fmt.Sprintf("https://twitch.tv/%v", stream.UserLogin)
	resp, err := dc.Post(dc.FormatMsg(
		stream.UserName,
		stream.GameName,
		stream.Title,
		streamUrl,
	), image)
	if err != nil {
		log.Printf("Error posting to discord: %v\n", err)
	}

	c <- *resp
}

func blueskyPostRoutine(stream *nosqldb.StreamHistoryDatum, category *nosqldb.CategoryDatum, image *utility.Image, c chan utility.PostResponse) {
	bc, _ := bluesky.NewClient()
	streamUrl := fmt.Sprintf("https://twitch.tv/%v", stream.UserLogin)

	msg := fmt.Sprintf(
		"%v is now streaming %v on Twitch: %v\n\n%v\n\n%v",
		stream.UserName,
		stream.GameName,
		streamUrl,
		stream.Title,
		strings.Join(category.BlueskyTags, " "),
	)

	resp, err := bc.Post(msg, image)
	if err != nil {
		log.Printf("Error posting to bluesky: %v\n", err)
	}

	c <- *resp
}

func mastodonPostRoutine(user *nosqldb.TwitchUserDatum, stream *nosqldb.StreamHistoryDatum, category *nosqldb.CategoryDatum, image *utility.Image, c chan utility.PostResponse) {
	mc, _ := mastodon.NewClient()
	streamUrl := fmt.Sprintf("https://twitch.tv/%v", stream.UserLogin)

	streamer := ""
	if user.MastodonUserId != "" {
		streamer = fmt.Sprintf("@%v", user.MastodonUserId)
	} else {
		streamer = stream.UserName
	}

	msg := fmt.Sprintf(
		"%v is now streaming %v on Twitch: %v\n\n%v\n\n%v",
		streamer,
		stream.GameName,
		streamUrl,
		stream.Title,
		strings.Join(category.MastodonTags, " "),
	)

	resp, err := mc.Post(msg, image)
	if err != nil {
		log.Printf("Error posting to mastodon: %v\n", err)
	}

	c <- *resp
}
