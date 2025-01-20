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

type BotClient struct {
	dc    *discordgo.Session
	ready bool
}

func NewBotClient() (*BotClient, error) {
	dc, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		return nil, err
	}
	defer dc.Close()

	client := BotClient{
		dc:    dc,
		ready: false,
	}
	dc.AddHandler(client.isReady)

	return &client, nil
}

func (c *BotClient) isReady(s *discordgo.Session, r *discordgo.Ready) {
	c.ready = true
}

func (c *BotClient) GetGuildMember(id string) (*discordgo.Member, error) {
	return c.dc.GuildMember(config.DiscordGuild, id)
}

func (c *BotClient) FormatMsg(userName string, category string, title string, url string) string {
	return fmt.Sprintf(
		"**%v** is now streaming **%v** on Twitch:\n%v\n\n%v",
		userName,
		category,
		title,
		url,
	)
}

func (c *BotClient) Post(msg string, image *utility.Image) (*utility.PostResponse, error) {
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
		Flags:   discordgo.MessageFlagsSuppressEmbeds,
	})
	if err != nil {
		return postResponse, err
	}

	_, err = c.dc.ChannelMessageCrosspost(config.DiscordChannel, res.ID)
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
