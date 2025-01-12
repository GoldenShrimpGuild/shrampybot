package twitch

import (
	"encoding/json"
	"shrampybot/config"
	"slices"
	"strings"

	"github.com/litui/helix/v3"
)

type Client struct {
	tc *helix.Client
}

func NewClient() (*Client, error) {
	tc, err := helix.NewClient(&helix.Options{
		ClientID:     config.TwitchApiKey,
		ClientSecret: config.TwitchApiSecret,
	})
	if err != nil {
		return &Client{}, err
	}

	resp, err := tc.RequestAppAccessToken([]string{})
	if err != nil {
		return &Client{}, err
	}
	tc.SetAppAccessToken(resp.Data.AccessToken)

	return &Client{
		tc: tc,
	}, err
}

func (c *Client) GetTeamMembers() (*[]map[string]string, error) {
	resp, err := c.tc.GetTeams(&helix.TeamsParams{
		Name: config.TwitchTeamName,
	})
	if err != nil {
		return &[]map[string]string{}, err
	}

	// if len(resp.Data.Teams) < 1 {
	// 	return &[]map[string]any{}, nil
	// }

	usersMap := []map[string]string{}
	uByte, _ := json.Marshal(resp.Data.Teams[0].Users)
	json.Unmarshal(uByte, &usersMap)

	return &usersMap, nil
}

func (c *Client) GetTeamMemberLoginsThreaded(ch chan string) {
	tms, _ := c.GetTeamMembers()

	for _, tm := range *tms {
		ch <- strings.ToLower(tm["user_login"])
	}
	close(ch)
}

func (c *Client) GetUsers(logins *[]string) (*[]map[string]string, error) {
	users := []helix.User{}

	// 100 item maximum for each call to GetUsers
	for subList := range slices.Chunk(*logins, 100) {
		resp, err := c.tc.GetUsers(&helix.UsersParams{
			Logins: subList,
		})
		if err != nil {
			return &[]map[string]string{}, err
		}
		users = append(users, resp.Data.Users...)
	}

	usersMap := []map[string]string{}
	uByte, _ := json.Marshal(users)
	json.Unmarshal(uByte, &usersMap)

	return &usersMap, nil
}
