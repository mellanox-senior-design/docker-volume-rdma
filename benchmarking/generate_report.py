#! /usr/bin/python

import json
import re
import glob

results = {}
for filename in glob.glob("**/bench_results*.json"):
    test_name = re.search('(.*)/bench_results.*\.json', filename).group(1)
    name = re.search('bench_results.?(.*)\.json', filename).group(1)
    if name == "":
        name = "none"

    with open(filename,"r") as fp:
        # print filename, name
        res = json.load(fp)
        res["name"] = name
        res["filename"] = filename
        if test_name in results.keys():
            results[test_name].append(res)
        else:
            results[test_name] = [res]
        
print json.dumps(results)
