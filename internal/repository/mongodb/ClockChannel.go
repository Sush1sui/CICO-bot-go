package mongodb

import (
	"context"
	"fmt"

	"github.com/Sush1sui/cico-bot-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (c *MongoClient) GetAllClockChannelInterface() ([]*models.ClockChannelModel, error) {
	var clockChannels []*models.ClockChannelModel

	cursor, err := c.Client.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve clock channels: %w", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var clockChannelSetup models.ClockChannelModel
		if err := cursor.Decode(&clockChannelSetup); err != nil {
			return nil, fmt.Errorf("failed to decode clock channel: %w", err)
		}
		clockChannels = append(clockChannels, &clockChannelSetup)
	}
	return clockChannels, nil
}

func (c *MongoClient) CreateClockChannelInterface(categoryId, clockInChannelId, clockInInterfaceId, clockOutChannelId, clockOutInterfaceId, adminChannelId, clockInRoleId string) (*models.ClockChannelModel, error) {
	clockChannels := &models.ClockChannelModel{
		CategoryID: categoryId,
		ClockInChannelID: clockInChannelId,
		ClockInInterfaceID: clockInInterfaceId,
		ClockOutChannelID: clockOutChannelId,
		ClockOutInterfaceID: clockOutInterfaceId,
		AdminChannelID: adminChannelId,
		ClockInRoleID: clockInRoleId,
	}

	if _, err := c.Client.InsertOne(context.Background(), clockChannels); err != nil {
		return nil, fmt.Errorf("failed to create clock channel: %w", err)
	}
	return clockChannels, nil
}