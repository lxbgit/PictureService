package models

import (
	//	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/lxbgit/PictureService/conf"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbEngineDefault *xorm.Engine
)

type Session struct {
	*xorm.Session
}

func init() {
	var err error
	var dsn, driver string

	dsn = fmt.Sprintf(filepath.Join(conf.SqliteDir, conf.SqliteFileName))

	driver = "sqlite3"

	if dbEngineDefault, err = xorm.NewEngine(driver, dsn); err != nil {
		log.Fatal("Failed to init db engine: " + err.Error())
	}
	dbEngineDefault.SetMaxOpenConns(100)
	dbEngineDefault.SetMaxIdleConns(50)
	dbEngineDefault.ShowSQL(conf.DebugMode)

	if err = dbEngineDefault.Sync2(
		&User{}, &App{},
	); err != nil {
		log.Panicf("Failed to sync db scheme: %s", err.Error())
	}

}

func NewSession() *Session {
	ms := new(Session)
	ms.Session = dbEngineDefault.NewSession()

	return ms
}

func newAutoCloseModelsSession() *Session {
	ms := new(Session)
	ms.Session = dbEngineDefault.NewSession()
	ms.IsAutoClose = true

	return ms
}

type DBModel interface {
	TableName() string
}

func InsertRow(s *Session, m DBModel) (err error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}
	_, err = s.AllCols().InsertOne(m)

	return
}

func InsertMultiRows(s *Session, m ...interface{}) (err error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}
	_, err = s.AllCols().Insert(m...)

	return
}

type UniqueDBModel interface {
	TableName() string
	UniqueCond() (string, []interface{})
}

func UpdateDBModel(s *Session, m UniqueDBModel) (err error) {
	whereStr, whereArgs := m.UniqueCond()
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	_, err = s.AllCols().Where(whereStr, whereArgs...).Update(m)

	return
}

func DeleteDBModel(s *Session, m UniqueDBModel) (err error) {
	whereStr, whereArgs := m.UniqueCond()

	if s == nil {
		s = newAutoCloseModelsSession()
	}

	_, err = s.Where(whereStr, whereArgs...).Delete(m)

	return
}
