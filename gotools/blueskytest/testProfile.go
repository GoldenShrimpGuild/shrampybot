package main

import (
	"context"
	"encoding/json"
	"fmt"

	bSky "github.com/tailscale/go-bluesky"
)

func testProfile(ctx *context.Context, bc *bSky.Client) error {
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
