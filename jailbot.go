package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/s1kx/unison"
)

func main() {
	// Create bot structure
	settings := &unison.BotSettings{
		Commands:   []*unison.Command{},
		EventHooks: []*unison.EventHook{},
		Services:   []*unison.Service{},
	}

	// Start the bot
	err := unison.Run(settings)
	if err != nil {
		logrus.Error(err)
	}
}
