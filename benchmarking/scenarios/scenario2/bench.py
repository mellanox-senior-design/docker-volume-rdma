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
    # subprocess.call(["hey", "-m=GET", "-disable-compression", "https://google.com"])
    # ret = subprocess.check_output(["hey", "-m=GET", "-disable-compression", "http://localhost:8000/?p=20"])
    # hey -m=GET -disable-compression https://google.com
    logging.debug("Starting...")
    ret = subprocess.check_output(["hey", "-m=GET", "-disable-compression", "https://google.com"])

    ret = re.sub(r'(\n\n)', "\n", ret)
    res = []

    ret = ret.splitlines()
    it = ret.__iter__()
    while True:
        try:
            line = it.next()
            while(line.lower() != "All requests done.".lower()): line = it.next()

            # summary results
            line = it.next()
            # res.append(line.lower())

            line = it.next()
            while(line.lower() != "Status code distribution:".lower()):

                tmp = re.sub(r'(\s)+', "", line)
                tmp = re.sub(r'(secs)', "", tmp)
                tmp = re.sub(r'(bytes)', "", tmp)
                tmp = re.sub(r':', ",", tmp).lower()
                res.append(tmp)
                line = it.next()

            while(line.lower() != "Response time histogram:".lower()): line = it.next()

            # response time histogram
            line = it.next()
            while(line.lower() != "Latency distribution:".lower()):
                tmp = re.sub(r'^[\s\W]+', "", line)
                tmp = re.sub(r'[\s\W]+$', "", tmp)
                tmp = re.sub(r'(\s\[)+', ",", tmp).lower()
                res.append(tmp)
                line = it.next()

            # latency distribution
            line = it.next()
            while(line):
                tmp = re.sub(r'^[\s\W]+', "", line)
                tmp = re.sub(r'(\D)+$', "", tmp)
                tmp = re.sub(r'(%\sin\s)', ",", tmp).lower()
                res.append(tmp)
                line = it.next()
        except StopIteration:
            break

    logging.info("Starting result save.")
    with open('/tmp/bench_results/result.json', 'w') as fp:
        results = {
            "hostname": hostname,
            "results": {
                "Load Post": {
                    "Results": res
                }
            }
        }

        logging.info(json.dumps(results))
        json.dump(results, fp)

if __name__ == '__main__':
    hostname = os.uname()[1]
    logging.basicConfig(format=hostname + ' %(asctime)s %(levelname)s: %(message)s', level=logging.DEBUG)
    main()
