package main

import (
	"fmt"
	"os"
)

func taskReport(config *ShrampyConfig) {
	var err error
	fmt.Printf("ShrampyBot Subscription Report\n\n")

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

	subs, err := twitchGetSubs(tc)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(7)
	}
	fmt.Printf("Total eventsub subscriptions: %v\n", len(*subs))
}
