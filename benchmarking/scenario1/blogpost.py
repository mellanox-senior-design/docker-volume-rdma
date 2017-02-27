import datetime, xmlrpclib
from time import time

wp_url = "http://localhost:8000/xmlrpc.php"
wp_username = "jc"
wp_password = "7$aYpIKvm1T16LZCxF"
wp_blogid = ""

status_draft = 0
status_published = 1

server = xmlrpclib.ServerProxy(wp_url)

title = "Title with no blerg"
content = "blerg blerg blerg"
date_created = xmlrpclib.DateTime(datetime.datetime.strptime("2017-02-27 2:05", "%Y-%m-%d %H:%M"))
categories = ["somecategory"]
tags = ["sometag", "othertag"]
data = {'title': title, 'description': content, 'dateCreated': date_created, 'categories': categories, 'mt_keywords': tags}

post_id = server.metaWeblog.newPost(wp_blogid, wp_username, wp_password, data, status_published)

starttime = time()
ret = server.metaWeblog.getPost(post_id, wp_username, wp_password)
timetaken = time() - starttime
print timetaken
