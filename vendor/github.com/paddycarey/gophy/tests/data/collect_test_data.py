#!/usr/bin/env python
"""collect_test_data.py

reads a list of urls from a text file (one per line), makes a HTTP request to
each URL, saving the result in the `tests/data/` directory for use as mock
responses during test runs of the gophy library.
"""

# future imports
from __future__ import absolute_import
from __future__ import print_function

# stdlib imports
import base64
import json
import os

# third-party imports
import requests


def main():

    # read list of urls from a text file
    with open("urls.txt", "r") as f:
        urls = f.read()
    urls = urls.splitlines()

    # make a simple get request to each url in turn and write the returned data
    # to a json file on disk, using a base64 representation of the url as a
    # filename.
    for url in urls:
        filename = base64.urlsafe_b64encode(url) + ".json"
        if os.path.exists(filename):
            continue
        response = requests.get("https://api.giphy.com/v1/" + url)
        response_data = {
            "status": response.status_code,
            "headers": dict(response.headers.items()),
            "body": response.text,
        }
        with open(filename, "w") as f:
            f.write(json.dumps(response_data))


if __name__ == "__main__":
    main()
