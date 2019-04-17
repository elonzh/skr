import click
import logging
import os
import requests
from bs4 import BeautifulSoup
from multiprocessing.dummy import Pool
from urllib.parse import urlparse


class SimpleDesktopsDownloader:
    def __init__(
        self,
        output="simpledesktops",
        max_threads=5,
        force=False,
        tree=False,
        logger=None,
    ):
        if not os.path.exists(output):
            os.mkdir(output)
        self.output = output
        self.pool = Pool(max_threads)
        self.force = force
        self.tree = tree

        self.session = requests.Session()
        self.skip_count = 0
        self.update_count = 0

        self.logger = logger or logging.getLogger()

    def join(self, path):
        return os.path.join(self.output, path)

    def download_job(self, img):
        # .295x184_q100.png
        img_url = img["src"][:-17]
        self.logger.info("start job: %s", img_url)
        # /uploads/desktops/
        path = urlparse(img_url).path[18:]
        if self.tree:
            dir = path.rsplit("/", 1)[0]
            os.makedirs(dir, exist_ok=True)
        else:
            path = path.replace("/", "-")

        path = self.join(path)
        if not self.force and os.path.exists(path):
            self.skip_count += 1
            self.logger.info("%s already exists! skip downloading ...", path)
            return
        with click.open_file(path, "wb") as fp:
            r = self.session.get(img_url)
            if not r.ok:
                self.logger.error("something wrong! [%d]%s", r.status_code, img_url)
                self.pool.terminate()
                return
            fp.write(r.content)
            self.update_count += 1
            self.logger.info("%s successfully downloaded.", path)

    def iter_download_job(self):
        url = "http://simpledesktops.com/browse/"
        page = 1
        while True:
            response = self.session.get(url + str(page))
            if response.status_code == 404:
                break
            page += 1
            bs = BeautifulSoup(response.text, "html.parser")
            for img in bs.select(".desktops .edge .desktop img"):
                yield img

    def download(self):
        self.logger.info("dispatching download jobs ...")
        self.pool.map(self.download_job, self.iter_download_job())
        self.logger.info(
            "all task done, %d updated, %d skipped, enjoy!",
            self.update_count,
            self.skip_count,
        )


@click.command()
@click.pass_context
@click.option(
    "-o",
    "--output",
    type=click.Path(file_okay=False, writable=True),
    default="simpledesktops",
    help="The directory to save files.",
)
@click.option(
    "-m",
    "--max-threads",
    type=click.IntRange(1, 20),
    default=5,
    help="Max number of thread pool to download file.",
)
@click.option(
    "-f", "--force", is_flag=True, help="Force download even file already exists."
)
@click.option("-t", "--tree", is_flag=True, help="Save files in to a directory tree.")
@click.option("-p", "--proxy", type=str, help="HTTP Proxy")
def cli(ctx, output, max_threads, force, tree, proxy):
    """
    Download wallpapers from http://simpledesktops.com/
    """
    s = SimpleDesktopsDownloader(
        output, max_threads, force, tree, logger=ctx.obj.logger
    )
    if proxy:
        s.session.proxies = {"http": proxy}
    s.download()
