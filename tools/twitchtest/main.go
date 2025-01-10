package main

import (
	"encoding/json"
	"fmt"
	"os"

	helix "github.com/litui/helix/v3"
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
		twitchApiKey    = env["TWITCH_API_KEY"]
		twitchApiSecret = env["TWITCH_API_SECRET"]
	)

	twC, err := helix.NewClient(&helix.Options{
		ClientID:     twitchApiKey,
		ClientSecret: twitchApiSecret,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
	}

	resp, err := twC.RequestAppAccessToken([]string{})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(4)
	}
	twC.SetAppAccessToken(resp.Data.AccessToken)

	users, err := twC.GetUsers(&helix.UsersParams{
		Logins: []string{"litui"},
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(5)
	}

	output, _ := json.MarshalIndent(users, "", "  ")
	fmt.Println(string(output))
}
