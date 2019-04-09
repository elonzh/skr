import random
import sys
import time

import click

FRAMES = ["🌝", "🌖", "🌗", "🌘", "🌚", "🌒", "🌓", "🌔"]
CHOICES = ["🌝", "🌚"]


def write(text):
    sys.stdout.write("\r" + text)
    sys.stdout.flush()


@click.command()
@click.argument("number", type=click.IntRange(min=1, max=10), default=1)
def cli(number):
    """
    When faced with two choices, simply toss a coin.
    It works not because it settles the question for you,
    but because in that brief moment when the coin is in the air,
    you suddenly know what you are hoping for.

    see: https://www.v2ex.com/t/546324
    """
    coins = [""] * number
    loop = 3
    for n in range(number):
        for i in range(loop * len(FRAMES)):
            coins[n] = FRAMES[i % len(FRAMES)]
            write(" ".join(coins))
            time.sleep(1.0 / (loop * len(FRAMES)))
        coins[n] = random.choice(CHOICES)
        write(" ".join(coins))
