package mastodon

import (
	"shrampybot/config"

	mast "github.com/mattn/go-mastodon"
)

type Client struct {
	mh *mast.Client
}

func NewClient() (*Client, error) {
	conf := &mast.Config{
		Server:      config.MastodonApiUrl,
		AccessToken: config.MastodonApiToken,
	}
	mh := mast.NewClient(conf)
}
