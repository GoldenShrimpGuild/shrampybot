package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func discordMain(testCase *string, env *map[string]string) error {
	var (
		discordToken = (*env)["DISCORD_TOKEN"]
		// discordGuild = (*env)["DISCORD_GUILD"]
		discordChannel = (*env)["DISCORD_CHANNEL"]
	)

	dc, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
	}
	defer dc.Close()

	dc.Identify.Intents = discordgo.IntentGuildMembers | discordgo.IntentsGuilds

	err = dc.Open()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(4)
	}

	switch *testCase {
	case "profile":
		err = discordTestProfile(dc)
	case "post":
		err = discordTestPost(&discordChannel, dc)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(5)
	}
	return nil
}

func discordTestProfile(dc *discordgo.Session) error {
	var err error
	user, err := dc.User("@me")
	if err != nil {
		return err
	}

	output, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}

func discordTestPost(channel *string, dc *discordgo.Session) error {
	fh, err := os.Open("../../attach/yui01.jpg")
	if err != nil {
		return err
	}
	defer fh.Close()
	var files []*discordgo.File
	files = append(files, &discordgo.File{
		Name:        "image.jpg",
		ContentType: "image/jpeg",
		Reader:      fh,
	})

	res, err := dc.ChannelMessageSendComplex(*channel, &discordgo.MessageSend{
		Content: "**Yui** is now streaming **Music** on Twitch:\nlol ADHD. A challenger appears!!!\n\nhttps://www.tbs.co.jp/anime/k-on/",
		Files:   files,
		Flags:   discordgo.MessageFlagsSuppressEmbeds | discordgo.MessageFlagsCrossPosted,
	})
	if err != nil {
		return err
	}

	output, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
