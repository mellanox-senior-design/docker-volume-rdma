#! /usr/bin/python

import os
import logging
import json

def main():
    logging.info("Starting result save.")
    result = {}
    with open("result.txt", "r") as res:
        result_arr = filter(lambda e: e != '' ,res.readline().replace("\n", "").split(" "))
        result = {
            "hostname": result_arr[0],
            "size":     result_arr[1],
            "sequential_output": {
                "per_chr": {
                    "kbs": result_arr[2],
                    "cp" : result_arr[3]
                },
                "block": {
                    "kbs": result_arr[4],
                    "cp" : result_arr[5]
                },
                "rewrite": {
                    "kbs": result_arr[6],
                    "cp" : result_arr[7]
                }
            },
            "sequential_input": {
                "per_chr": {
                    "kbs": result_arr[8],
                    "cp" : result_arr[9]
                },
                "block": {
                    "kbs": result_arr[10],
                    "cp" : result_arr[11]
                }
            },
            "random": {
                "seeks": {
                    "kbs": result_arr[12],
                    "cp" : result_arr[13]
                }
            }
        }

    with open('/tmp/bench_results/result.json', 'w') as fp:
        results = {
            "hostname": hostname,
            "results": result
        }

        logging.info(json.dumps(results))
        json.dump(results, fp)

if __name__ == '__main__':
    hostname = os.uname()[1]
    logging.basicConfig(format=hostname + ' %(asctime)s %(levelname)s: %(message)s', level=logging.DEBUG)
    main()
