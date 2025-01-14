package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	blueSky "github.com/tailscale/go-bluesky"
)

func blueskyMain(testCase *string, env *map[string]string) {
	var (
		blueskyLogin    = (*env)["BLUESKY_LOGIN"]
		blueskyPassword = (*env)["BLUESKY_PASSWORD"]
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

	switch *testCase {
	case "profile":
		err = blueskyTestProfile(&ctx, bc)
	case "post":
		err = blueskyTestPost(&ctx, bc)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(5)
	}
}

func blueskyTestProfile(ctx *context.Context, bc *blueSky.Client) error {
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

func blueskyTestPost(ctx *context.Context, bc *blueSky.Client) error {
	now := time.Now()

	fh, err := os.Open("assets/adi-goldstein-unsplash.jpg")
	if err != nil {
		return err
	}
	defer fh.Close()
	fmt.Println("Opened image file.")

	var imageBlob *util.LexBlob
	var createResp *atproto.RepoCreateRecord_Output

	err = bc.CustomCall(func(api *xrpc.Client) error {
		resp, err := atproto.RepoUploadBlob(*ctx, api, fh)
		if err != nil {
			return err
		}
		imageBlob = resp.Blob

		return nil
	})
	if err != nil {
		return err
	}
	fmt.Println("Uploaded image blob.")

	embed := bsky.FeedPost_Embed{
		EmbedImages: &bsky.EmbedImages{
			Images: []*bsky.EmbedImages_Image{
				{
					Alt:   "Alt text",
					Image: imageBlob,
					AspectRatio: &bsky.EmbedDefs_AspectRatio{
						Height: 720,
						Width:  1280,
					},
				},
			},
		},
	}

	post := bsky.FeedPost{
		Text:      "Hello World",
		CreatedAt: now.Format(time.RFC3339),
		Embed:     &embed,
		// Facets:    []*bsky.RichtextFacet{},
		// Reply:     &bsky.FeedPost_ReplyRef{},
	}

	// Gotta use the CustomCall API to post since the library is minimal
	err = bc.CustomCall(func(api *xrpc.Client) error {
		createResp, err = atproto.RepoCreateRecord(*ctx, api, &atproto.RepoCreateRecord_Input{
			Repo:       api.Auth.Did,
			Collection: "app.bsky.feed.post",
			Record: &util.LexiconTypeDecoder{
				Val: &post,
			},
		})

		return err
	})
	if err != nil {
		return err
	}
	fmt.Println("Posted to Bluesky.")

	output, err := json.MarshalIndent(createResp, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
