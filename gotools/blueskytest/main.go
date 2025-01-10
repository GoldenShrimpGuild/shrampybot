package main

import (
	"context"
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	blueSky "github.com/tailscale/go-bluesky"
)

func main() {
	parser := argparse.NewParser("blueskytest", "Test Calls for Bluesky")
	packageName := parser.Selector("p", "packageName", []string{"shrampybot-dev", "shrampybot-prod"}, &argparse.Options{
		Required: false,
		Help:     "Name of the package as it exists in project.yml.",
		Default:  "shrampybot-dev",
	})
	projectPath := parser.String("j", "projectPath", &argparse.Options{
		Required: false,
		Help:     "Path (relative or absolute) to project.yml",
		Default:  projectFilename,
	})
	test := parser.Selector("t", "test", []string{
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

	project, _ := readProject(projectPath)
	env, err := project.getAllEnv(packageName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	var (
		blueskyLogin    = env["BLUESKY_LOGIN"]
		blueskyPassword = env["BLUESKY_PASSWORD"]
	)

	ctx := context.Background()

	bc, err := blueSky.Dial(ctx, blueSky.ServerBskySocial)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
	}
	defer bc.Close()

	err = bc.Login(ctx, blueskyLogin, blueskyPassword)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(4)
	}

	switch *test {
	case "profile":
		err = testProfile(&ctx, bc)
	case "post":
		err = testPost(&ctx, bc)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(5)
	}
}
