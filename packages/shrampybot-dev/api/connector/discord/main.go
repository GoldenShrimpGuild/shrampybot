package discord

import (
	"shrampybot/config"

	"github.com/bwmarrin/discordgo"
)

type Client struct {
	dc    *discordgo.Session
	ready bool
}

func NewClient() (*Client, error) {
	dc, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		return nil, err
	}
	defer dc.Close()

	client := Client{
		dc:    dc,
		ready: false,
	}
	dc.AddHandler(client.isReady)

	return &client, nil
}

func (c *Client) isReady(s *discordgo.Session, r *discordgo.Ready) {
	c.ready = true
}

func (c *Client) Post(msg string, image []byte) (string, error) {

}
