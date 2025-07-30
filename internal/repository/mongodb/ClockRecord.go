package mongodb

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/Sush1sui/cico-bot-go/internal/config"
	"github.com/Sush1sui/cico-bot-go/internal/models"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (c *MongoClient) ClockIn(userId string) (*models.ClockRecordModel, error) {
	clockInTime := time.Now()
	
	var existingRecord models.ClockRecordModel
	err := c.Client.FindOne(context.Background(), bson.M{"userId": userId}).Decode(&existingRecord)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("error checking existing clock record: %w", err)
	}

	// if record exists
	if err == nil {
		// if already clocked in
		if !existingRecord.ClockOutTime.IsZero() && (existingRecord.ClockOutTime == nil || existingRecord.ClockOutTime.IsZero()) {
			return nil, fmt.Errorf("user %s is already clocked in", userId)
		}

		// record exists but clocked out
		update := bson.M{
			"$set": bson.M{
				"clockInTime": clockInTime,
			},
			"$unset": bson.M{"clockOutTime": ""},
		}

		_, err := c.Client.UpdateByID(context.Background(), existingRecord.ID, update)
		if err != nil {
			return nil, fmt.Errorf("error updating existing clock record: %w", err)
		}

		// fetch updated record
		err = c.Client.FindOne(context.Background(), bson.M{"_id": existingRecord.ID}).Decode(&existingRecord)
		if err != nil {
			return nil, fmt.Errorf("error fetching updated clock record: %w", err)
		}
		return &existingRecord, nil
	}

	// if no existing record, create a new one
	newRecord := models.ClockRecordModel{
		UserID:      userId,
		ClockInTime: &clockInTime,
	}
	res, err := c.Client.InsertOne(context.Background(), newRecord)
	if err != nil {
		return nil, fmt.Errorf("error inserting new clock record: %w", err)
	}

	newRecord.ID = res.InsertedID.(bson.ObjectID)
	return &newRecord, nil
}

func (c *MongoClient) ClockOut(userId string) (*models.ClockRecordModel, error) {
	clockOutTime := time.Now()

	var existingRecord models.ClockRecordModel
	err := c.Client.FindOne(context.Background(), bson.M{
		"userId": userId,
		"clockOutTime": nil,
	}).Decode(&existingRecord)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no active clock record found for user %s", userId)
		}
		return nil, fmt.Errorf("error fetching clock record: %w", err)
	}

	if existingRecord.ClockInTime == nil || existingRecord.ClockInTime.IsZero() {
		return nil, fmt.Errorf("user %s has not clocked in", userId)
	}

	totalHours := clockOutTime.Sub(*existingRecord.ClockInTime).Hours()
	if existingRecord.TotalHours != nil {
		totalHours += *existingRecord.TotalHours
	}

	update := bson.M{
		"$set": bson.M{
			"clockOutTime": &clockOutTime,
			"totalHours":   totalHours,
		},
		"$unset": bson.M{"clockInTime": ""},
	}

	_, err = c.Client.UpdateByID(context.Background(), existingRecord.ID, update)
	if err != nil {
		return nil, fmt.Errorf("error updating clock record: %w", err)
	}

	// fetch updated record
	err = c.Client.FindOne(context.Background(), bson.M{"_id": existingRecord.ID}).Decode(&existingRecord)
	if err != nil {
		return nil, fmt.Errorf("error fetching updated clock record: %w", err)
	}

	return &existingRecord, nil
}

func (c *MongoClient) CheckForExpiredClock(s *discordgo.Session) error {
	fmt.Println("Checking for expired clock records...")

	cursor, err := c.Client.Find(context.Background(), bson.M{})
	if err != nil {
		return fmt.Errorf("error fetching clock records: %w", err)
	}
	defer cursor.Close(context.Background())

	guild, err := s.State.Guild(config.GlobalConfig.GuildID)
	if err != nil {
		return fmt.Errorf("error fetching guild: %w", err)
	}

	// build a map of user IDs to member for quick lookup
	memberMap := make(map[string]*discordgo.Member)
	for _, member := range guild.Members {
		memberMap[member.User.ID] = member
	}

	for cursor.Next(context.Background()) {
		var record models.ClockRecordModel
		if err := cursor.Decode(&record); err != nil {
			return fmt.Errorf("error decoding clock record: %w", err)
		}

		member := memberMap[record.UserID]
		if member == nil {
			fmt.Printf("User %s not found in guild %s\n", record.UserID, config.GlobalConfig.GuildID)
			continue
		}

		if slices.Contains(member.Roles, config.GlobalConfig.TL_ROLE_ID) {
			c.HandleIfExpiredClock(s, record.UserID, config.GlobalConfig.TL_ROLE_ID)
		} else if slices.Contains(member.Roles, config.GlobalConfig.CHATTER_ROLE_ID) {
			c.HandleIfExpiredClock(s, record.UserID, config.GlobalConfig.CHATTER_ROLE_ID)
		}
	}
 
	fmt.Println("Finished checking for expired clock records.")
	return nil
}

