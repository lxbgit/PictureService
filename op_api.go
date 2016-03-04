package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/PictureService/conf"
	"github.com/PictureService/models"
	"github.com/PictureService/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var data struct {
		Name     string `json:"name" binding:"required"`
		PassCode string `json:"pass_code" binding:"required"`
	}
	if err := c.BindJSON(&data); err != nil {
		Error(c, BAD_POST_DATA, err.Error())
		return
	}

	memConfMux.RLock()
	user := memConfUsersByName[data.Name]
	memConfMux.RUnlock()

	if user == nil {
		Error(c, USER_NOT_EXIST)
		return
	}
	if user.PassCode != encryptUserPassCode(data.PassCode) {
		Error(c, PASS_CODE_ERROR)
		return
	}

	setUserKeyCookie(c, user.Key)
	Success(c, nil)
}

func Logout(c *gin.Context) {
	deleteUserKeyCookie(c)
	Success(c, nil)
}

type newUserData struct {
	Name     string `json:"name" binding:"required"`
	PassCode string `json:"pass_code" binding:"required"`
	AuxInfo  string `json:"aux_info"`
}

func InitUser(c *gin.Context) {
	confWriteMux.Lock()
	defer confWriteMux.Unlock()

	memConfMux.RLock()
	if len(memConfUsers) > 0 {
		memConfMux.RUnlock()
		Error(c, BAD_REQUEST, "users already exists")
		return
	}
	memConfMux.RUnlock()

	data := &newUserData{}
	if err := c.BindJSON(data); err != nil {
		Error(c, BAD_POST_DATA, err.Error())
		return
	}

	if err := verifyNewUserData(data); err != nil {
		Error(c, BAD_REQUEST, err.Error())
		return
	}

	key := utils.GenerateKey()
	user, err := newUserWithNewUserData(data, key, key)
	if err != nil {
		Error(c, SERVER_ERROR, err.Error())
		return
	}
	setUserKeyCookie(c, user.Key)

	Success(c, nil)

}

func NewUser(c *gin.Context) {
	confWriteMux.Lock()
	defer confWriteMux.Unlock()

	data := &newUserData{}
	if err := c.BindJSON(data); err != nil {
		Error(c, BAD_POST_DATA, err.Error())
		return
	}

	if err := verifyNewUserData(data); err != nil {
		Error(c, BAD_REQUEST, err.Error())
		return
	}

	_, err := newUserWithNewUserData(data, utils.GenerateKey(), getOpUserKey(c))
	if err != nil {
		Error(c, SERVER_ERROR, err.Error())
		return
	}

	Success(c, nil)

}

func verifyNewUserData(data *newUserData) error {
	memConfMux.RLock()
	defer memConfMux.RUnlock()

	if memConfUsersByName[data.Name] != nil {
		return fmt.Errorf("user [%s] already exists", data.Name)
	}

	if len(data.Name) < 3 {
		return fmt.Errorf("user name too short, length must bigger than 2")
	}

	if len(data.PassCode) < 6 {
		return fmt.Errorf("user passcode too short, length must bigger than 6")
	}

	return nil
}

func newUserWithNewUserData(data *newUserData, userKey, creatorKey string) (*models.User, error) {
	user := &models.User{
		Name:       data.Name,
		PassCode:   encryptUserPassCode(data.PassCode),
		CreatorKey: creatorKey,
		CreatedUTC: utils.GetNowSecond(),
		AuxInfo:    data.AuxInfo,
		Key:        userKey}

	return updateUser(user)
}

type updateUserData struct {
	Name    string `json:"name" binding:"required"`
	AuxInfo string `json:"aux_info"`
}

func UpdateUser(c *gin.Context) {
	confWriteMux.Lock()
	defer confWriteMux.Unlock()

	data := &updateUserData{}
	if err := c.BindJSON(data); err != nil {
		Error(c, BAD_POST_DATA, err.Error())
		return
	}

	if err := verifyUpdateUserData(data, getOpUserKey(c)); err != nil {
		Error(c, BAD_REQUEST, err.Error())
		return
	}

	_, err := updateUserWithUpdateData(data, getOpUserKey(c))
	if err != nil {
		Error(c, SERVER_ERROR, err.Error())
		return
	}

	Success(c, nil)
}

func verifyUpdateUserData(data *updateUserData, userKey string) error {
	memConfMux.RLock()
	defer memConfMux.RUnlock()

	if memConfUsersByName[data.Name] != nil && memConfUsersByName[data.Name].Key != memConfUsers[userKey].Key {
		return fmt.Errorf("user name [%s] already exists", data.Name)
	}

	if len(data.Name) < 3 {
		return fmt.Errorf("user name too short, length must bigger than 2")
	}

	return nil
}

