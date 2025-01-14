package main

import (
	"fmt"
	"os"
	"time"

	"github.com/litui/helix/v3"
)

func taskReconcile(config *ShrampyConfig) {
	var err error
	fmt.Printf("ShrampyBot Reconcile Subscriptions\n\n")

	tc, err := connectToTwitch(config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(5)
	}

	logins := getUserLogins(config)
	if logins == nil || len(*logins) == 0 {
		fmt.Println("ShrampyBot collection has not yet been populated.")
		fmt.Println("Try running with \"-t populate\"")
		os.Exit(6)
	}
	fmt.Printf("Twitch logins tracked: %v\n", len(*logins))

	users, err := twitchGetUsers(tc, logins)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(7)
	}

	subs, err := twitchGetSubs(tc)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(8)
	}
	fmt.Printf("Total eventsub subscriptions before: %v\n", len(*subs))

	rateLimit := 800

	sleepRate := time.Duration(60.0/float32(rateLimit)*1000) * time.Millisecond

	time.Sleep(1 * time.Second)

	for _, user := range *users {
		subStreamOnlineExists := false
		subStreamOfflineExists := false

		for _, sub := range *subs {
			if sub.Type == "stream.online" && sub.Condition.BroadcasterUserID == user.ID {
				subStreamOnlineExists = true
			}
			if sub.Type == "stream.offline" && sub.Condition.BroadcasterUserID == user.ID {
				subStreamOfflineExists = true
			}
		}

		if !subStreamOnlineExists {
			fmt.Printf("Subscribing %v to %v\n", user.Login, "stream.online")
			tc.CreateEventSubSubscription(&helix.EventSubSubscription{
				Type:    "stream.online",
				Version: "1",
				Condition: helix.EventSubCondition{
					BroadcasterUserID: user.ID,
				},
				Transport: helix.EventSubTransport{
					Method:   "webhook",
					Callback: config.Url + "event/webhook",
					Secret:   config.TwitchEventSecret,
				},
			})
			time.Sleep(sleepRate)
		}
		if !subStreamOfflineExists {
			fmt.Printf("Subscribing %v to %v\n", user.Login, "stream.offline")
			tc.CreateEventSubSubscription(&helix.EventSubSubscription{
				Type:    "stream.offline",
				Version: "1",
				Condition: helix.EventSubCondition{
					BroadcasterUserID: user.ID,
				},
				Transport: helix.EventSubTransport{
					Method:   "webhook",
					Callback: config.Url + "event/webhook",
					Secret:   config.TwitchEventSecret,
				},
			})
			time.Sleep(sleepRate)
		}
	}

	subs, err = twitchGetSubs(tc)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(8)
	}
	fmt.Printf("Total eventsub subscriptions after: %v\n", len(*subs))
}
