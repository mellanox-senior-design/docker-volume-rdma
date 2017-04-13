# loads page using http load generator from https://github.com/rakyll/hey
# also outputs in csv format the return results from the generator

# layout of csv file is:
# summary
# first five are field, time (in secs)
# last two are field, amt (in bytes)

# response time historgram
# time (in secs), number of occurences

# latency distribution
# percent in this time, time (in secs)

import os
import logging
import time
import sys
import subprocess
import re
import json

def main():
    logging.debug("Starting...")

    fps = 0
    s = time.time()
    data = subprocess.check_output(["sh", "test.sh"])
    e = time.time()
    test_duration = e - s
    lines = data.split("\n")
    for i in lines:
        l=i[11:]
        if l[0:7] == "work: a":
            fps = re.search("([0-9.]+)", l).groups()[0]

    logging.info("Starting result save.")
    with open('/tmp/bench_results/result.json', 'w') as fp:
        results = {
            "hostname": hostname,
            "results": {
                "bbb": {
                    "fps": fps,
                    "time": test_duration
                }
            }
        }

        logging.info(json.dumps(results))
        json.dump(results, fp)

if __name__ == '__main__':
    hostname = os.uname()[1]
    logging.basicConfig(format=hostname + ' %(asctime)s %(levelname)s: %(message)s', level=logging.DEBUG)
    main()
