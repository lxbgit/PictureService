#!/usr/bin/env python
# -*- coding: utf-8 -*-

import urllib
import urllib2
import json
import upyun
import requests


if __name__ == '__main__':
    url = "http://127.0.0.1:8080/1/upload/auth"

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
    policy = decoded['data']['policy']
    print 'start upload'



    up_url = "http://v0.api.upyun.com/appwill123"
    data = {
        "file":"/Users/tomgreen/Desktop/aaaa.jpg",
        "signature":token,
        "policy":policy,
    }

    session = requests.Session()
    resp = session.request('POST', up_url, files=data)

    print(resp.content)
