package deploy

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var eventHandlers = []any{}

func DeployEvents(s *discordgo.Session) {
	if len(eventHandlers) == 0 { return }

	for _, handler := range eventHandlers {
		s.AddHandler(handler)
	}

	fmt.Println("Event handlers deployed successfully")
}