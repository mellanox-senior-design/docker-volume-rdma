import os
import logging
import time
import sys
import subprocess

def main():
    results = subprocess.call(["hey", "-m=GET", "-disable-compression", "disable-keepalive", "https://google.com"])

    results = subprocess.call(["hey", "-m=GET", "-disable-compression", "disable-keepalive", "http://localhost:8000/?p=20"])

if __name__ == '__main__':
    hostname = os.uname()[1]
    logging.basicConfig(format=hostname + ' %(asctime)s %(levelname)s: %(message)s', level=logging.DEBUG)
    main()
