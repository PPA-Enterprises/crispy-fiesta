package db

import (
	"os"
	"time"

	"gopkg.in/mgo.v2"
)

type DBConnection struct {
	session *mgo.Session
}

func NewConnection(host string, dbName string) (conn *DBConnection) {
	info := &mgo.DialInfo{
		Addrs:    []string{host},
		Timeout:  60 * time.Second,
		Database: dbName,
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PWD"),
	}

	session, err := mgo.DialWithInfo(info)

	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	conn = &DBConnection{session}
	return conn
}

func (conn *DBConnection) Use(dbName, tableName string) (collection *mgo.Collection) {
	return conn.session.DB(dbName).C(tableName)
}

func (conn *DBConnection) Close() {
	conn.session.Close()
	return
}
