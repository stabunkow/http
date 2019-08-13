package models

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id  bson.ObjectId `bson:"_id,omitempty" json:"id" redis:"id"`
	Sid string        `bson:"sid,omitempty" json:"-" redis:"sid"` // like token, as user certificate, default ''

	WechatOpenId     string `bson:"wechat_open_id,omitempty" json:"-" redis:"wechat_open_id"`
	WechatSessionKey string `bson:"wechat_session_key,omitempty" json:"-" redis:"wechat_session_key"`

	Coins    string `bson:"coins" json:"coins" redis:"coins"`
	Diamonds string `bson:"diamonds" json:"diamonds" redis:"diamonds"`

	Restaurant *Restaurant `bson:"restaurant" json:"restaurant" redis:"-"`
	Chefs      []*Chef     `bson:"chefs" json:"chefs" redis:"-"`

	CreatedAt int64 `bson:"created_at" json:"created_at" redis:"created_at"`
	UpdatedAt int64 `bson:"updated_at" json:"updated_at" redis:"updated_at"`
}

func (m *User) GetId() string {
	return m.Id.Hex()
}

func (m *User) SetId(id string) {
	m.Id = bson.ObjectIdHex(id)
}
