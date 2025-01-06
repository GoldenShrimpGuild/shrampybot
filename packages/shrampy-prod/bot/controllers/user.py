import asyncio
import logging
from logging import DEBUG, WARN, ERROR, INFO
from lib.s3 import S3
from lib.mastodon import MastodonHandler
from lib.twitch import TwitchHandler
from controllers.generic import GenericController
from auth.admin import AdminAuthenticator
from auth.twitch import TwitchAuthenticator

class UserController(GenericController):
    def __init__(self, *args, **kwargs):
        logger = logging.getLogger("UserController")
        self.l = logger.log
        self.l(INFO, "Initializing UserController.")
        self._admin_auth = AdminAuthenticator(*args, **kwargs)
        self._twitch_auth = TwitchAuthenticator(*args, **kwargs)
        super().__init__(*args, **kwargs)

        self._mh = MastodonHandler()
        self._th = TwitchHandler()
        self._s3 = S3()

    async def entry_point(self):
        if not await self._admin_auth.check_credentials():
            return {
                "body": self._router.call_error(14),
                "statusCode": 403
            }
        
        await self._th.authenticate()
        return await super().entry_point()

    async def _delete__user_custom(self):
        """Delete one or more entries from custom pairs"""
        custom_users = self._body.get("custom_users", [])
        if not custom_users:
            return {
                "body": self._router.call_error(15, "custom_users"),
                "statusCode": 400
            }


        # Ensure we validate the incoming json.
        for cu in custom_users:
            if not cu.get("mastodon_id") and not cu.get("twitch_login"):
                return {
                    "body": self._router.call_error(16),
                    "statusCode": 400
                }

        for cu in custom_users:
            self._s3.del_custom_pair(
                cu["mastodon_id"],
                cu["twitch_login"]
            )
        self._s3.refresh_maps()

        return {"body": {"success": "true"}}

    async def _get__user_custom(self):
        """Fetch list of custom pairs"""
        return {"body": self._s3.custom_pairs, "statusCode": 200}

    async def _get__user_map(self):
        """Endpoint: /user/map; Method: GET
        
        Return s3 key-value store of account
        mapping
        """

        full_map = {
            "status": "success"
        }
        status_code = 200
        try:
            full_map['m_to_t'] = self._s3.mt_map
            full_map['t_to_m'] = self._s3.tm_map
        except:
            full_map["status"] = "failure"
            status_code = 200
        
        return {"body": full_map, "statusCode": status_code}

    async def _get__user_twitch(self):
        out_dict = self._s3.twitch_users
        out_dict["status"] = "success"

        return {"body": out_dict, "statusCode": 200}

    async def _patch__user_map(self):
        """Endpoint: /user/map; Method: PATCH
        
        Update s3 key-value store with current
        account mapping.
        """

        # return {"body": {"status": "Made it to user_map function"}, "statusCode": 200}

        full_map = self._mh.twitch_map
        self._s3.commit_maps(full_map["m_to_t"], full_map["t_to_m"])

        # Compare results
        fm_mt_count = len(full_map["m_to_t"].keys())
        fm_tm_count = len(full_map["t_to_m"].keys())
        s3_mt_count = len(self._s3.mt_map.keys())
        s3_tm_count = len(self._s3.tm_map.keys())

        out_dict = {
            "map_counts_in": {
                "mastodon_to_twitch": fm_mt_count,
                "twitch_to_mastodon": fm_tm_count
            },
            "map_counts_out": {
                "mastodon_to_twitch": s3_mt_count,
                "twitch_to_mastodon": s3_tm_count
            }
        }

        if fm_mt_count <= s3_mt_count and \
                fm_tm_count <= s3_tm_count:
            out_dict["status"] = "success"
            status_code = 200
        else:
            out_dict["status"] = "failure"
            status_code = 500

        return {"body": out_dict, "statusCode": status_code}

    async def _patch__user_twitch(self):
        """Load Twitch users into S3 for each we have mapped."""

        # Retrieve cached mastodon users
        user_list = list(set(self._s3.tm_map.keys()))
        print("Unique users found in S3 userlist: {}".format(len(user_list)))
        # Retrieve list of Twitch team members (if any)
        team_users = [u.user_login for u in (await self._th.get_team_info()).users]
        print("Users found in Twitch Team: {}".format(len(team_users)))

        user_list.extend(team_users)
        logins = list(set(user_list))
        user_map = {}
        async for ul in self._th.get_users(user_logins=logins):
            user_map.update(ul)
        self._s3.commit_twitch_users(users=user_map)

        return {"body": {"users": list(set(user_map.keys()))}, "statusCode": 200}

    async def _post__user_custom(self):
        """Submit a new entry to the custom user map.
        
        Mastodon format is user@instance (no initial @). 
        @instance can be dropped for the local instance.

        custom_users = [{
            "mastodon_id": "",
            "twitch_login": ""
        }]
        """
        custom_users = self._body.get("custom_users", [])
        if not custom_users:
            return self._router.call_error(15, "custom_users")

        # Ensure we validate the incoming json.
        for cu in custom_users:
            if not cu.get("mastodon_id") and not cu.get("twitch_login"):
                return self._router.call_error(16)

        # Load 'em in.
        for cu in custom_users:
            self._s3.add_custom_pair(
                cu["mastodon_id"],
                cu["twitch_login"]
            )
        self._s3.refresh_maps()
            
        return {"body": {"status": "success"}, "statusCode": 200}