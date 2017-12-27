package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/andersfylling/jailbot/command"
	"github.com/andersfylling/unison"
)

func main() {
	// Create bot structure
	settings := &unison.Config{
		Commands: []*unison.Command{
			command.BanCommand,
		},
		EventHooks: []*unison.EventHook{},
		Services:   []*unison.Service{},
	}

	// Start the bot
	err := unison.Run(settings)
	if err != nil {
		logrus.Error(err)
	}
}
