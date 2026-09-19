package models

import "time"

type Vote struct {
	ID        interface{} `bson:"_id,omitempty" json:"id"`
	PollID    interface{} `bson:"pollId" json:"pollId"`
	OptionID  string      `bson:"optionId" json:"optionId"`
	VoterID   interface{} `bson:"voterId" json:"voterId"`
	CreatedAt time.Time   `bson:"createdAt" json:"createdAt"`
}
