"""
抓取抖音名片并推送到钉钉群
"""
import json
import os
import subprocess
import time
from datetime import datetime

import requests

ENCODING = "utf8"
PROJECT_URL = "https://github.com/elonzh/skr"


def render_msg(user, history=None):
    """
    根据名片数据生成消息内容, 具体规则见:

    https://open-doc.dingtalk.com/docs/doc.htm?treeId=257&articleId=105735&docType=1
    """
    if not history:
        text = (
            "#### @所有人, 今天 **{user[NickName]}** 涨粉了吗?\n"
            "TA(ID: {user[ID]}) 的数据如下:\n"
            "> 关注数: {user[FocusNumStr]:>7}\n"
            "> 粉丝数: {user[FollowerNumStr]:>7}\n"
            "> 点赞数: {user[LikesNumStr]:>7}\n"
            "> 作品数: {user[PostNumStr]:>7}\n"
            "> 喜欢数: {user[LikedNumStr]:>7}\n".format(user=user)
        )
    else:
        updated_at = datetime.fromtimestamp(history["UpdatedAt"])
        sentences = [
            "#### @所有人, 今天 **{user[NickName]}** 涨粉了吗?\n"
            "自 {updated_at}, TA(ID: {user[ID]}) 的数据变化如下:\n".format(
                updated_at=updated_at.strftime("%m-%e %H:%M"), user=user
            )
        ]
        for name, num_key, num_str_key in (
            ("关注数", "FocusNum", "FocusNumStr"),
            ("粉丝数", "FollowerNum", "FollowerNumStr"),
            ("点赞数", "LikesNum", "LikesNumStr"),
            ("作品数", "PostNum", "PostNumStr"),
            ("喜欢数", "LikedNum", "LikedNumStr"),
        ):
            changes = user[num_key] - history[num_key]
            if changes > 0:
                symbol = "🔺"
            elif changes < 0:
                symbol = "🔻"
            else:
                symbol = "➖"
            sentences.append(
                "> {name}: {num_str:<7} {symbol} {changes}\n".format(
                    name=name, num_str=user[num_str_key], symbol=symbol, changes=changes
                )
            )
            text = "\n".join(sentences)
    return {
        "msgtype": "actionCard",
        "actionCard": {
            "title": user["NickName"] + " 涨粉了吗?",
            "text": text,
            "hideAvatar": "true",
            "btnOrientation": "1",
            "btns": [
                {"title": "💃 查看详情", "actionURL": user["URL"]},
                {"title": "🌟 Star", "actionURL": PROJECT_URL},
            ],
        },
        "at": {"isAtAll": "true"},
    }


def main(config_path, user_histories_path):
    print("配置文件路径:", config_path)
    print("历史数据路径:", user_histories_path)
    # 读取配置文件
    with open(config_path) as fp:
        config = json.load(fp)
    skr_path = config.get("skr_path", "./skr")
    url_configs = config.get("url_configs", {})
    # 生成 skr 命令行参数
    args = [skr_path, "douyin", "--silent"]
    for url in url_configs:
        args.append("-u")
        args.append(url)
    # 调用 skr 获取数据
    ret = subprocess.check_output(args)
    users = json.loads(ret)
    # 历史数据
    user_histories = {}
    new_user_histories = {}
    if os.path.exists(user_histories_path):
        with open(user_histories_path, "rt", encoding=ENCODING) as fp:
            user_histories = json.load(fp)
    # 使用钉钉机器人发送消息
    session = requests.Session()
    for user in users:
        print("开始处理 User[NickName:{user[NickName]}, ID:{user[ID]}]".format(user=user))
        user["UpdatedAt"] = time.time()
        new_user_histories[user["ID"]] = user

        web_hook_urls = url_configs.get(user["URL"])
        if not web_hook_urls:
            continue
        for url in web_hook_urls:
            res = session.post(
                url, json=render_msg(user, user_histories.get(user["ID"]))
            )
            print("消息发送结束:", url, ",", res.json())

    with open(user_histories_path, "wt", encoding=ENCODING) as fp:
        json.dump(new_user_histories, fp, ensure_ascii=False, indent=2)
    print("新的历史纪录导出成功")


if __name__ == "__main__":
    # 配置文件路径
    config_path = "config.json"
    # 历史数据路径
    user_histories_path = "user_histories.json"
    main(config_path, user_histories_path)
