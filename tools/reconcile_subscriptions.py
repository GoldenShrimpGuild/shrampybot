#!/usr/bin/env python

import asyncio
import boto3
import yaml
import json
import sys
from botocore.config import Config as BotoConfig
from twitchAPI.twitch import Twitch
from twitchAPI.type import EventSubSubscriptionConflict, TwitchResourceNotFound
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
    twitch_ids = set([v["id"] for i,v in twitch_users.items()])

    th = await Twitch(twitch_api_key, twitch_api_secret)
    eh: EventSubWebhook = EventSubWebhook(eventsub_url, 443, twitch=th)
    eh.secret = twitch_event_secret
    eh.unsubscribe_on_stop = False
    eh.wait_for_subscription_confirm = False

    total_subs = []

    next_cursor = None
    while True:
        es_subs = await th.get_eventsub_subscriptions(after=next_cursor)
        total_subs.extend(es_subs.data)
        next_cursor = es_subs.current_cursor()
        if len(es_subs.data) < 100:
            break

    sub_ids = set([i.condition["broadcaster_user_id"] for i in total_subs])
    outlier_subs = [i for i in total_subs if not i.transport["callback"].startswith(eventsub_url)]
    types = set([i.type for i in total_subs])

    new_twitch_ids = twitch_ids.difference(sub_ids)
    pending_delete = sub_ids.difference(twitch_ids)

    print(f"{sys.argv[1]} Twitch list contains {len(twitch_ids)} ids.")
    print(f"There are {len(new_twitch_ids)} ids awaiting event subscription.")
    print(f"There are {len(pending_delete)} ids awaiting event unsubscription.")

    new_id_count = 0
    new_sub_count = 0

    for login, user in twitch_users.items():
        if user["id"] in new_twitch_ids:
            print(f"Requesting subscriptions for user: {login}")
            try:
                await eh.listen_stream_online(user["id"], _null_cb)
                new_sub_count += 1
            except EventSubSubscriptionConflict:
                print(f"Already subscribed to online status for {user["id"]}.")
            try:
                await eh.listen_stream_offline(user["id"], _null_cb)
                new_sub_count += 1
            except EventSubSubscriptionConflict:
                print(f"Already subscribed to offline status for {user["id"]}.")
            
            new_id_count += 1

    print(f"Requested {new_sub_count} subscriptions across {new_id_count} ids.")
    
    deleted_ids = set()
    deleted_sub_count = 0

    for sub in total_subs:
        if sub.condition.get("broadcaster_user_id", "") in pending_delete:
            print(f"Requesting unsubscribe for id: {sub.condition["broadcaster_user_id"]}")
            try:
                await th.delete_eventsub_subscription(sub.id)
                deleted_sub_count += 1
            except TwitchResourceNotFound:
                print(f"Could not delete subscription {sub.type} for {sub.condition["broadcaster_user_id"]}.")
            
            deleted_ids.add(sub.condition["broadcaster_user_id"])

    print(f"Deleted {deleted_sub_count} subscriptions across {len(deleted_ids)} ids.")

    await th.close()

if __name__ == "__main__":
    asyncio.run(main())