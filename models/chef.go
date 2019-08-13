package models

type Chef struct {
	Idx      int    `bson:"idx" json:"idx" redis:"idx"`
	No       string `bson:"no" json:"no" redis:"no"`
	Level    int    `bson:"level" json:"level" redis:"level"`
	StoveIdx int    `bson:"stove_idx" json:"stove_idx" redis:"stove_idx"`
}
