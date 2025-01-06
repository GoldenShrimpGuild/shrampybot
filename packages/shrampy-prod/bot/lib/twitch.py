import asyncio
import os
import logging
from logging import DEBUG, WARN, ERROR, INFO
from functools import cached_property
from twitchAPI import twitch, eventsub, types, object
from lib.helper import chunkify, async_chunkify

# MAX_SUBS_PER_CALL = 30


class TwitchHandler:
    TWITCH_MESSAGE_ID = 'twitch-eventsub-message-id'
    TWITCH_MESSAGE_TIMESTAMP = 'twitch-eventsub-message-timestamp'
    TWITCH_MESSAGE_SIGNATURE = 'twitch-eventsub-message-signature'
    TWITCH_MESSAGE_TYPE = 'twitch-eventsub-message-type'

    def __init__(self):
        logger = logging.getLogger("TwitchHandler")
        self.l = logger.log
        self.l(INFO, "Initializing TwitchHandler")
        self._team_name = os.environ["TWITCH_TEAM_NAME"]
        self._event_subs = []
        self._total_event_subs = 0
        self._team_info = None # twitch.ChannelTeam()

    @cached_property
    def _th(self):
        """Twitch Handler"""
        th = twitch.Twitch(
            app_id=os.environ["TWITCH_API_KEY"],
            app_secret=os.environ["TWITCH_API_SECRET"],
            authenticate_app=True,
            target_app_auth_scope=None
        )
        return th
    
    async def authenticate(self):
        return await self._th.authenticate_app(scope=[])

    @cached_property
    def _eh(self):
        """Event Handler"""

        eh = eventsub.EventSub(
            callback_url=os.environ["EVENTSUB_URL"],
            api_client_id=os.environ["TWITCH_API_KEY"],
            port=443,
            twitch=self._th
        )
        eh.secret = os.environ["TWITCH_EVENT_SECRET"]
        eh.unsubscribe_on_stop = False
        eh.wait_for_subscription_confirm = False
        return eh

    async def get_event_subs(self):
        self._event_subs = []

        eventsubs_raw: object.GetEventSubSubscriptionResult = await self._th.get_eventsub_subscriptions()
        self._event_subs.extend(eventsubs_raw.data)
        self._total_event_subs = eventsubs_raw.total

        return self._event_subs

    async def get_team_info(self) -> twitch.ChannelTeam:
        self.l(DEBUG, "Accessed team_info.")
        if self._team_name:
            team = await self._th.get_teams(name=self._team_name)
            if team:
                self._team_info = team

        return self._team_info

    @property
    async def _event_streamon_users(self):
        return [
            i.condition["broadcaster_user_id"]
            for i in await self.get_event_subs()
            if i.type == "stream.online"
        ]

    @property
    async def _event_streamoff_users(self):
        return [
            i.condition["broadcaster_user_id"]
            for i in await self.get_event_subs()
            if i.type == "stream.offline"
        ]

    @property
    async def _event_raid_users(self):
        return [
            i.condition["from_broadcaster_user_id"]
            for i in await self.get_event_subs()
            if i.type == "channel.raid"
        ]

    async def get_users(self, user_logins=[]):
        # Limit to 100 per request
        for user_chunk in chunkify(user_logins):
            yield {
                i.login: i.to_dict()
                async for i in self._th.get_users(
                    logins=user_chunk
                )
            }
            await asyncio.sleep(0.0)
    
    async def get_stream_by_user_id(self, user_id):
        async for stream in self._th.get_streams(first=1, user_id=user_id):
            return stream

    async def _null_cb(self, data):
        '''NULL callback for use with listen_ eventsub calls.'''
        pass

    async def subscribe_to_events(self, twitch_users):
        new_streamon_count = 0
        new_streamoff_count = 0

        existing_online = await self._event_streamon_users
        existing_offline = await self._event_streamoff_users

        for u, data in twitch_users.items():
            uid = data["id"]

            if not uid in existing_online:
                try:
                    await self._eh.listen_stream_online(
                        broadcaster_user_id=uid,
                        callback=self._null_cb
                    )
                    new_streamon_count += 1
                except types.EventSubSubscriptionConflict as e:
                    self.l(DEBUG, "Event conflict on uid '{}', event '{}'"
                        .format(uid, "stream.online")
                    )
                except types.EventSubSubscriptionTimeout as e:
                    self.l(DEBUG, f"Timeout trying to subscribe to event: {e}")
                except types.EventSubSubscriptionError as e:
                    self.l(DEBUG, f"Error trying to subscribe to event: {e}")
                except types.TwitchBackendException as e:
                    self.l(DEBUG, f"Twitch backend exception trying to subscribe to event: {e}")

            if not uid in existing_offline:
                try:
                    await self._eh.listen_stream_offline(
                        broadcaster_user_id=uid,
                        callback=self._null_cb
                    )
                    new_streamoff_count += 1
                except types.EventSubSubscriptionConflict as e:
                    self.l(DEBUG, "Event conflict on uid '{}', event '{}'"
                        .format(uid, "stream.offline")
                    )
                except types.EventSubSubscriptionTimeout as e:
                    self.l(DEBUG, f"Timeout trying to subscribe to event: {e}")
                except types.EventSubSubscriptionError as e:
                    self.l(DEBUG, f"Error trying to subscribe to event: {e}")
                except types.TwitchBackendException as e:
                    self.l(DEBUG, f"Twitch backend exception trying to subscribe to event: {e}")
                
        newbies = new_streamon_count + \
            new_streamoff_count
        if newbies:
            try:
                self._event_subs = []
            except AttributeError:
                pass

        return await self.get_event_subs()

    async def unsubscribe_all_events(self):
        await self._eh.unsubscribe_all()
        return True
