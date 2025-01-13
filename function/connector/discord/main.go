package discord

import (
	"fmt"
	"log"
	"shrampybot/config"
	"shrampybot/utility"

	"github.com/bwmarrin/discordgo"
)

const (
	PlatformName = "discord"
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

func (c *Client) Post(msg string, image utility.Image) (*utility.PostResponse, error) {
	var files []*discordgo.File

	files = append(files, &discordgo.File{
		Name:        "image.jpg",
		ContentType: "image/jpeg",
		Reader:      image.GetReader(),
	})

	postResponse := &utility.PostResponse{}

	log.Printf("Sending Discord message...")
	res, err := c.dc.ChannelMessageSendComplex(config.DiscordChannel, &discordgo.MessageSend{
		Content: msg,
		Files:   files,
		Flags:   discordgo.MessageFlagsSuppressEmbeds | discordgo.MessageFlagsCrossPosted,
	})
	if err != nil {
		return postResponse, err
	}

	log.Printf("Discord message sent: %v\n", res.ID)

	postResponse.Platform = PlatformName
	postResponse.Id = res.ID
	postResponse.Url = fmt.Sprintf(
		"https://discord.com/channels/%v/%v/%v",
		config.DiscordGuild,
		config.DiscordChannel,
		res.ID,
	)

	return postResponse, nil
}
