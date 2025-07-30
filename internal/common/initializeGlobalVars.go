package common

import (
	"fmt"

	"github.com/Sush1sui/cico-bot-go/internal/config"
	"github.com/Sush1sui/cico-bot-go/internal/repository"
)

func InitializeGlobalVars() error {
	clockChannelSetups, err := repository.ClockChannelService.GetAllClockChannelInterface()
	if err != nil {
		return fmt.Errorf("error retrieving clock channel setups: %v", err)
	}
	if len(clockChannelSetups) == 0 {
		return fmt.Errorf("no clock channel setups found in the database")
	}

	clockChannel := clockChannelSetups[0]
	if clockChannel == nil {
		return fmt.Errorf("no valid clock channel found in the database")
	}
	if clockChannel.AdminChannelID == "" {
		return fmt.Errorf("admin channel ID is missing in the clock channel")
	}
	if clockChannel.ClockInRoleID == "" {
		return fmt.Errorf("clock in role ID is missing in the clock channel")
	}
	config.GlobalConfig.AdminChannelID = clockChannel.AdminChannelID
	config.GlobalConfig.ClockInRoleID = clockChannel.ClockInRoleID
	return nil
}