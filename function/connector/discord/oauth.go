package discord

import "github.com/bwmarrin/discordgo"

type OAuthClient struct {
	dc    *discordgo.Session
	ready bool
}

func NewOAuthClient(accessToken string) (*OAuthClient, error) {
	dc, err := discordgo.New("Bearer " + accessToken)
	if err != nil {
		return nil, err
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