func updateUserWithUpdateData(data *updateUserData, userKey string) (*models.User, error) {
	memConfMux.RLock()
	user := *memConfUsers[userKey]
	memConfMux.RUnlock()

	user.AuxInfo = data.AuxInfo
	user.Name = data.Name
	return updateUser(&user)
}

func updateUser(user *models.User) (*models.User, error) {
	s := models.NewSession()
	defer s.Close()
	if err := s.Begin(); err != nil {
		s.Rollback()
		return nil, err
	}

	memConfMux.RLock()
	oldUser := memConfUsers[user.Key]
	memConfMux.RUnlock()

	if oldUser == nil {
		if err := models.InsertRow(s, user); err != nil {
			s.Rollback()
			fmt.Println(err)
			return nil, err
		}
	} else {
		if err := models.UpdateDBModel(s, user); err != nil {
			s.Rollback()
			return nil, err
		}
	}

	if err := s.Commit(); err != nil {
		s.Rollback()
		return nil, err
	}

	memConfMux.Lock()
	defer memConfMux.Unlock()

	if oldUser != nil {
		memConfUsersByName[oldUser.Name] = nil
	}
	memConfUsers[user.Key] = user
	memConfUsersByName[user.Name] = user

	return user, nil
}

func GetUsers(c *gin.Context) {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		Error(c, BAD_REQUEST, "page not number")
		return
	}
	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		Error(c, BAD_REQUEST, "count not number")
		return
	}

	users, err := models.GetUsers(nil, page, count)
	if err != nil {
		Error(c, SERVER_ERROR, err.Error())
		return
	}

	totalCount, err := models.GetUserCount(nil)
	if err != nil {
		Error(c, SERVER_ERROR, err.Error())
		return
	}

	memConfMux.RLock()
	for _, user := range users {
		user.PassCode = ""
		if memConfUsers[user.CreatorKey] != nil {
			user.CreatorName = memConfUsers[user.CreatorKey].Name
		}
	}
	memConfMux.RUnlock()

	Success(c, map[string]interface{}{
		"total_count": totalCount,
		"list":        users,
	})
}

func GetApp(c *gin.Context) {
	memConfMux.RLock()
	app := memConfApps[c.Param("app_key")]
	memConfMux.RUnlock()

	if app == nil {
		Success(c, nil)
		return
	}

	returnApp := *app

	Success(c, &returnApp)
}

func SearchApps(c *gin.Context) {
	apps, err := searchApps(c.Query("q"))
	if err != nil {
		Error(c, SERVER_ERROR, err.Error())
		return
	}

	Success(c, apps)
}

func searchApps(q string) ([]*models.App, error) {
	return models.SearchAppByName(nil, q)
}

func NewApp(c *gin.Context) {
	confWriteMux.Lock()
	defer confWriteMux.Unlock()

	var data struct {
		Name       string `json:"name" binding:"required"`
		CloudName  string `json:"cloudname" binding:"required"`
		DateFormat string `json:"dateformat"`
		Domain     string `json:"domain"`
		Bucket     string `json:"bucket"`
		AuxInfo    string `json:"aux_info"`
	}
	if err := c.BindJSON(&data); err != nil {
		Error(c, BAD_POST_DATA, err.Error())
		return
	}

	memConfMux.RLock()
	if memConfAppsByName[data.Name] != nil {
		memConfMux.RUnlock()
		Error(c, BAD_REQUEST, "appname already exists: "+data.Name)
		return
	}
	memConfMux.RUnlock()

	app := &models.App{
		Key:        utils.GenerateKey(),
		Name:       data.Name,
		CloudName:  data.CloudName,
		Domain:     data.Domain,
		Bucket:     data.Bucket,
		DateFormat: data.DateFormat,
		AuxInfo:    data.AuxInfo,
		CreatedUTC: utils.GetNowSecond(),
		UserKey:    getOpUserKey(c),
	}

	if _, err := updateApp(app); err != nil {
		Error(c, SERVER_ERROR, err.Error())
		return
	}

	Success(c, nil)

}

func updateApp(app *models.App) (*models.App, error) {
	s := models.NewSession()
	defer s.Close()
	if err := s.Begin(); err != nil {
		s.Rollback()
		return nil, err
	}

	memConfMux.RLock()
	oldApp := memConfApps[app.Key]
	memConfMux.RUnlock()

	if oldApp == nil {
		if err := models.InsertRow(s, app); err != nil {
			s.Rollback()
			return nil, err
		}
	} else {
		if err := models.UpdateDBModel(s, app); err != nil {
			s.Rollback()
			return nil, err
		}
	}

	if err := s.Commit(); err != nil {
		s.Rollback()
		return nil, err
	}

	memConfMux.Lock()
	defer memConfMux.Unlock()

	if oldApp != nil {
		memConfAppsByName[oldApp.Name] = nil
	}
	memConfApps[app.Key] = app
	memConfAppsByName[app.Name] = app

	return app, nil
}

