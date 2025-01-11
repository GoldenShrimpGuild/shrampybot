package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("services_test", "Test Calls for API Social Media Endpoints")
	envName := parser.Selector("e", "envName", []string{"shrampybot-dev", "shrampybot-prod"}, &argparse.Options{
		Required: false,
		Help:     "Name of the deployment environment.",
		Default:  "shrampybot-dev",
	})
	serviceName := parser.Selector("s", "serviceName", []string{"bluesky", "discord", "mastodon", "twitch"}, &argparse.Options{
		Required: true,
		Help:     "Name of desired service.",
	})
	testCase := parser.Selector("t", "test", []string{
		"profile",
		"post",
	}, &argparse.Options{
		Required: false,
		Help:     "Test case to run",
		Default:  "profile",
	})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	project, _ := readProject(envName)
	env, err := project.getAllEnv(envName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	switch *serviceName {
	case "bluesky":
		blueskyMain(testCase, &env)
	case "twitch":
		twitchMain(testCase, &env)
	}
}
