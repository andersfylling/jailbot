package main

import (
	"github.com/andersfylling/jailbot/command"
	"github.com/andersfylling/jailbot/hook"
	"github.com/andersfylling/unison"
	"github.com/sirupsen/logrus"
)

func main() {

	// Create bot structure
	settings := &unison.Config{
		Commands: []*unison.Command{
			command.BanCommand,
		},
		EventHooks: []*unison.EventHook{
			hook.SubscribeGuildHook,
			hook.MemberEventHook,
		},
		Services: []*unison.Service{},
	}

	// Start the bot
	err := unison.Run(settings)
	if err != nil {
		logrus.Error(err)
	}
}
