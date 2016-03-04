package main

import (
	"log"
	"sync"

	"github.com/lxbgit/PictureService/models"
)

var (
	confWriteMux       = sync.Mutex{}
	memConfUsers       map[string]*models.User
	memConfUsersByName map[string]*models.User
	memConfAppsByName  map[string]*models.App
	memConfApps        map[string]*models.App

	memConfMux = sync.RWMutex{}
)

func init() {
	users, err := models.GetAllUser(nil)
	if err != nil {
		log.Panicf("Failed to load user info: %s", err.Error())
	}

	apps, err := models.GetAllApps(nil)
	if err != nil {
		log.Panicf("Failed to load app info: %s", err.Error())
	}

	fillMemConfData(users, apps)

}

func fillMemConfData(users []*models.User, apps []*models.App) {

	memConfUsers = make(map[string]*models.User)
	memConfUsersByName = make(map[string]*models.User)
	memConfApps = make(map[string]*models.App)
	memConfAppsByName = make(map[string]*models.App)

	for _, user := range users {
		memConfUsers[user.Key] = user
		memConfUsersByName[user.Name] = user
	}

	for _, app := range apps {
		memConfApps[app.Key] = app
		memConfAppsByName[app.Name] = app
	}

}
