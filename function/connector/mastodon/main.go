package mastodon

import (
	"context"
	"log"
	"regexp"
	"shrampybot/config"
	"strings"
	"time"

	mast "github.com/litui/go-mastodon"
)

type Client struct {
	mh  *mast.Client
	ctx context.Context
}

func NewClient() (*Client, error) {
	var c Client
	c.ctx = context.Background()
	conf := &mast.Config{
		Server:      config.MastodonApiUrl,
		AccessToken: config.MastodonApiToken,
	}
	c.mh = mast.NewClient(conf)

	return &c, nil
}

func (c *Client) GetMappedTwitchLoginsThreaded(ch chan string) {
	log.Println("Entered function: GetMappedTwitchLoginsThreaded")
	var twitchMatch = regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?twitch\.tv\/([A-Za-z0-9_-]+)\/?`)
	var accounts []*mast.AdminAccount
	var pg mast.Pagination
	pg.Limit = 200

	log.Println("Iterating through pages of accounts...")
	for {
		acct, err := c.mh.AdminViewAccounts(c.ctx, &mast.AdminViewAccountsInput{}, &pg)
		if err != nil {
			log.Fatal(err)
			break
		}
		accounts = append(accounts, acct...)
		if pg.MaxID == "" {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	log.Println("Matching regex to Twitch URLs in profiles...")
	for _, acct := range accounts {
		for _, field := range acct.Account.Fields {
			if twitchMatch.MatchString(field.Value) {
				values := twitchMatch.FindStringSubmatch(field.Value)
				if len(values) == 2 {
					ch <- strings.ToLower(values[1])
					break
				}

			}
		}
	}
	close(ch)
	log.Println("Exiting function: GetMappedTwitchLoginsThreaded")
}
