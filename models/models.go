package models

import (
	"gopkg.in/mgo.v2/bson"
)

type ObjectId bson.ObjectId

// func (id *ObjectId) RedisScan(src interface{}) error {
// 	b, ok := src.([]byte)

// 	if !ok {
// 		return fmt.Errorf("cannot convert objectId from %T to %T", src, b)
// 	}

// 	*id = bson.ObjectIdHex(string(b))

// 	return nil
// }

// func (id ObjectId) RedisArg() interface{} {
// 	return id.Hex()
// }
