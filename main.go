package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/PictureService/conf"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
	"qiniupkg.com/x/log.v7"
)

func main() {

	wd, _ := os.Getwd()
	pidFile, err := os.OpenFile(filepath.Join(wd, "picture.pid"), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("failed to create pid file: %s", err.Error())
		os.Exit(1)
	}
	pidFile.WriteString(strconv.Itoa(os.Getegid()))
	pidFile.Close()

	if conf.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ginIns := gin.New()
	ginIns.Use(gin.Recovery())
	if conf.DebugMode {
		ginIns.Use(gin.Logger())
	}

	if conf.WebDebugMode {
		// static
		ginIns.Static("/web", "./web")
	}
	//} else {
	//	// bin static
	//	ginIns.GET("/web/*file",
	//		func(c *gin.Context) {
	//			fileName := c.Param("file")
	//			if fileName == "/" {
	//				fileName = "/index.html"
	//			}
	//			data, err := Asset("web" + fileName)
	//			if err != nil {
	//				c.String(http.StatusNotFound, err.Error())
	//				return
	//			}
	//
	//			switch {
	//			case strings.LastIndex(fileName, ".html") == len(fileName)-5:
	//				c.Header("Content-Type", "text/html; charset=utf-8")
	//			case strings.LastIndex(fileName, ".css") == len(fileName)-4:
	//				c.Header("Content-Type", "text/css")
	//			}
	//			c.String(http.StatusOK, string(data))
	//		})
	//}

	//client api
	clientAPIGroup := ginIns.Group("/1")
	{
		clientAPIGroup.POST("/upload/auth", UploadAuth)
		clientAPIGroup.GET("/image/:appname/*key", RedirectThumbImage)
	}

	// op api
	opAPIGroup := ginIns.Group("/op")
	{
		opAPIGroup.POST("/login", Login)
		opAPIGroup.POST("/logout", OpAuth, Logout)

		opAPIGroup.GET("/users/:page/:count", InitUserCheck, OpAuth, GetUsers)
		opAPIGroup.POST("/user", OpAuth, NewUser)
		opAPIGroup.PUT("/user", OpAuth, UpdateUser)
		opAPIGroup.POST("/user/init", InitUser)
		opAPIGroup.GET("/user/info", OpAuth, GetLoginUserInfo)

		opAPIGroup.GET("/apps/user/:user_key", OpAuth, GetApps)
		opAPIGroup.GET("/apps/all/:page/:count", OpAuth, GetAllApps)

		opAPIGroup.POST("/app", OpAuth, NewApp)
		opAPIGroup.PUT("/app", OpAuth, UpdateApp)

	}

	err = gracehttp.Serve(&http.Server{Addr: fmt.Sprintf(":%d", conf.HttpPort), Handler: ginIns})
	if err != nil {
		log.Printf("fatal error: %s", err.Error())
	}

}
