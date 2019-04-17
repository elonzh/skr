import click
import logging
import os
import requests
import xml.etree.ElementTree as ET
from multiprocessing.dummy import Pool


class OctocatsDownloader:
    feeds_url = "http://feeds.feedburner.com/Octocats"

    def __init__(self, output="octocats", max_threads=5, force=False, logger=None):
        if not os.path.exists(output):
            os.mkdir(output)
        self.output = output

        self.session = requests.Session()
        self.pool = Pool(max_threads)
        self.force = force
        self.skip_count = 0
        self.update_count = 0
        self.feeds = None

        self.logger = logger or logging.getLogger()

    def join_path(self, path):
        return os.path.join(self.output, path)

    def download_job(self, img_element):
        src = img_element.get("src")
        filename = src.rsplit("/", 1)[-1]
        path = self.join_path(filename)
        if not self.force and os.path.exists(path):
            self.skip_count += 1
            self.logger.info("%s already exists! skip downloading ...", filename)
            return
        img = self.session.get(src).content
        with click.open_file(path, "wb") as fp:
            fp.write(img)
            self.update_count += 1
            self.logger.info("%s successfully downloaded.", filename)

    def fetch_feeds(self):
        self.logger.info("fetching RSS feeds ...")
        response = self.session.get(self.feeds_url)
        with click.open_file(self.join_path("Octocats.xml"), "w") as fp:
            fp.write(response.text)
        self.feeds = ET.fromstring(response.text)
        return self.feeds

    def download(self):
        feeds = self.feeds or self.fetch_feeds()
        # http://www.w3school.com.cn/xml/xml_namespaces.asp
        img_elements = feeds.iterfind(
            ".//atom:entry/atom:content/atom:div/atom:a/atom:img",
            {"atom": "http://www.w3.org/2005/Atom"},
        )
        self.logger.info("dispatching download jobs ...")
        self.pool.map(self.download_job, img_elements)
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
    default="octocats",
    help="The directory to save images.",
)
@click.option(
    "-m",
    "--max-threads",
    type=click.IntRange(1, 10),
    default=5,
    help="Max number of thread pool to download image.",
)
@click.option("-p", "--proxy", type=str, help="HTTP Proxy")
@click.option(
    "-f", "--force", is_flag=True, help="Fore download images even they exists."
)
def cli(ctx, output, max_threads, proxy, force):
    """
    Download Octocats from https://octodex.github.com
    """
    o = OctocatsDownloader(output, max_threads, force, logger=ctx.obj.logger)
    if proxy:
        o.session.proxies = {"http": proxy}
    o.download()
