import os
import logging
from logging import DEBUG, WARN, ERROR, INFO
from functools import cached_property
import hikari
from hikari import HikariError as DiscordError


class DiscordHandler:

    def __init__(self):
        logger = logging.getLogger("DiscordHandler")
        self.l = logger.log
        self.l(INFO, "Initializing DiscordHandler")

    async def get_dh(self):
        dh = None
        self._rest = hikari.RESTApp()
        await self._rest.start()
        dh = self._rest.acquire(
            token=os.environ['DISCORD_TOKEN'],
            token_type=hikari.applications.TokenType.BOT
        )
        dh.start()
        return dh
    
    async def close_dh(self, dh):
        await dh.close()
        await self._rest.close()

    @cached_property
    async def _me(self):
        dh = await self.get_dh()
        my_user = await dh.fetch_my_user()
        return my_user

    async def send_message(self, msg, image=None):
        image_attach = None
        dh = await self.get_dh()

        if image:
            image_attach = hikari.Bytes(
                image,
                "image.jpg",
                mimetype="image/jpeg"
            )
        message = await dh.create_message(
            channel=os.environ["DISCORD_CHANNEL"],
            content=msg,
            attachment=image_attach,
            flags=hikari.MessageFlag.SUPPRESS_EMBEDS
        )
        try:
            await dh.crosspost_message(
                channel=os.environ["DISCORD_CHANNEL"],
                message=message
            )
        except DiscordError as e:
            self.l(WARN, "Could not crosspost message {}".format(
                message.id
            ))
        
        return message