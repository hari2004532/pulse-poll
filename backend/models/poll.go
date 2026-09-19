package models

import "time"

type PollOption struct {
	ID   string `bson:"id" json:"id"`
	Text string `bson:"text" json:"text"`
}

type Poll struct {
	ID        interface{}  `bson:"_id,omitempty" json:"id"`
	Question  string       `bson:"question" json:"question"`
	Options   []PollOption `bson:"options" json:"options"`
	CreatedBy interface{}  `bson:"createdBy" json:"createdBy"`
	CreatedAt time.Time    `bson:"createdAt" json:"createdAt"`
	IsActive  bool         `bson:"isActive" json:"isActive"`
}
