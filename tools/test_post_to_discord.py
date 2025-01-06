#!/usr/bin/env python

import asyncio
import hikari
import yaml
import json
import sys

async def send_message(dh: hikari.impl.RESTClientImpl, channel, msg, image=None):
    image_attach = None
    if image:
        image_attach = hikari.Bytes(
            image,
            "image.jpg",
            mimetype="image/jpeg"
        )
    
    message = await dh.create_message(
        channel=channel,
        content=msg,
        attachment=image_attach,
        flags=hikari.MessageFlag.SUPPRESS_EMBEDS
    )

    try:
        await dh.crosspost_message(
            channel=channel,
            message=message
        )
    except hikari.HikariError as e:
        print("Couldn't crosspost. No big deal if this is a dev channel.")

    print("Message posted to channel.")

async def main():
    project_conf: dict = yaml.load(open('project.yml'), Loader=yaml.SafeLoader)
    try:
        package: dict = [i for i in project_conf["packages"] if i["name"] == sys.argv[1]][0]
    except IndexError:
        print("You must specify which environment you will use; eg: 'shrampy-dev'")
        sys.exit(1)

    function: dict = [i for i in package["functions"] if i["name"] == "bot"][0]

    discord_token: str = package["environment"]["DISCORD_TOKEN"]
    discord_guild: str = package["environment"]["DISCORD_GUILD"]
    discord_channel: str = function["environment"]["DISCORD_CHANNEL"]

    rest = hikari.RESTApp()
    await rest.start()
    dh: hikari.impl.RESTClientImpl = rest.acquire(discord_token, hikari.applications.TokenType.BOT)
    dh.start()
    print(await dh.fetch_my_user())

    with open("./assets/technotronic.png", "b+r") as fh:
        png = fh.read()

    await send_message(dh, discord_channel, "Technotronic is now streaming Music.", png)

    await dh.close()
    await rest.close()

if __name__ == "__main__":
    asyncio.run(main())