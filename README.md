#PictureService

- cd PictureServic && go build
- 在配置文件conf/service.json里设置运营配置访问地址
- 访问运营地址，添加相应的app并配置参数
- ./PictureServic --web-debug  现在必须这样启动使用静态文件，下一步会把静态文件加载到生成的二进制文件中

app访问接口：

- 上传认证

POST : /1/upload/auth
body {
"appname":"videosns", //在运营界面添加的app名称，必填
"file_type":"image", //必填，值为image 或video
"key":"", //选填，可自己定义文件名
"file_hash":"", //选填，文件哈希值
}

RESP:
```
{
  "data": {
    "cloud": "qiniuyun",
    "key": "02-25-2016-f20a50b2-f2b6-4f4c-9251-65d15429714c",
    "token": "o0p03FklQWXZ5n4uX85D9CAwUQadX_ZyRI0Tm-51:cDWZ9PkRZTupsOm2Sl2PbptLBFU=:eyJzY29wZSI6InZpZGVvc25zOjAyLTI1LTIwMTYtZjIwYTUwYjItZjJiNi00ZjRjLTkyNTEtNjVkMTU0Mjk3MTRjIiwiZGVhZGxpbmUiOjE0NTYzODA2MzJ9",
    "url": "http://beijing6.appdao.com:5551/1/image/videosns/02-25-2016-f20a50b2--4f4c-9251-65d15429714c" 
  },
  "status": true
}
```

- 访问图片
	
	GET: 图片的地址+"?"+相应图片格式参数;
		
	七牛支持的参数为crop,width,size;
	
	UPYun需要在upyun控制台里设置自定义版本,使用版本名作为参数值,size=xxx