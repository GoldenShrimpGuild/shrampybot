package main

import (
	"context"
	"encoding/json"
	"fmt"

	blueSky "github.com/tailscale/go-bluesky"
)

func testProfile(ctx *context.Context, bc *blueSky.Client) error {
	profile, err := bc.FetchProfile(*ctx, "litui.ca")
	if err != nil {
		return err
	}

	output, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
