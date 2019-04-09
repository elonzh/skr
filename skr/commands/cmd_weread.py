import os

import click


@click.command()
@click.pass_context
@click.option("-d", "--duration", type=click.INT, default=6 * 3600)
@click.option(
    "-k",
    "--adbkey",
    type=click.Path(exists=True, dir_okay=False, resolve_path=True),
    default=os.path.expanduser("~/.android/adbkey"),
)
@click.option(
    "-p", "--port_path", type=click.STRING, help="The filename of usb port to use."
)
@click.option(
    "-s",
    "--serial",
    type=click.STRING,
    help="The serial number of the device to use."
    "If serial specifies a TCP address:port, "
    "then a TCP connection is used instead of a USB connection.",
)
def cli(ctx, duration, adbkey, port_path, serial):
    """
    微信读书自动读书
    """
    import random
    import re
    import time
    from datetime import datetime

    from adb import adb_commands, sign_m2crypto

    logger = ctx.obj.logger

    # KitKat+ devices require authentication
    signer = sign_m2crypto.M2CryptoSigner(adbkey)
    device = adb_commands.AdbCommands()
    device.ConnectDevice(
        port_path=port_path.encode() if port_path else None,
        serial=serial.encode() if serial else None,
        rsa_keys=[signer],
        default_timeout_ms=3000,
    )

    logger.info(
        "设备信息: %s-%s",
        device.Shell("getprop ro.product.brand").strip(),
        device.Shell("getprop ro.product.model").strip(),
    )

    pattern = re.compile(r"(\d+)x(\d+)")
    width, height = pattern.search(device.Shell("wm size")).groups()
    width = int(width)
    height = int(height)
    logger.info("屏幕尺寸: %dx%d", width, height)

    # 点亮屏幕
    device.Shell("input keyevent 224")
    # 关闭自动亮度
    device.Shell("settings put system screen_brightness_mode 0")
    # 将亮度调到最低
    device.Shell("settings put system screen_brightness 0")

    now = time.time()
    end_time = now + duration
    end_datetime = datetime.fromtimestamp(end_time)
    logger.info("截止时间: %s", end_datetime.isoformat())
    pages = 0
    while now < end_time:
        point = {
            "X": width * random.uniform(0.93, 0.96),
            "Y": width * random.uniform(0.93, 0.96),
        }
        device.Shell("input tap {X} {Y}".format(**point))
        pages += 1
        logger.info("进行翻页, 第 %d 页，坐标 %s", pages, point)
        time.sleep(random.randint(30, 40))
        now = time.time()
    logger.info("到达截止时间:%s, 停止运行", end_datetime.isoformat())
