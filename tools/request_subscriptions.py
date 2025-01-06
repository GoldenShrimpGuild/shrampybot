#!/usr/bin/env python

import asyncio
import boto3
import yaml
import json
import sys
from botocore.config import Config as BotoConfig
from twitchAPI.twitch import Twitch
from twitchAPI.type import EventSubSubscriptionConflict
from twitchAPI.eventsub.webhook import EventSubWebhook

async def _null_cb(data):
    pass

async def main():
    project_conf: dict = yaml.load(open('project.yml'), Loader=yaml.SafeLoader)
    try:
        package: dict = [i for i in project_conf["packages"] if i["name"] == sys.argv[1]][0]
    except IndexError:
        print("You must specify which environment you will use; eg: 'shrampy-dev'")
        sys.exit(1)

    functions: dict = [i for i in package["functions"] if i["name"] == "bot"][0]

    twitch_api_key: str = package["environment"]["TWITCH_API_KEY"]
    twitch_api_secret: str = package["environment"]["TWITCH_API_SECRET"]
    twitch_event_secret: str = package["environment"]["TWITCH_EVENT_SECRET"]
    aws_access_key: str = project_conf["environment"]["AWS_ACCESS_KEY_ID"]
    aws_secret_key: str = project_conf["environment"]["AWS_SECRET_ACCESS_KEY"]
    aws_endpoint_url: str = package["environment"]["AWS_ENDPOINT_URL"]
    aws_bucket: str = package["environment"]["AWS_BUCKET"]
    aws_region: str = package["environment"]["AWS_DEFAULT_REGION"]

    eventsub_url: str = functions["environment"]["EVENTSUB_URL"]

    s3: boto3.Session = boto3.client(
        's3',
        aws_access_key_id=aws_access_key,
        aws_secret_access_key=aws_secret_key,
        endpoint_url=aws_endpoint_url
    )
    s3obj: dict = s3.get_object(Bucket=aws_bucket, Key="twitch_users")
    document = ""
    for line in s3obj["Body"]:
        document += line.decode(encoding="utf-8")
    s3.close()
    twitch_users: dict = json.loads(document)

    th = await Twitch(twitch_api_key, twitch_api_secret)
    eh: EventSubWebhook = EventSubWebhook(eventsub_url, 443, twitch=th)
    eh.secret = twitch_event_secret
    eh.unsubscribe_on_stop = False
    eh.wait_for_subscription_confirm = False

    # await eh.unsubscribe_all()

    count = 0

    for login, user in twitch_users.items():
        print(f"Requesting for user: {login}")
        try:
            await eh.listen_stream_online(user["id"], _null_cb)
        except EventSubSubscriptionConflict:
            print(f"Already subscribed to online status for {user["id"]}.")
        try:
            await eh.listen_stream_offline(user["id"], _null_cb)
        except EventSubSubscriptionConflict:
            print(f"Already subscribed to offline status for {user["id"]}.")
        count += 1

    await th.close()

    print(f"Requested subscriptions for {count} users.")

if __name__ == "__main__":
    asyncio.run(main())