# creates a 1 GB file and writes to file. Use this function to verify 1 GB file is created.
import os
import logging
import time
import sys
import threading
import random

def main():
    f = open('words')
    words = f.read().strip('\n').split('\n')
    f.close()

    numWords = len(words) - 2
    line = []
    bodyList = []

    for x in range(0, 10000000):
        for y in range(0, 20):
            i = random.randint(0, numWords)
            line.append(words[i])
        line.append("\n")
        bodyList.append(' '.join(line))
        line = []

    body = ''.join(bodyList).strip('\n')

    f = open('file.txt', 'w')
    f.write(body)
    f.close()


if __name__ == '__main__':
    # hostname = os.uname()[1]
    # logging.basicConfig(format=hostname + ' %(asctime)s %(levelname)s: %(message)s', level=logging.DEBUG)
    main()
