package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	mast "github.com/litui/go-mastodon"
)

func mastodonMain(testCase *string, env *map[string]string) error {
	var err error
	ctx := context.Background()

	var (
		mastodonApiUrl   = (*env)["MASTODON_API_URL"]
		mastodonApiToken = (*env)["MASTODON_API_TOKEN"]
	)

	conf := &mast.Config{
		Server:      mastodonApiUrl,
		AccessToken: mastodonApiToken,
	}

	mc := mast.NewClient(conf)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(3)
	// }

	switch *testCase {
	case "profile":
		err = mastodonTestProfile(ctx, mc)
	case "post":
		err = mastodonTestPost(ctx, mc)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(5)
	}
	return nil
}

func mastodonTestProfile(ctx context.Context, mc *mast.Client) error {
	user, err := mc.AccountLookup(ctx, "@litui")
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

func mastodonTestPost(ctx context.Context, mc *mast.Client) error {
	toot := &mast.Toot{
		Status: "Hello world.",
	}

	stat, err := mc.PostStatus(ctx, toot)
	if err != nil {
		return err
	}

	output, err := json.MarshalIndent(stat, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
