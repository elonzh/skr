import importlib
import logging
import os
import sys
from collections import namedtuple

import click
import colorlog


def is_root():
    return os.getegid() == 0


def sudo():
    args = ["sudo", sys.executable] + sys.argv + [os.environ]
    # the next line replaces the currently-running process with the sudo
    os.execlpe("sudo", *args)


def root_required(func):
    if not is_root():
        sudo()

    return func


CONTEXT_SETTINGS = dict(auto_envvar_prefix="SKR")

cmd_folder = os.path.abspath(os.path.join(os.path.dirname(__file__), "commands"))


class CLI(click.MultiCommand):
    def list_commands(self, ctx):
        rv = []
        for filename in os.listdir(cmd_folder):
            if not filename.startswith("cmd_"):
                continue
            path = os.path.join(cmd_folder, filename)

            if os.path.isfile(path) and filename.endswith(".py"):
                rv.append(filename[4:-3])
            elif os.path.isdir(path):
                rv.append(filename[4:])
        rv.sort()
        return rv

    def get_command(self, ctx, name):
        mod = importlib.import_module(".".join(["commands", "cmd_" + name]))
        return mod.cli


@click.command(cls=CLI, context_settings=CONTEXT_SETTINGS)
@click.pass_context
@click.option("-v", "--verbose", is_flag=True, help="Enables verbose mode.")
def cli(ctx, verbose):
    """🏁  skr~ skr~"""
    logger = logging.Logger("skr", logging.INFO)
    handler = logging.StreamHandler()
    handler.setFormatter(
        colorlog.ColoredFormatter(
            "[%(log_color)s%(levelname)-8s%(reset)s]%(blue)s%(message)s",
            datefmt=None,
            reset=True,
            log_colors={
                "DEBUG": "cyan",
                "INFO": "green",
                "WARNING": "yellow",
                "ERROR": "red",
                "CRITICAL": "red,bg_white",
            },
            secondary_log_colors={},
            style="%",
        )
    )
    logger.addHandler(handler)

    if verbose:
        logger.setLevel(logging.DEBUG)

    ctx.obj = namedtuple("GlobalContext", ("logger",))(logger=logger)


if __name__ == "__main__":
    cli()
