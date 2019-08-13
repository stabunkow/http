package models

type Stove struct {
	Idx     int `bson:"idx" json:"idx" redis:"idx"`
	Level   int `bson:"level" json:"level" redis:"level"`
	ChefIdx int `bson:"chef_idx" json:"chef_idx" redis:"chef_idx"`
}
