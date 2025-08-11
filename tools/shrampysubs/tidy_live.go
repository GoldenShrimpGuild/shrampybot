package main

import (
	"fmt"
	"os"
)

// Clean up any streams listed as live that are no longer running.
func taskTidyLive(config *ShrampyConfig) {
	var err error
	fmt.Printf("ShrampyBot Tidy Live Streams\n\n")

	tc, err := connectToTwitch(config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(5)
	}

	liveStreams := getLiveStreams(config)
	liveLogins := []string{}
	for _, s := range liveStreams {
		liveLogins = append(liveLogins, s.UserLogin)
	}
	fmt.Printf("Found %v live streams in ShrampyBot.\n", len(liveStreams))

	actuallyLiveStreams, err := twitchGetStreams(tc, &liveLogins)
	if err != nil {
		fmt.Printf("Encountered an error trying to get live streams from Twitch.\n\n")
		os.Exit(10)
	}
	fmt.Printf("Found %v corresponding live streams in Twitch.\n", len(*actuallyLiveStreams))

	offlineStreams := []*LiveStream{}
	for _, ols := range liveStreams {
		foundLs := false
		for _, als := range *actuallyLiveStreams {
			if als.UserLogin == ols.UserLogin {
				foundLs = true
			}
		}

		if !foundLs {
			offlineStreams = append(offlineStreams, ols)
		}
	}

	for _, os := range offlineStreams {
		err = endLiveStream(os.ID, config)
		if err != nil {
			fmt.Printf("Could not end active stream for %v on ShrampyBot.\n", os.UserLogin)
		} else {
			fmt.Printf("Ended live stream for %v.\n", os.UserLogin)
		}
	}
}
