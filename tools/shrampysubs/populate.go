package main

import (
	"fmt"
	"os"
)

func taskPopulate(config *ShrampyConfig) {
	var err error
	fmt.Printf("ShrampyBot User Population (Remote)\n\n")

	logins := getUserLogins(config)
	if logins == nil || len(*logins) == 0 {
		fmt.Println("ShrampyBot collection has not yet been populated.")
		fmt.Println("Try running with \"-t populate\"")
		os.Exit(6)
	}
	fmt.Printf("Twitch logins tracked before: %v\n", len(*logins))

	newLogins, err := populateUserLogins(config)
	if err != nil || newLogins == nil || len(*logins) == 0 {
		fmt.Printf("Population failed: %v\n", err)
		os.Exit(8)
	}
	fmt.Printf("Twitch logins tracked after: %v\n", len(*logins))
}
