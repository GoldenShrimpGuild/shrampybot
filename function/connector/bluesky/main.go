package bluesky

import (
	"context"
	"log"
	"regexp"
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
	bc  *blueSky.Client
	ctx context.Context
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
		bc:  bc,
		ctx: ctx,
	}, nil
}

func (c *Client) Post(msg string, thumb *utility.Image) (*utility.PostResponse, error) {
	var err error
	now := time.Now()

	var imageBlob *util.LexBlob

	log.Println("Uploading image blob to Bluesky.")
	// Pre-upload the image "blob"
	err = c.bc.CustomCall(func(api *xrpc.Client) error {
		resp, err := atproto.RepoUploadBlob(c.ctx, api, thumb.GetReader())
		if err != nil {
			log.Printf("Issue uploading image blob: %v\n", err)
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
		Facets:    compileFacets(msg),
		// Reply: &bsky.FeedPost_ReplyRef{},
	}

	postResponse := &utility.PostResponse{}

	log.Println("Posting message to Bluesky.")
	// Gotta use the CustomCall API to post since the library is minimal
	err = c.bc.CustomCall(func(api *xrpc.Client) error {
		record, err := atproto.RepoCreateRecord(c.ctx, api, &atproto.RepoCreateRecord_Input{
			Repo:       api.Auth.Did,
			Collection: "app.bsky.feed.post",
			Record: &util.LexiconTypeDecoder{
				Val: &post,
			},
		})
		if err != nil {
			log.Printf("Error posting to Bluesky: %v\n", err)
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

func compileFacets(msg string) []*bsky.RichtextFacet {
	facets := []*bsky.RichtextFacet{}

	urlRegex := regexp.MustCompile(`https?://[A-Za-z0-9._\-/]+`)
	hashTagRegex := regexp.MustCompile(`#([A-Za-z0-9]([^ \n]*[A-Za-z0-9]{1}))`)
	// atUserRegex := regexp.MustCompile(`@([A-Za-z0-9\-.]+)`)

	urlIndices := urlRegex.FindAllIndex([]byte(msg), 10)
	urlMatches := urlRegex.FindAllStringSubmatch(msg, 10)
	for i, indices := range urlIndices {
		facet := bsky.RichtextFacet{
			Index: &bsky.RichtextFacet_ByteSlice{
				ByteStart: int64(indices[0]),
				ByteEnd:   int64(indices[1]),
			},
			Features: []*bsky.RichtextFacet_Features_Elem{{
				RichtextFacet_Link: &bsky.RichtextFacet_Link{
					LexiconTypeID: "app.bsky.richtext.facet#link",
					Uri:           urlMatches[i][0],
				},
			}},
		}
		facets = append(facets, &facet)
	}

	hashTagIndices := hashTagRegex.FindAllIndex([]byte(msg), 10)
	hashTagMatches := hashTagRegex.FindAllStringSubmatch(msg, 10)
	for i, indices := range hashTagIndices {
		facet := bsky.RichtextFacet{
			Index: &bsky.RichtextFacet_ByteSlice{
				ByteStart: int64(indices[0]),
				ByteEnd:   int64(indices[1]),
			},
			Features: []*bsky.RichtextFacet_Features_Elem{{
				RichtextFacet_Tag: &bsky.RichtextFacet_Tag{
					LexiconTypeID: "app.bsky.richtext.facet#tag",
					Tag:           hashTagMatches[i][1],
				},
			}},
		}
		facets = append(facets, &facet)
	}

	// atUserIndices := atUserRegex.FindAllIndex([]byte(msg), 10)
	// atUserMatches := atUserRegex.FindAllStringSubmatch(msg, 10)

	return facets
}
