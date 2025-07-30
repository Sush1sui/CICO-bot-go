package repository

import "github.com/Sush1sui/cico-bot-go/internal/models"

type ClockChannelInterface interface {
	GetAllClockChannelInterface() ([]*models.ClockChannelModel, error)
	CreateClockChannelInterface(categoryId, clockInChannelId, clockInInterfaceId, clockOutChannelId, clockOutInterfaceId, adminChannelId, clockInRoleId string) (*models.ClockChannelModel, error)
	DeleteAllClockChannelInterface() error
}

var ClockChannelService ClockChannelInterface