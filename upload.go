package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"qiniupkg.com/api.v7/kodo"

	"github.com/PictureService/conf"
	"github.com/PictureService/logger"
	"github.com/PictureService/utils"
	"github.com/gin-gonic/gin"
)

func generateQiniuToken(bucket, key string, video bool) string {
	kodo.SetMac(conf.QiniuConfig.Key, conf.QiniuConfig.Secret)

	policy := &kodo.PutPolicy{}
	if bucket == "" {
		bucket = conf.QiniuConfig.Bucket
	}

	policy.Scope = bucket + ":" + key
	if video {
		policy.PersistentOps = "avthumb/m3u8/ab/128k/vb/640k/wmImage/aHR0cDovLzd4bDg4My5tZWRpYTEuejAuZ2xiLmNsb3VkZG4uY29tL2xvZ28ucG5n"
		policy.PersistentPipeline = conf.QiniuConfig.Pipeline
		policy.PersistentNotifyUrl = "http://" + conf.ServiceAddr + "/1/video/upload/qiniu/notify"
	}

	client := kodo.New(0, nil)
	return client.MakeUptoken(policy)
}

// make md5 from string
func md5Str(s string) (ret string) {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func generateUPYunToken(bucket, key string, video bool) (string, string, error) {
	policy := make(map[string]string)

	if bucket == "" {
		bucket = conf.UPYunConfig.Bucket
	}
	policy["bucket"] = bucket
	policy["save-key"] = key
	policy["expiration"] = strconv.FormatInt(time.Now().Unix()+int64(conf.UPYunConfig.Expiration), 10)

	args, err := json.Marshal(policy)
	if err != nil {
		return "", "", err
	}

	policyStr := base64.StdEncoding.EncodeToString(args)
	sig := md5Str(policyStr + "&" + conf.UPYunConfig.Secret)

	return policyStr, sig, nil
}

func UploadAuth(c *gin.Context) {
	var data struct {
		APPName  string `json:"appname" binding:"required"`
		FileType string `json:"file_type" binding:"required"`
		Key      string `json:"key"`
		FileHash string `json:"file_hash"`
	}
	err := c.BindJSON(&data)
	if err != nil {
		Error(c, BAD_POST_DATA, nil, err.Error())
		return
	}

	key := data.Key
	if key == "" {
		key = utils.GenerateUUID()
	}

	appinfo, ok := memConfAppsByName[data.APPName]
	if !ok {
		Error(c, BAD_POST_DATA, "app name not found")
		return
	}

	if data.Key == "" && appinfo.DateFormat != "" {
		t := time.Now()
		prefix := t.Format(appinfo.DateFormat)
		key = prefix + "-" + key
	}

	var policy, token string
	if appinfo.CloudName == "qiniu" {
		token = generateQiniuToken(appinfo.Bucket, key, data.FileType == "video")
	} else if appinfo.CloudName == "upyun" {
		policy, token, err = generateUPYunToken(appinfo.Bucket, key, data.FileType == "video")
		if err != nil {
			Error(c, SERVER_ERROR, nil, err.Error())
			return
		}
	}

	ret := map[string]interface{}{
		"cloud":  appinfo.CloudName,
		"key":    key,
		"token":  token,
		"policy": policy,
	}

	if data.FileType == "image" {
		ret["url"] = fmt.Sprintf("http://%s/1/image/%s/%s", conf.ServiceAddr, data.APPName, key)
	}

	logger.RequestLogger.Info(map[string]interface{}{
		"type":     "upload auth",
		"url":      c.Request.URL.Path,
		"request":  data,
		"response": ret,
	})

	Success(c, ret)
}

func RedirectThumbImage(c *gin.Context) {
	key := c.Param("key")
	width := c.Query("width")
	crop := c.Query("crop")
	size := c.Query("size")
	appname := c.Param("appname")
	if size != "" && len(size) > 2 && width == "" {
		width = size[1:]
	}

	appinfo, ok := memConfAppsByName[appname]
	if !ok {
		Error(c, BAD_POST_DATA, "app name not found")
		return
	}

	var location string
	if appinfo.Domain != "" {
		location = "http://" + appinfo.Domain + key
	} else if appinfo.CloudName == "qiniu" {
		location = "http://" + conf.QiniuConfig.Domain + key
	} else if appinfo.CloudName == "upyun" {
		location = "http://" + conf.UPYunConfig.Domain + key
	}

	if location == "" {
		Error(c, SERVER_ERROR, "location set error")
		return
	}

	if appinfo.CloudName == "qiniuyun" {
		if width != "" {
			if _, err := strconv.Atoi(width); err == nil {
				if crop == "1" {
					location = location + "?imageView2/1/w/" + width
				} else {
					location = location + "?imageView2/2/w/" + width
				}
			} else {
				logger.ErrorLogger.Error(map[string]interface{}{
					"type":    "image_thumb",
					"url":     c.Request.URL.Path,
					"err_msg": err.Error(),
				})
			}
		}
	}

	if appinfo.CloudName == "upyun" && size != "" {
		location = location + "!" + size
	}

	c.Redirect(302, location)
}
