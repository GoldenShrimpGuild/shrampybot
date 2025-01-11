package twitch

import (
	"shrampybot/config"
	"slices"

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

func (c *Client) GetTeamMembers() (*[]helix.TeamUser, error) {
	resp, err := c.tc.GetTeams(&helix.TeamsParams{
		Name: config.TwitchTeamName,
	})
	if err != nil {
		return &[]helix.TeamUser{}, err
	}

	if len(resp.Data.Teams) > 0 {
		return &resp.Data.Teams[0].Users, err
	}

	return &[]helix.TeamUser{}, nil
}

func (c *Client) GetUsers(logins *[]string) (*[]helix.User, error) {
	users := []helix.User{}

	// 100 item maximum for each call to GetUsers
	for subList := range slices.Chunk(*logins, 100) {
		resp, err := c.tc.GetUsers(&helix.UsersParams{
			Logins: subList,
		})
		if err != nil {
			return &users, err
		}
		users = append(users, resp.Data.Users...)
	}

	return &users, nil
}
