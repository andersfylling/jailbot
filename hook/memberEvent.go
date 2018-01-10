package hook

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/andersfylling/jailbot/database/document"
	"github.com/andersfylling/jailbot/notify"
	"github.com/andersfylling/unison"
	"github.com/andersfylling/unison/events"
	"github.com/sirupsen/logrus"
	"gopkg.in/bwmarrin/Discordgo.v0"
)

const memberEventAuditLogEntriesLimit = 10

var MemberEventHook = &unison.EventHook{
	Name:    "memberevent",
	Usage:   "When a member is banned, kicked, etc. this hook handles job creation",
	OnEvent: unison.EventHandlerFunc(memberEventHookAction),
	Events: []events.EventType{
		//events.GuildBanAddEvent,       // ban
		events.GuildBanRemoveEvent,    // ban revoked
		events.GuildMemberRemoveEvent, // kick, leave, ban
	},
}

func memberEventHookAction(ctx *unison.Context, event *events.DiscordEvent, self bool) (handled bool, err error) {
	var guildID string
	var user *discordgo.User
	var removed bool

	// Convert the interface into its correct type
	switch ev := event.Event.(type) {
	default:
		return true, nil
	case *discordgo.GuildBanRemove:
		user = ev.User
		guildID = ev.GuildID
		break
	case *discordgo.GuildMemberRemove:
		user = ev.User
		guildID = ev.GuildID
		removed = true
		break
	}

	logrus.Info("[memberEvent] converted interface")

	guild, err := ctx.Discord.Guild(guildID)
	if err != nil {
		return true, err
	}
	logrus.Info("[memberEvent] got guild info")

	// make sure the bot can read the audit log
	channels, err := ctx.Discord.GuildChannels(guildID)
	if err != nil {
		return true, err
	}
	var someChannel *discordgo.Channel

	logrus.Infof("[memberEvent] got %d channels", len(channels))
	for _, channel := range channels {
		if channel.Type == 0 || channel.Type == 2 {
			someChannel = channel
			break
		}
	}
	// Don't need a specific channel as audit log reading isn't determined by channel Permissions
	// 	but "global" guild permissions.
	permissions, err := ctx.Discord.UserChannelPermissions(ctx.Bot.User.ID, someChannel.ID)
	if err != nil {
		return true, err
	} else if (permissions & discordgo.PermissionViewAuditLogs) != discordgo.PermissionViewAuditLogs {
		err = fmt.Errorf("Guild{%s, id:%s} has not granted bot permissions to read Audit logs", guild.Name, guild.ID)
		return true, err
	}
	logrus.Info("[memberEvent] got permission flags")

	// Get the last N audit log entries
	bytes, err := ctx.Discord.Request("GET", discordgo.EndpointGuilds+guildID+"/audit-logs?limit="+strconv.Itoa(memberEventAuditLogEntriesLimit), nil)

	//logrus.Info(string(bytes[:]))

	auditLog := &unison.AuditLog{}
	err = json.Unmarshal(bytes, &auditLog)
	if err != nil {
		return true, err
	}

	var eventType notify.NotificationType
	var auditEntry *unison.AuditLogEntry

	// find out what happened
	var ok bool
	var foundUser bool
	for i, entry := range auditLog.AuditLogEntries {
		if entry.TargetID == user.ID {
			foundUser = true

			if removed {
				if entry.ActionType == 20 {
					// kick
					eventType = notify.TypeKick
					ok = true
				} else if entry.ActionType == 22 {
					// ban
					eventType = notify.TypeBan
					ok = true
				}
			} else if entry.ActionType == 23 {
				// removed ban
				eventType = notify.TypeUnban
				ok = true
			}

			if ok {
				auditEntry = entry
				break
			} else if i > memberEventAuditLogEntriesLimit {
				break
			}
		}
	}
	if !ok {
		if !foundUser {
			return true, errors.New("unable to find user{id:" + user.ID + "} in audit log")
		} else {
			return true, errors.New("found user in audit log, but weren't able to get correct log entry")
		}
	}

	// save ban to database
	// 	1. check if user exist in db
	userDoc := &document.UserDocument{
		DiscordID: user.ID,
	}
	err = userDoc.GetExisting()
	if err != nil { // user not found
		userDoc.Avatar = user.Avatar
		userDoc.Bot = user.Bot
		userDoc.Discriminator = user.Discriminator
		userDoc.Email = user.Email
		userDoc.MFAEnabled = user.MFAEnabled
		userDoc.Token = user.Token
		userDoc.Username = user.Username
		userDoc.Verified = user.Verified

		// 1.1 if it doesnt exist, save it to db
		id, err := userDoc.Insert()
		if err != nil {
			// unable to save...
			return true, err
		}

		// saved
		userDoc.ID = id
	}
	//  2. check if guild exist in db
	guildDoc := &document.GuildDocument{
		DiscordID: guild.ID,
	}
	err = guildDoc.GetExisting()
	if err != nil { // guild not found
		guildDoc.Icon = guild.Icon
		guildDoc.Name = guild.Name
		guildDoc.OwnerID = guild.OwnerID
		guildDoc.Region = guild.Region

		// 2.1 if it doesnt exist, save it to db
		id, err := guildDoc.Insert()
		if err != nil {
			// unable to save...
			return true, err
		}

		// saved
		guildDoc.ID = id
	}
	//  3. create new EventDocument with Type == ban
	eventDoc := &document.EventDocument{
		GuildID: guildDoc.DiscordID,
		UserID:  userDoc.DiscordID,
		Type:    eventType,
		Reason:  auditEntry.Reason,
	}
	//  4. save
	_, err = eventDoc.Insert()
	if err != nil {
		return true, err
	}

	// publish
	notification := notify.NewNotification2(eventType, user, guild, auditEntry)
	notify.Publish(ctx, notification)

	return true, nil
}
