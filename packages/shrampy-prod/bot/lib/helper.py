import asyncio
import os
import re
import urllib.request
from datetime import datetime


def chunkify(input_list, count=100):
    for i in range(0, len(input_list), count):
        yield input_list[i:i + count]

async def async_chunkify(input_list, count=100):
    i = 0
    while (i < len(input_list)):
        imod = (i + 1) % count
        if imod == 0:
            yield (input_list[i:i+count])
        elif i >= len(input_list):
            yield (input_list[i:])
            break
        i += 1
        asyncio.sleep(0.0)

def twitch_date(twitch_time):
    return datetime.strptime(twitch_time, "%Y-%m-%dT%H:%M:%S.%fZ")

def fetch_twitch_thumb(thumb_url):
    retval = b""
    proper_url = re.sub(
        r"\{width\}x\{height\}",
        os.environ.get("STREAM_THUMB_RESOLUTION"),
        thumb_url
    )
    with urllib.request.urlopen(proper_url) as f:
        retval = f.read()
    return retval