func (c *MongoClient) HandleIfExpiredClock(s *discordgo.Session, userId, roleId string) bool {
	var userRecord models.ClockRecordModel
	err := c.Client.FindOne(context.Background(), bson.M{"userId": userId}).Decode(&userRecord)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("no clock record found for user:", userId)
		}
		fmt.Printf("Error fetching user clock record: %v\n", err)
		return  false
	}

	totalHours := time.Since(*userRecord.ClockInTime).Hours()
	guild, err := s.State.Guild(config.GlobalConfig.GuildID)
	if err != nil {
		fmt.Printf("Error fetching guild: %v\n", err)
		return false
	}
	var discordMember *discordgo.Member
	hasRole := false
	for _, member := range guild.Members {
		if member.User.ID == userId {
			discordMember = member
			break
		}
	}
	if discordMember == nil {
		fmt.Printf("User %s not found in guild %s\n", userId, config.GlobalConfig.GuildID)
		return false
	}
	if slices.Contains(discordMember.Roles, roleId) {
			hasRole = true
	}
	if !hasRole {
		fmt.Printf("User %s does not have the required role %s\n", userId, roleId)
		return false
	}

	timeLimit, exists := config.GlobalConfig.TimeLimit[roleId]
	if !exists {
		fmt.Printf("No time limit found for role %s\n", roleId)
		return false
	}

	if totalHours > timeLimit {
		clockOutTime := time.Now()
		fmt.Println("User has exceeded the time limit:", discordMember.User.Username)

		_, err = c.Client.UpdateOne(context.Background(), bson.M{"userId": userId}, bson.M{
			"$set": bson.M{
				"clockOutTime": &clockOutTime,
			},
			"$unset": bson.M{"clockInTime": ""},
		})
		if err != nil {
			fmt.Printf("Error updating clock record for user %s: %v\n", userId, err)
			return false
		}

		adminChannel, err := s.State.Channel(config.GlobalConfig.AdminChannelID)
		if err != nil {
			fmt.Printf("Error fetching admin channel: %v\n", err)
			return false
		}

		_, err = s.ChannelMessageSend(adminChannel.ID, fmt.Sprintf("⚠️ <@%s> has exceeded the time limit of %.2f hours for role %s. Clocking out now.", userId, timeLimit, roleId))
		if err != nil {
			fmt.Printf("Error sending message to admin channel: %v\n", err)
			return false
		}
		return false
	}
	return true
}

func (c *MongoClient) GetUserClockRecord(userId string) (*models.ClockRecordModel, error) {
	var record models.ClockRecordModel
	err := c.Client.FindOne(context.Background(), bson.M{"userId": userId}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no clock record found for user %s", userId)
		}
		return nil, fmt.Errorf("error fetching clock record: %w", err)
	}
	return &record, nil
}

func (c *MongoClient) GetAllClockRecords() ([]*models.ClockRecordModel, error) {
	var records []*models.ClockRecordModel
	cursor, err := c.Client.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error fetching clock records: %w", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var record models.ClockRecordModel
		if err := cursor.Decode(&record); err != nil {
			return nil, fmt.Errorf("error decoding clock record: %w", err)
		}
		records = append(records, &record)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cursor: %w", err)
	}
	return records, nil
}

func (c *MongoClient) RemoveClockRecord(userId string) (int, error) {
	result, err := c.Client.DeleteOne(context.Background(), bson.M{"userId": userId})
	if err != nil {
		return 0, fmt.Errorf("error removing clock record: %w", err)
	}
	return int(result.DeletedCount), nil
}