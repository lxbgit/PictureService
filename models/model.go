package models

import "fmt"

type User struct {
	Key        string `xorm:"key TEXT PK " json:"key"`
	PassCode   string `xorm:"pass_code TEXT " json:"pass_code"`
	Name       string `xorm:"name TEXT  UNIQUE" json:"name"`
	CreatorKey string `xorm:"creator_key TEXT " json:"creator_key"`
	CreatedUTC int    `xorm:"created_utc INT " json:"created_utc"`
	AuxInfo    string `xorm:"aux_info TEXT" json:"aux_info"`

	CreatorName string `xorm:"-" json:"creator_name"`
}

func (*User) TableName() string {
	return "user"
}

func (m *User) UniqueCond() (string, []interface{}) {
	return "key=?", []interface{}{m.Key}
}

func GetAllUser(s *Session) ([]*User, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	var res []*User
	if err := s.Find(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func GetUsers(s *Session, page, count int) ([]*User, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	var res []*User
	if err := s.OrderBy("name desc").Limit(count, (page-1)*count).Find(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func GetUserCount(s *Session) (int, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	count, err := s.Count(&User{})

	return int(count), err
}

type App struct {
	Key        string `xorm:"key TEXT PK " json:"key"`
	Name       string `xorm:"name TEXT not NULL" json:"name"`
	CloudName  string `xorm:"cloudname TEXT not NULL" json:"cloudname"`
	Domain     string `xorm:"domain TEXT" json:"domain"`
	Bucket     string `xorm:"bucket TEXT" json:"bucket"`
	DateFormat string `xorm:"data_sign TEXT " json:"dateformat"`
	AuxInfo    string `xorm:"aux_info TEXT" json:"aux_info"`
	UserKey    string `xorm:"user_key TEXT " json:"creator_key"`

	UserName   string `xorm:"-" json:"creator_name"`
	CreatedUTC int    `xorm:"created_utc INT " json:"created_utc"`
}

func (*App) TableName() string {
	return "app"
}

func (m *App) UniqueCond() (string, []interface{}) {
	return "key=?", []interface{}{m.Key}
}

func GetAllApps(s *Session) ([]*App, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	var res []*App
	if err := s.Find(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func GetAppsByUserKey(s *Session, userKey string) ([]*App, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	var res []*App
	if err := s.Where("user_key=?", userKey).OrderBy("created_utc desc").Find(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func GetAllAppsPage(s *Session, page int, count int) ([]*App, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	var res []*App
	if err := s.OrderBy("created_utc desc").Limit(count, (page-1)*count).Find(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func GetAppCount(s *Session) (int, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	count, err := s.Count(&App{})

	return int(count), err
}

func SearchAppByName(s *Session, q string) ([]*App, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	var res []*App
	err := s.Where("like(?, name)=1", "%"+q+"%").OrderBy("").Find(&res)

	return res, err
}

func GetAppByName(s *Session, q string) (*App, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}
	var res []*App
	if err := s.Where("name=?", q).Find(&res); err != nil {
		fmt.Println(q, err)
		return nil, err
	}

	return res[0], nil
}
