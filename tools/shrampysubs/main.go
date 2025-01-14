package main

import (
	"fmt"
	"io"
	"os"

	"github.com/akamensky/argparse"
	"gopkg.in/yaml.v3"
)

type ShrampyConfig struct {
	Url               string `yaml:"shrampybot_url"`
	AdminToken        string `yaml:"shrampybot_admin_token"`
	TwitchClientId    string `yaml:"twitch_client_id"`
	TwitchSecretKey   string `yaml:"twitch_secret_key"`
	TwitchEventSecret string `yaml:"twitch_event_secret"`
}

func main() {
	var err error
	parser := argparse.NewParser("shrampysubs", "ShrampyBot Subscription Manager")

	configFile := parser.File("c", "configFile", os.O_RDONLY, 0, &argparse.Options{
		Required: false,
		Default:  "./shrampysubs-dev.yml",
		Help:     "YAML config file",
	})
	task := parser.Selector("t", "task", []string{
		"report",
		"populate",
		"reconcile",
		"unsubscribe_all",
	}, &argparse.Options{
		Required: true,
		Help:     "Task to execute",
	})

	err = parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	defer configFile.Close()
	buf, err := io.ReadAll(configFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}

	config := &ShrampyConfig{}
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(3)
	}

	switch *task {
	case "report":
		taskReport(config)
	case "populate":
		taskPopulate(config)
	case "reconcile":
		taskReconcile(config)
	case "unsubscribe_all":
		taskUnsubscribeAll(config)
	}
}
