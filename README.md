# ShrampyBot - The Golden Shrimp Guild NOW Bot

This is a bot being used and developed actively by the [Golden Shrimp Guild](https://gsg.live) for announcing when streamers on the gsg.live Twitch group and Mastodon instance go live.

This is a "serverless function" running in Python on DigitalOcean Functions. It leverages the Mastodon API for tooting, Discord API for announcing in the guild's NOW channel (along with crossposts), and Twitch API EventSub Webhooks to receive event notifications. It uses the Twitch REST API endpoint for getting more detailed information.

A `project-example.yml` script is present with example options. The authoritative `project.yml` file needs to be in the directory to deploy using the doctl utility from Digital Ocean (available via apt or homebrew).

There are a number of settings which can be updated via REST calls. These will be documented below in the fullness of time.