# ShrampyBot - The Golden Shrimp Guild NOW Bot

This is a bot being used and developed actively by the [Golden Shrimp Guild](https://gsg.live) for announcing when streamers on the gsg.live Twitch group and Mastodon instance go live.

This is a "serverless function" running in Python on DigitalOcean Functions. It leverages the Mastodon API for tooting, Discord API for announcing in the guild's NOW channel (along with crossposts), and Twitch API EventSub Webhooks to receive event notifications. It uses the Twitch REST API endpoint for getting more detailed information.

A `project-example.yml` script is present with example options. The authoritative `project.yml` file needs to be in the directory to deploy using the doctl utility from Digital Ocean (available via apt or homebrew).

Several scripts exist in the `tools` directory which should be run from the root of the project when used. Please ensure you have built and activated a local python virtual environment and have installed everything in the `requirements.txt` file:

* `subscription_stats.py` - Get a complete breakdown of what event subscriptions are active and for how many users. Also breaks down discrepancies between the stored Twitch user list in the s3-bucket and the current subscriptions.
* `reconcile_subscriptions.py` - Takes the most efficient path to reconciling the existing subscriptions with any user additions or removals needed. This was written to take some of the load off the serverless function's most intensive call and some of its logic will likely be moved directly into the function once it's using the latest version of the TwitchAPI python library.

There are a number of settings which can be updated via REST calls. These will be documented below in the fullness of time.