package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	bSky "github.com/tailscale/go-bluesky"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("You must specify the package name, eg: shrampybot-dev")
		os.Exit(1)
	}

	project, _ := readProject()
	env, err := project.getAllEnv(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	var (
		blueskyLogin    = env["BLUESKY_LOGIN"]
		blueskyPassword = env["BLUESKY_PASSWORD"]
	)

	ctx := context.Background()

	bc, err := bSky.Dial(ctx, bSky.ServerBskySocial)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer bc.Close()

	err = bc.Login(ctx, blueskyLogin, blueskyPassword)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	profile, err := bc.FetchProfile(ctx, "litui.ca")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
	}

	output, _ := json.MarshalIndent(profile, "", "  ")
	fmt.Println(string(output))
}
