package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"

	"github.com/litui/helix/v3"
)

type UserLoginsBody struct {
	Data *[]string `json:"data"`
}

type LiveStream struct {
	helix.Stream       `tstype:",extends,required"`
	DiscordPostId      string    `json:"discord_post_id,omitempty"`
	DiscordPostUrl     string    `json:"discord_post_url,omitempty"`
	MastodonPostId     string    `json:"mastodon_post_id,omitempty"`
	MastodonPostUrl    string    `json:"mastodon_post_url,omitempty"`
	BlueskyPostId      string    `json:"bluesky_post_id,omitempty"`
	BlueskyPostUrl     string    `json:"bluesky_post_url,omitempty"`
	ShrampybotFiltered bool      `json:"shrampybot_filtered"`
	EndedAt            time.Time `json:"ended_at,omitempty"`
}

type PutLiveStreamRequest struct {
	EndNow bool `json:"endNow,omitempty"`
}

type PutLiveStreamResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type LiveStreamResponse struct {
	Count int           `json:"count"`
	Data  []*LiveStream `json:"data"`
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

func getLiveStreams(config *ShrampyConfig) []*LiveStream {
	body := LiveStreamResponse{}
	client := http.Client{}
	request, _ := http.NewRequest("GET", config.Url+"public/stream", bytes.NewReader([]byte{}))
	request.Header.Add("Content-type", "application/json")
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

func endLiveStream(id string, config *ShrampyConfig) error {
	respBody := PutLiveStreamResponse{}
	reqBody := PutLiveStreamRequest{
		EndNow: true,
	}
	client := http.Client{}

	marshalledReq, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	request, _ := http.NewRequest("PUT", fmt.Sprintf("%sadmin/stream/status/%s", config.Url, id), bytes.NewReader(marshalledReq))
	request.Header.Add("Content-type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("bearer %v", config.AdminToken))
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &respBody)

	if respBody.Status != "success" {
		return fmt.Errorf("did not receive a success message")
	}
	return nil
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

func twitchGetStreams(tc *helix.Client, logins *[]string) (*[]helix.Stream, error) {
	streams := []helix.Stream{}

	// 100 item maximum for each call to GetUsers
	for subList := range slices.Chunk(*logins, 100) {
		resp, err := tc.GetStreams(&helix.StreamsParams{
			UserLogins: subList,
		})
		if err != nil {
			return &streams, err
		}
		streams = append(streams, resp.Data.Streams...)
	}

	return &streams, nil
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
