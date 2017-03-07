# using WordPress API
# https://python-wordpress-xmlrpc.readthedocs.io/en/latest/index.html
import os
import logging
import time
import sys
import threading

from wordpress_xmlrpc import Client, WordPressPost
from wordpress_xmlrpc.methods.posts import GetPosts, NewPost, GetPost
from wordpress_xmlrpc.methods.users import GetUserInfo
from wordpress_xmlrpc.methods import posts

wp_url = "http://localhost:8000/xmlrpc.php"
wp_username = "jc"
wp_password = "7$aYpIKvm1T16LZCxF"
wp_blogid = ""

wp = Client(wp_url, wp_username, wp_password)

body = "But I must explain to you how all this mistaken idea of denouncing pleasure and praising pain was born and I will give you a complete account of the system, and expound the actual teachings of the great explorer of the truth, the master-builder of human happiness. No one rejects, dislikes, or avoids pleasure itself, because it is pleasure, but because those who do not know how to pursue pleasure rationally encounter consequences that are extremely painful. Nor again is there anyone who loves or pursues or desires to obtain pain of itself, because it is pain, but because occasionally circumstances occur in which toil and pain can procure him some great pleasure. To take a trivial example, which of us ever undertakes laborious physical exercise, except to obtain some advantage from it? But who has any right to find fault with a man who chooses to enjoy a pleasure that has no annoying consequences, or one who avoids a pain that produces no resultant pleasure?\n\nOn the other hand, we denounce with righteous indignation and dislike men who are so beguiled and demoralized by the charms of pleasure of the moment, so blinded by desire, that they cannot foresee the pain and trouble that are bound to ensue; and equal blame belongs to those who fail in their duty through weakness of will, which is the same as saying through shrinking from toil and pain. These cases are perfectly simple and easy to distinguish. In a free hour, when our power of choice is untrammelled and when nothing prevents our being able to do what we like best, every pleasure is to be welcomed and every pain avoided. But in certain circumstances and owing to the claims of duty or the obligations of business it will frequently occur that pleasures have to be repudiated and annoyances accepted. The wise man therefore always holds in these matters to this principle of selection: he rejects pleasures to secure other greater pleasures, or else he endures pains to avoid worse pains."

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
    post.user = 'mario'

    post.id = wp.call(NewPost(post))
    return post

def getPost(post):
    GetPost(post.id)

readerThreads = []
def loadPostTest(post, readers):
    logging.debug("Making " + str(readers) + " readers to get popular post...")
    sys.stdout.flush()

    start = time.time()
    for i in readers:
        id = i
        t = threading.Thread(target=getPost, args= (post, ))
        t.daemon = True
        readerThreads.append(t)

    for x in readerThreads:
        x.start()

    for x in readerThreads:
        x.join()

    end = time.time()
    logging.debug("Getting the popular post for " + str(readers) + " readers took: " + str(end-start))

    return end

def main():
    logging.debug("Starting...")

    argv = sys.argv
    if len(argv) != 3:
        sys.exit(1)

    post = 0
    if argv[1] == '1':
        logging.debug("Making new post...")
        post = makePost(
        'A post about how popular something cool is',
        body,
        {
        'post_tag': ['test', 'firstpost'],
        'category': ['Introductions', 'Tests']})
        logging.debug("Post was made...")

    else:
        logging.debug("Get the popular post...")
        offset = 0
        increment = 20
        while True:
            results = wp.call(posts.GetPosts({'number': increment, 'offset': offset}))
            if len(results) == 0:
                break
            for post in results:
                if(post.user == 'mario'):
                    post = loadPostTest(post, argv[2])

            offset = offset + increment

    if(post == 0):
        logging.debug("derp")
    else:
        logging.debug("Found the popular post...")

    loadPostTest(post, argv[2])

if __name__ == '__main__':
    hostname = os.uname()[1]
    logging.basicConfig(format=hostname + ' %(asctime)s %(levelname)s: %(message)s', level=logging.DEBUG)
    main()
