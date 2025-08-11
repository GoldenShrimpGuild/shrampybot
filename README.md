# ShrampyBot - The Golden Shrimp Guild NOW Bot

This is an API being used and developed actively by the [Golden Shrimp Guild](https://gsg.live) primarily for announcing when streamers on the gsg.live Twitch group and Mastodon instance go live.

## Backend

ShrampyBot's backend is a "serverless function" running in Go 1.23 on AWS Lambda with DynamoDB. Its main functionality involves management of [Twitch EventSub](https://dev.twitch.tv/docs/eventsub/) subscriptions and receipt of eventsub webhook notifications.

When receiving `stream.online` notifications for members of the [GSG Twitch Team](https://www.twitch.tv/team/gsg) in relevant (mostly music) categories, it will store a record of the stream status and send notifications out on Discord, [Mastodon](https://soc.gsg.live/@shrampybot), and [Bluesky](https://bsky.app/profile/shrampybot.gsg.live).

Most of the ShrampyBot API is not public, though there are some basic exceptions.

## Frontend

The ShrampyBot frontend is written in Vue3 + TypeScript. It is also not intended to be a public-facing UI for the most part, but again there are exceptions. At present the most useful public endpoints are:

- [GSG Multi-Twitch](https://goldenshrimpguild.github.io/shrampybot/multi) - Inspired by https://www.multitwitch.tv/, this makes use of the ShrampyBot API and GSG Events API to dynamically tile embedded streams for each of the currently online GSG members.

