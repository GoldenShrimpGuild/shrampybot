package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/litui/helix/v3"
)

func twitchMain(testCase *string, env *map[string]string) {
	var (
		twitchApiKey    = (*env)["TWITCH_API_KEY"]
		twitchApiSecret = (*env)["TWITCH_API_SECRET"]
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

	switch *testCase {
	case "profile":
		err = twitchTestProfile(twC)
	case "post":
		err = errors.New("test case 'post' unsupported for twitch")
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(5)
	}
}

func twitchTestProfile(twC *helix.Client) error {
	users, err := twC.GetUsers(&helix.UsersParams{
		Logins: []string{"litui"},
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(6)
	}

	output, _ := json.MarshalIndent(users.Data.Users, "", "  ")
	fmt.Println(string(output))

	return nil
}
