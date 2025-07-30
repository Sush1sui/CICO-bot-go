package models

import "go.mongodb.org/mongo-driver/v2/bson"

type ClockChannelModel struct {
	ID bson.ObjectID `bson:"_id,omitempty"`
	CategoryID   string        `bson:"categoryId"`
	ClockInChannelID  string     `bson:"clockInChannelId"`
	ClockInInterfaceID string 	  `bson:"clockInInterfaceId"`
	ClockOutChannelID string     `bson:"clockOutChannelId"`
	ClockOutInterfaceID string `bson:"clockOutInterfaceId"`
	AdminChannelID string     `bson:"adminChannelId"`
	ClockInRoleID string        `bson:"clockInRoleId"`
}