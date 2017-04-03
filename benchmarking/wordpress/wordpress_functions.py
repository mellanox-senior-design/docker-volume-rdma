# using WordPress API
# https://python-wordpress-xmlrpc.readthedocs.io/en/latest/index.html
import os
import logging
import time
import sys
import threading
import random

from wordpress_xmlrpc import Client, WordPressPost
from wordpress_xmlrpc.methods.posts import GetPosts, NewPost, GetPost
from wordpress_xmlrpc.methods.users import GetUserInfo
from wordpress_xmlrpc.methods import posts

wp_url = "http://wordpress/xmlrpc.php"
wp_username = "test_account"
wp_password = "wordpress"
wp_blogid = ""

wp = Client(wp_url, wp_username, wp_password)

body = ""

def getPosts():
    logging.debug("Get Posts...")

    # posts = wp.call(GetPosts())
    logging.debug(posts)

    return posts

def makePost(title, content, terms_names):
    post = WordPressPost()
    post.title = title
    post.content = content
    post.post_status = 'publish'
    post.terms_names = terms_names
    post.user = 'Mario'

    post.id = wp.call(NewPost(post))
    return post

def getPost(post):
    GetPost(post.id)

def main():
    logging.debug("Starting...")

    argv = sys.argv
    if len(argv) != 2:
        sys.exit(1)

    post = 0

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

    body = ''.join(bodyList).strip('\n').strip()

    if argv[1] == 'post':
        logging.debug("Making new post...")
        post = makePost(
        'A post about how popular something cool is',
        body,
        {
        'post_tag': ['test', 'firstpost'],
        'category': ['Introductions', 'Tests']})
        logging.debug("Post was made...")
        print post.id



if __name__ == '__main__':
    hostname = os.uname()[1]
    logging.basicConfig(format=hostname + ' %(asctime)s %(levelname)s: %(message)s', level=logging.DEBUG)
    main()
