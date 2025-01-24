package discord

import (
	"errors"
	"fmt"
	"log"
	"shrampybot/config"
	"shrampybot/utility"
	"slices"

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

func (c *BotClient) UserIsAdmin(id string) bool {
	membership, err := c.GetGuildMember(id)
	if err != nil {
		return false
	}
	if slices.Contains(membership.Roles, config.DiscordAdminRole) {
		return true
	}

	return false
}

func (c *BotClient) UserIsDev(id string) bool {
	membership, err := c.GetGuildMember(id)
	if err != nil {
		return false
	}
	if slices.Contains(membership.Roles, config.DiscordDevRole) {
		return true
	}

	return false
}

func (c *BotClient) LocalScopesFromMembership(id string) ([]string, error) {
	// determine scopes
	scopes := []string{}
	membership, err := c.GetGuildMember(id)
	if err != nil {
		log.Printf("Error retrieving scopes for user %v: %v\n", id, err)
		return scopes, err
	}
	if membership != nil {
		scopes = append(scopes, "login")
		scopes = append(scopes, "self")
		if slices.Contains(membership.Roles, config.DiscordDevRole) {
			scopes = append(scopes, "dev")
		}
		if slices.Contains(membership.Roles, config.DiscordAdminRole) {
			scopes = append(scopes, "admin")
		}
	}
	if len(scopes) == 0 {
		log.Printf("No membership on the Discord found for user %v.\n", id)
		return scopes, errors.New("no membership found for user")
	}
	return scopes, nil
}
