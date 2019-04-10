import itertools
import logging
from concurrent.futures import ThreadPoolExecutor
from pprint import pprint

import bson
import click
import pendulum
import requests


def iter_date(
    start_date: pendulum.datetime, end_date: pendulum.datetime, chunk_size=59
):
    if end_date < start_date:
        raise ValueError(
            "start_date:%s should not large than end_date:%s", start_date, end_date
        )
    while start_date <= end_date:
        new_end_date = min(start_date.add(days=chunk_size), end_date)
        yield start_date, new_end_date
        start_date = new_end_date.add(days=1)


class ShanbayClient:
    def __init__(self, token):
        self.token = token
        self._session = requests.Session()
        if token:
            self._session.cookies.set("auth_token", token)
        self._session.headers.update(
            {
                "X-Device": "android.ticktick, MI 8, 5003, 5cab671122d4db0dde5a941c, android_xiaomi_dd,",
                "User-Agent": "Dalvik/2.1.0 (Linux; U; Android 9; MI 8 MIUI/9.3.28)",
            }
        )

    def checkin_calendar(self, user_id, start_date, end_date):
        """
        区间两边闭合，最多返回 60 条，按打卡天数降序
        """
        return self._session.get(
            "https://apiv3.shanbay.com/uc/checkin/calendar/dates",
            params={
                "user_id": user_id,
                "start_date": start_date.to_date_string() or "",
                "end_date": end_date.to_date_string() or "",
            },
        ).json()


class DiDaClient:
    def __init__(self, token):
        self.token = token
        self._session = requests.Session()
        self._session.headers.update(
            {
                "Authorization": "OAuth " + self.token,
                "X-Device": "android.ticktick, MI 8, 5003, 5cab671122d4db0dde5a941c, android_xiaomi_dd,",
                "User-Agent": "Dalvik/2.1.0 (Linux; U; Android 9; MI 8 MIUI/9.3.28)",
            }
        )

    def get_habits(self):
        return self._session.get("https://api.dida365.com/api/v2/habits").json()

    def batch_habits(self, data):
        # 无法修改 createdTime
        return self._session.post(
            "https://api.dida365.com/api/v2/habits/batch", json=data
        ).json()

    def batch_checkin(self, data):
        return self._session.post(
            "https://api.dida365.com/api/v2/habitCheckins/batch", json=data
        ).json()


def calendar(user_id, start_date, end_date):
    """
    返回 API 一样的结果，但没有长度限制，日志按打卡日期递增排序
    """
    client = ShanbayClient(None)
    start_date = pendulum.parse(start_date)
    end_date = pendulum.parse(end_date)
    checkin_days_num = -1

    def fn(args):
        result = client.checkin_calendar(user_id, *args)
        nonlocal checkin_days_num
        checkin_days_num = result["checkin_days_num"]
        print("fetched:", args, len(result["logs"]))
        return result["logs"]

    with ThreadPoolExecutor() as executor:
        logs = sorted(
            itertools.chain(*executor.map(fn, iter_date(start_date, end_date))),
            key=lambda l: l["date"],
        )
        return {"checkin_days_num": checkin_days_num, "logs": logs}


def log_to_checkin(habit_id, logs):
    for log in logs:
        date = pendulum.parse(log["date"], tz="Asia/Shanghai")
        checkin = {
            "id": str(bson.ObjectId()),
            "habitId": habit_id,
            "checkinStamp": date.year * 10000 + date.month * 100 + date.day,
            "checkinTime": date.isoformat(),
        }
        yield checkin


@click.command()
@click.pass_context
@click.option("-t", "--token", type=click.STRING, required=True, help="滴答清单认证token")
@click.option("-u", "--user_id", type=click.STRING, required=True, help="滴答清单认证token")
@click.option(
    "-s", "--start_date", type=click.STRING, default="2016-01-01", help="起始日期"
)
@click.option(
    "-e",
    "--end_date",
    type=click.STRING,
    default=pendulum.now().to_date_string(),
    help="结束日期",
)
@click.option("-d", "--delete", type=click.BOOL, help="是否删除相同名字的习惯")
def cli(ctx, token, user_id, start_date, end_date, delete):
    logger: logging.Logger = ctx.obj.logger

    result = calendar(user_id, start_date, end_date)
    checkin_days_num, logs = result["checkin_days_num"], result["logs"]
    assert checkin_days_num == len(logs)
    if checkin_days_num == 0:
        logger.warning("No checkin logs, exit")
        return

    client = DiDaClient(token)

    name = "扇贝打卡"
    if delete:
        habits = client.get_habits()
        delete_ids = []
        for h in habits:
            if h["name"] == name:
                print("delete matched habit")
                pprint(h)
                delete_ids.append(h["id"])
        result = client.batch_habits({"add": [], "delete": delete_ids, "update": []})
        pprint(result)

    habit_id = str(bson.ObjectId())
    tz = "Asia/Shanghai"
    habit = {
        "name": name,
        "id": habit_id,
        "createdTime": pendulum.parse(logs[0]["date"], tz=tz).isoformat(),
        "modifiedTime": pendulum.parse(logs[1]["date"], tz=tz).isoformat(),
        "totalCheckIns": checkin_days_num,
        "color": "#209E85",
        "encouragement": "Shanbay, feel the change",
        "iconRes": "habit_learn_words",
        "sortOrder": -1374389534720,
        "status": 0,
    }
    result = client.batch_habits({"add": [habit], "delete": [], "update": []})
    pprint(result)

    result = client.batch_checkin(
        {"add": list(log_to_checkin(habit_id, logs)), "delete": [], "update": []}
    )
    pprint(result)
