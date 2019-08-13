package models

type Restaurant struct {
	Level           int      `bson:"level" json:"level" redis:"level"`
	Scene           int      `bson:"scene" json:"scene" redis:"scene"`
	ConveyBeltLevel int      `bson:"convey_belt_level" json:"convey_belt_level" redis:"convey_belt_level"`
	Stoves          []*Stove `bson:"stoves" json:"stoves" redis:"-"`
}
