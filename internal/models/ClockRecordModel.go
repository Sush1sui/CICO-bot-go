package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ClockRecordModel struct {
	ID bson.ObjectID `bson:"_id,omitempty"`
	UserID string 	  `bson:"userId"`
	ClockInTime *time.Time `bson:"clockInTime, omitempty"`
	ClockOutTime *time.Time `bson:"clockOutTime, omitempty"`
	TotalHours *float64 `bson:"totalHours,omitempty"`
}