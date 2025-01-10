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

// discordgo seems to only want to deal with URLs rather than bytes
func (c *Client) Post(msg string, imageUrl string) (string, error) {
	// discordgo.MessageEmbed{
	// 	Image: &discordgo.MessageEmbedImage{

	// 	}
	// }

	// c.dc.ChannelMessageSendEmbed()

	return "", nil
}
