#!/usr/bin/env python
# -*- coding: utf-8 -*-

import urllib
import urllib2
import json
from qiniu import Auth
from qiniu import put_file
from qiniu import put_data
from qiniu import etag



if __name__ == '__main__':
    url = "http://182.92.64.58:3000/1/upload/auth"
    # url = "http://127.0.0.1:8080/1/upload/auth"

    data = {
        "appname":"videosns",
        "file_type":"image",
        "key":"",
        "file_hash":"",
    }
    jdata = json.dumps(data)
    req = urllib2.Request(url,jdata)
    response = urllib2.urlopen(req)

    qiniudata =  response.read()

    print(qiniudata)
    decoded = json.loads(qiniudata)
    json.dumps(decoded)
    key = decoded['data']['key']
    print 'key:', key
    token =  decoded['data']['token']

    print 'start upload'

    mime_type = "image/png"
    params = {'x:a': 'a'}
    localfile = "/Users/tomgreen/Desktop/test.jpeg"
    progress_handler = lambda progress, total: progress
    ret, info = put_file(token, key, localfile, params, mime_type, progress_handler=progress_handler)
    assert ret['key'] == key


    print(ret)
    print(info)