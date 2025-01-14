package bluesky

import (
	"context"
	"log"
	"time"

	"shrampybot/config"
	"shrampybot/utility"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	blueSky "github.com/tailscale/go-bluesky"
)

const (
	PlatformName = "bluesky"
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

func (c *Client) Post(msg string, thumb *utility.Image) (*utility.PostResponse, error) {
	var err error
	now := time.Now()

	var imageBlob *util.LexBlob

	log.Println("Uploading image blob to Bluesky.")
	// Pre-upload the image "blob"
	err = c.bc.CustomCall(func(api *xrpc.Client) error {
		ctx := context.Background()

		resp, err := atproto.RepoUploadBlob(ctx, api, thumb.GetReader())
		if err != nil {
			return err
		}
		imageBlob = resp.Blob

		return nil
	})
	if err != nil {
		return &utility.PostResponse{}, err
	}

	embed := bsky.FeedPost_Embed{
		EmbedImages: &bsky.EmbedImages{
			Images: []*bsky.EmbedImages_Image{
				{
					Alt:   thumb.AltText,
					Image: imageBlob,
					AspectRatio: &bsky.EmbedDefs_AspectRatio{
						Height: int64(thumb.Height),
						Width:  int64(thumb.Width),
					},
				},
			},
		},
	}

	post := bsky.FeedPost{
		Text:      msg,
		CreatedAt: now.Format(time.RFC3339),
		Embed:     &embed,
		Facets:    []*bsky.RichtextFacet{},
		Reply:     &bsky.FeedPost_ReplyRef{},
	}

	postResponse := &utility.PostResponse{}

	log.Println("Posting message to Bluesky.")
	// Gotta use the CustomCall API to post since the library is minimal
	err = c.bc.CustomCall(func(api *xrpc.Client) error {
		ctx := context.Background()
		record, err := atproto.RepoCreateRecord(ctx, api, &atproto.RepoCreateRecord_Input{
			Repo:       api.Auth.Did,
			Collection: "app.bsky.feed.post",
			Record: &util.LexiconTypeDecoder{
				Val: &post,
			},
		})
		if err != nil {
			return err
		}

		log.Printf("Posted to Bluesky, id: %v\n", record.Uri)
		postResponse.Platform = PlatformName
		postResponse.Id = record.Cid
		postResponse.Url = record.Uri

		return nil
	})

	return postResponse, err
}
