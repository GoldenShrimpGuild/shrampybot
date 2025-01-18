package config

import "os"

var (
	BotName                  = os.Getenv("BOT_NAME")
	GsgAdminToken            = os.Getenv("GSG_ADMIN_TOKEN")
	StreamupDebounceInterval = os.Getenv("STREAMUP_DEBOUNCE_INTERVAL")
	StreamThumbResolution    = os.Getenv("STREAM_THUMB_RESOLUTION")

	AwsAccessKeyId     = os.Getenv("AWS_ACCESS_KEY_ID")
	AwsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	AwsEndpointUrl     = os.Getenv("AWS_ENDPOINT_URL")
	AwsDefaultRegion   = os.Getenv("AWS_DEFAULT_REGION")
	AwsBucket          = os.Getenv("AWS_BUCKET")

	TwitchApiKey      = os.Getenv("TWITCH_API_KEY")
	TwitchApiSecret   = os.Getenv("TWITCH_API_SECRET")
	TwitchEventSecret = os.Getenv("TWITCH_EVENT_SECRET")
	TwitchTeamName    = os.Getenv("TWITCH_TEAM_NAME")
	EventsubUrl       = os.Getenv("EVENTSUB_URL")

	MastodonApiUrl   = os.Getenv("MASTODON_API_URL")
	MastodonApiToken = os.Getenv("MASTODON_API_TOKEN")
	MastodonPostMode = os.Getenv("MASTODON_POST_MODE")

	BlueskyLogin    = os.Getenv("BLUESKY_LOGIN")
	BlueskyPassword = os.Getenv("BLUESKY_PASSWORD")

	DiscordToken   = os.Getenv("DISCORD_TOKEN")
	DiscordGuild   = os.Getenv("DISCORD_GUILD")
	DiscordChannel = os.Getenv("DISCORD_CHANNEL")

	DBCryptKey = os.Getenv("DB_CRYPT_KEY")
)
