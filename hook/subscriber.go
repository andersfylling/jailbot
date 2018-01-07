package hook

import (
	"errors"

	"github.com/andersfylling/jailbot/jailbotconst"
	"github.com/andersfylling/jailbot/notify"
	"github.com/andersfylling/unison"
	"github.com/andersfylling/unison/events"
	"github.com/sirupsen/logrus"
	"gopkg.in/bwmarrin/Discordgo.v0"
)

var SubscribeGuildHook = &unison.EventHook{
	Name:    "subscribeguild",
	Usage:   "Check if a guild has the required channel and auto-subscribe if so",
	OnEvent: unison.EventHandlerFunc(subscribeGuildHookAction),
	Events: []events.EventType{
		events.GuildCreateEvent,
		events.GuildUpdateEvent,
	},
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func subscribeGuildHookAction(ctx *unison.Context, event *events.DiscordEvent, self bool) (handled bool, err error) {
	var guild *discordgo.Guild

	// Make sure it's a new message
	if event.Type != events.GuildCreateEvent {
		return true, nil
	}

	// Convert the interface into its correct Message type
	switch ev := event.Event.(type) {
	default:
		return true, nil
	case *discordgo.GuildCreate:
		guild = ev.Guild
		break
	case *discordgo.GuildUpdate:
		guild = ev.Guild
		break
	}

	channels, err := ctx.Discord.GuildChannels(guild.ID)
	if err != nil {
		return true, err
	}

	// check that the notification channel exists
	var channel *discordgo.Channel
	for _, c := range channels {
		if stringInSlice(c.Name, jailbotconst.ChannelNamesForNotifications) {
			channel = c
			break
		}
	}

	if channel == nil {
		return true, errors.New("Could not find a notification channel for jailbot")
	}

	// subscribe to all events
	err = notify.SubscribeToAll(guild.ID, channel.ID)
	if err == nil {
		logrus.Infof("Guild %s subscribed", guild.Name)
	}

	return true, err
}
