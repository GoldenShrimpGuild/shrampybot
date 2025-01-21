package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"shrampybot/utility/nosqldb"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type OAuthClient struct {
	dc    *discordgo.Session
	ready bool
}

func NewOAuthClient(oauth *nosqldb.DiscordOAuthDatum) (*OAuthClient, error) {
	// Check accessToken expiry and refresh if expired or nearly so, updating original datum
	if time.Now().After(oauth.ExpiresAt.Add(-60 * time.Second)) {
		hc := http.Client{}

		query_data := url.Values{}
		query_data.Set("grant_type", "refresh_token")
		query_data.Set("refresh_token", oauth.RefreshToken)

		request, err := http.NewRequest("POST", "https://discord.com/api/oauth2/token", strings.NewReader(query_data.Encode()))
		if err != nil {
			log.Println("Could not create new refresh Discord token request.")
			return &OAuthClient{}, err
		}
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err := hc.Do(request)
		if err != nil {
			log.Printf("Unsuccessful request to Discord for new token: %v\n", err)
			return &OAuthClient{}, err
		}
		if resp.StatusCode > 399 {
			log.Printf("Error when requesting new Discord token: %v\n", resp.StatusCode)
			bodyBytes, _ := io.ReadAll(resp.Body)
			log.Printf("New discord token response body: %v\n", string(bodyBytes))
			return &OAuthClient{}, fmt.Errorf("error when requesting new discord token")
		}
		responseMap := map[string]any{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &responseMap)

		oauth.TokenType = responseMap["token_type"].(string)
		oauth.AccessToken = responseMap["access_token"].(string)
		oauth.AccessTokenEnc = ""
		oauth.AccessTokenIV = ""
		oauth.RefreshToken = responseMap["refresh_token"].(string)
		oauth.RefreshTokenEnc = ""
		oauth.RefreshTokenIV = ""
		oauth.ExpiresIn = responseMap["expires_in"].(float64)
		oauth.ExpiresAt = time.Now().Add(time.Duration(oauth.ExpiresIn) * time.Second)
		oauth.Scope = responseMap["scope"].(string)

		oauth.Refreshed = true
		// Don't forget to save the datum after returning.
	}

	// And try again
	dc, err := discordgo.New("Bearer " + oauth.AccessToken)
	if err != nil {
		log.Println("Could not create session with new Access Token.")
		return &OAuthClient{}, err
	}
	defer dc.Close()

	client := OAuthClient{
		dc:    dc,
		ready: false,
	}
	dc.AddHandler(client.isReady)

	return &client, nil
}

func (c *OAuthClient) isReady(s *discordgo.Session, r *discordgo.Ready) {
	c.ready = true
}

func (c *OAuthClient) GetSelf() (*discordgo.User, error) {
	return c.dc.User("@me")
}

func (c *OAuthClient) GetConnections() ([]*discordgo.UserConnection, error) {
	return c.dc.UserConnections()
}
