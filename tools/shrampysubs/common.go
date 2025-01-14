package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/litui/helix/v3"
)

type UserLoginsBody struct {
	Data *[]string `json:"data"`
}

func connectToTwitch(config *ShrampyConfig) (*helix.Client, error) {
	tc, err := helix.NewClient(&helix.Options{
		ClientID:     config.TwitchClientId,
		ClientSecret: config.TwitchSecretKey,
	})
	if err != nil {
		return &helix.Client{}, err
	}

	resp, err := tc.RequestAppAccessToken([]string{})
	if err != nil {
		return &helix.Client{}, err
	}
	tc.SetAppAccessToken(resp.Data.AccessToken)

	return tc, err
}

func getUserLogins(config *ShrampyConfig) *[]string {
	body := UserLoginsBody{}
	client := http.Client{}
	request, _ := http.NewRequest("GET", config.Url+"admin/collection", bytes.NewReader([]byte{}))
	request.Header.Add("Content-type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("bearer %v", config.AdminToken))
	resp, err := client.Do(request)
	if err != nil {
		return body.Data
	}

	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)

	json.Unmarshal(bodyBytes, &body)
	return body.Data
}

func populateUserLogins(config *ShrampyConfig) (*[]string, error) {
	body := UserLoginsBody{}
	client := http.Client{}
	request, _ := http.NewRequest("PATCH", config.Url+"admin/collection", bytes.NewReader([]byte{}))
	request.Header.Add("Content-type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("bearer %v", config.AdminToken))
	resp, err := client.Do(request)
	if err != nil {
		return body.Data, err
	}

	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)

	json.Unmarshal(bodyBytes, &body)
	return body.Data, nil
}

func twitchGetUsers(tc *helix.Client, logins *[]string) (*[]helix.User, error) {
	users := []helix.User{}

	// 100 item maximum for each call to GetUsers
	for subList := range slices.Chunk(*logins, 100) {
		resp, err := tc.GetUsers(&helix.UsersParams{
			Logins: subList,
		})
		if err != nil {
			return &users, err
		}
		users = append(users, resp.Data.Users...)
	}

	return &users, nil
}

func twitchGetSubs(tc *helix.Client) (*[]helix.EventSubSubscription, error) {
	subs := []helix.EventSubSubscription{}
	var after string

	for {
		resp, err := tc.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{
			After: after,
		})
		if err != nil {
			return &subs, err
		}
		subs = append(subs, resp.Data.EventSubSubscriptions...)

		if resp.Data.Pagination.Cursor == "" {
			break
		}
		after = resp.Data.Pagination.Cursor
	}

	return &subs, nil
}
