package main

import (
	"fmt"
	"os"
	"time"
)

func taskUnsubscribeAll(config *ShrampyConfig) {
	var err error
	fmt.Printf("ShrampyBot Unsubscribe All\n\n")

	tc, err := connectToTwitch(config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(5)
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

	for _, sub := range *subs {
		tc.RemoveEventSubSubscription(sub.ID)

		time.Sleep(sleepRate)
	}

	subs, err = twitchGetSubs(tc)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(8)
	}
	fmt.Printf("Total eventsub subscriptions after: %v\n", len(*subs))
}
