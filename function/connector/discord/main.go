package discord

import (
	"bytes"
	"log"
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
	var files []*discordgo.File
	imageReader := bytes.NewReader(image)

	files = append(files, &discordgo.File{
		Name:        "image.jpg",
		ContentType: "image/jpeg",
		Reader:      imageReader,
	})

	log.Printf("Sending Discord message...")
	res, err := c.dc.ChannelMessageSendComplex(config.DiscordChannel, &discordgo.MessageSend{
		Content: msg,
		Files:   files,
		Flags:   discordgo.MessageFlagsSuppressEmbeds | discordgo.MessageFlagsCrossPosted,
	})
	if err != nil {
		return "", err
	}

	return res.ID, nil
}
