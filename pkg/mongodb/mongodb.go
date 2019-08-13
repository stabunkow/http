package mongodb

import (
	"stabunkow/http/pkg/setting"

	"gopkg.in/mgo.v2"
)

type Db struct {
	*mgo.Session
	Database string
	Addr     string
}

var DefaultDb *Db

func GetDefaultDb() *Db {
	return DefaultDb
}

func Setup() error {
	DefaultDb = NewMongodb(
		setting.MongodbSetting.Addr,
		setting.MongodbSetting.Database,
	)

	return nil
}

func NewMongodb(addr, database string) *Db {
	session, err := mgo.Dial(addr)

	if err != nil {
		panic("mgo Dial failed: " + err.Error())
	}

	session.SetMode(mgo.Monotonic, true)
	return &Db{
		Session:  session,
		Addr:     addr,
		Database: database,
	}
}

type MgoSession struct {
	*mgo.Session
	Database string
}

func (db *Db) Conn() *MgoSession {
	return &MgoSession{
		Session:  db.Session.Copy(),
		Database: db.Database,
	}
}

func (session *MgoSession) Collection(table string) *mgo.Collection {
	return session.DB(session.Database).C(table)
}
