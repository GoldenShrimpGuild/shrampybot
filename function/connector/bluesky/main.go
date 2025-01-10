package bluesky

import (
	"bytes"
	"context"
	"time"

	"shrampybot/config"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	blueSky "github.com/tailscale/go-bluesky"
)

type Client struct {
	bc *blueSky.Client
}

func NewClient() (*Client, error) {
	ctx := context.Background()

	bc, err := blueSky.Dial(ctx, blueSky.ServerBskySocial)
	if err != nil {
		return &Client{}, err
	}
	defer bc.Close()

	err = bc.Login(ctx, config.BlueskyLogin, config.BlueskyPassword)
	if err != nil {
		return &Client{}, err
	}

	return &Client{
		bc: bc,
	}, nil
}

func (c *Client) Post(msg string, image []byte, imageAltText string) (string, error) {
	now := time.Now()
	imageReader := bytes.NewReader(image)

	var imageBlob *util.LexBlob

	err := c.bc.CustomCall(func(api *xrpc.Client) error {
		ctx := context.Background()

		resp, err := atproto.RepoUploadBlob(ctx, api, imageReader)
		if err != nil {
			return err
		}
		imageBlob = resp.Blob

		return nil
	})
	if err != nil {
		return "", err
	}

	embed := bsky.FeedPost_Embed{
		EmbedImages: &bsky.EmbedImages{
			Images: []*bsky.EmbedImages_Image{
				{
					Alt:   imageAltText,
					Image: imageBlob,
				},
			},
		},
	}

	post := bsky.FeedPost{
		Text:      "Hello World",
		CreatedAt: now.Format(time.RFC3339),
		Embed:     &embed,
		Facets:    []*bsky.RichtextFacet{},
		Reply:     &bsky.FeedPost_ReplyRef{},
	}

	// Gotta use the CustomCall API to post since the library is minimal
	err = c.bc.CustomCall(func(api *xrpc.Client) error {
		ctx := context.Background()
		_, err := atproto.RepoCreateRecord(ctx, api, &atproto.RepoCreateRecord_Input{
			Repo:       api.Auth.Did,
			Collection: "app.bsky.feed.post",
			Record: &util.LexiconTypeDecoder{
				Val: &post,
			},
		})

		return err
	})

	return "", err
}