func UpdateApp(c *gin.Context) {
	confWriteMux.Lock()
	defer confWriteMux.Unlock()

	var data struct {
		Key        string `json:"key" binding:"required"`
		Name       string `json:"name" binding:"required"`
		CloudName  string `json:"cloudname" binding:"required"`
		Domain     string `json:"domain"`
		Bucket     string `json:"bucket"`
		DateFormat string `json:"dateformat"`
		AuxInfo    string `json:"aux_info"`
	}
	if err := c.BindJSON(&data); err != nil {
		Error(c, BAD_POST_DATA, err.Error())
		return
	}
	memConfMux.RLock()
	oldaapp := memConfApps[data.Key]
	if oldaapp == nil {
		Error(c, BAD_REQUEST, "app key not exists: "+data.Key)
		memConfMux.Unlock()
		return
	}

	if memConfApps[data.Key].Name != data.Name {
		for _, app := range memConfApps {
			if app.Name == data.Name {
				Error(c, BAD_REQUEST, "appname already exists: "+data.Name)
				memConfMux.Unlock()
				return
			}
		}
	}

	memConfMux.RUnlock()

	app := *oldaapp
	app.Name = data.Name
	app.CloudName = data.CloudName
	app.Domain = data.Domain
	app.Bucket = data.Bucket
	app.DateFormat = data.DateFormat
	app.AuxInfo = data.AuxInfo
	if _, err := updateApp(&app); err != nil {
		Error(c, SERVER_ERROR, err.Error())
		return
	}

	Success(c, nil)
}

func GetApps(c *gin.Context) {
	userKey := c.Param("user_key")
	apps, err := models.GetAppsByUserKey(nil, userKey)
	if err != nil {
		Error(c, SERVER_ERROR, err.Error())
		return
	}

	for _, app := range apps {
		app.UserName = memConfUsers[app.UserKey].Name
	}

	Success(c, apps)
}

func GetAllApps(c *gin.Context) {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		Error(c, BAD_REQUEST, "page not number")
		return
	}
	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		Error(c, BAD_REQUEST, "count not number")
		return
	}

	apps, err := models.GetAllAppsPage(nil, page, count)
	if err != nil {
		Error(c, SERVER_ERROR, err.Error())
		return
	}

	totalCount, err := models.GetAppCount(nil)
	if err != nil {
		Error(c, SERVER_ERROR, err.Error())
		return
	}

	Success(c, map[string]interface{}{
		"total_count": totalCount,
		"list":        apps,
	})
}

func OpAuth(c *gin.Context) {
	cookie, err := c.Request.Cookie("op_user")
	if err != nil {
		Error(c, NOT_LOGIN, err.Error())
		c.Abort()
		return
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.UserPassCodeEncryptKey), nil
	})
	if err != nil {
		Error(c, NOT_LOGIN, err.Error())
		c.Abort()
		return
	}

	if !token.Valid {
		Error(c, NOT_LOGIN, "cookie token invalid")
		c.Abort()
		return
	}

	userKey := token.Claims["uky"].(string)
	memConfMux.RLock()
	if memConfUsers[userKey] == nil {
		memConfMux.RUnlock()
		Error(c, NOT_LOGIN, "user not exist")
		c.Abort()
		return
	}
	memConfMux.RUnlock()

	setOpUserKey(c, userKey)
}

func InitUserCheck(c *gin.Context) {
	memConfMux.RLock()
	userCount := len(memConfUsers)
	memConfMux.RUnlock()

	if userCount == 0 {
		Error(c, USER_NOT_INIT)
		c.Abort()
	}
}

func GetLoginUserInfo(c *gin.Context) {
	key := getOpUserKey(c)
	memConfMux.RLock()
	user := *memConfUsers[key]
	user.CreatorName = memConfUsers[user.CreatorKey].Name
	memConfMux.RUnlock()

	user.PassCode = ""

	Success(c, user)
}

func encryptUserPassCode(code string) string {
	hash := hmac.New(sha256.New, []byte(conf.UserPassCodeEncryptKey))
	hash.Write([]byte(code))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func setUserKeyCookie(c *gin.Context, userKey string) {
	jwtIns := jwt.New(jwt.SigningMethodHS256)
	jwtIns.Claims["uky"] = userKey

	encStr, _ := jwtIns.SignedString([]byte(conf.UserPassCodeEncryptKey))
	cookie := new(http.Cookie)
	cookie.Name = "op_user"
	cookie.Expires = time.Now().Add(time.Duration(30*86400) * time.Second)
	cookie.Value = encStr
	cookie.Path = "/op"
	http.SetCookie(c.Writer, cookie)
}

func deleteUserKeyCookie(c *gin.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "op_user"
	cookie.Value = ""
	http.SetCookie(c.Writer, cookie)
}
