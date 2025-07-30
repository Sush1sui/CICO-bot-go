package repository

import (
	"github.com/Sush1sui/cico-bot-go/internal/models"
	"github.com/bwmarrin/discordgo"
)

type ClockRecordInterface interface {
	ClockIn(userId string) (*models.ClockRecordModel, error)
	ClockOut(userId string) (*models.ClockRecordModel, error)
	CheckForExpiredClock(s *discordgo.Session) error
	HandleIfExpiredClock(s *discordgo.Session, userId, roleId string) bool
	GetUserClockRecord(userId string) (*models.ClockRecordModel, error)
	GetAllClockRecords() ([]*models.ClockRecordModel, error)
	ReClockUser(userId string) (*models.ClockRecordModel, error)
	RemoveClockRecordOfThoseNotClockedIn() error
}

var ClockRecordService ClockRecordInterface