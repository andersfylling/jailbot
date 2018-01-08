package command

import (
	"errors"
	"fmt"

	"github.com/andersfylling/jailbot/database/document"
	"github.com/andersfylling/jailbot/notify"
	"github.com/andersfylling/unison"
	"github.com/sirupsen/logrus"
	"gopkg.in/bwmarrin/Discordgo.v0"
)

var BanCommandArgs struct {
	User string `arg:"-u" help:"mention targeted member"`
	//	UserName string `arg:"-un" help:"use username#discriminator as target value"`
	UserID string `arg:"-i" help:"specify a member id"`
	Reason string `arg:"-r" help:"reason for banning member"`
	Days   int    `arg:"-d" help:"Delete target members message N days back in time"`
}

var BanCommand = &unison.Command{
	Name:        "ban",
	Usage:       "Ban a member",
	Action:      banCommandAction,
	Deactivated: false,
	Flags:       &BanCommandArgs,
}

// PermissionRole Check if user has permission to deal with roles
func permissionRole(ctx *unison.Context, m *discordgo.Message) bool {
	authorPermissions, _ := ctx.Bot.Discord.UserChannelPermissions(m.Author.ID, m.ChannelID)
	required := discordgo.PermissionBanMembers | discordgo.PermissionKickMembers

	return (authorPermissions & required) == required
}

func retrieveTargetUserID() string {
	// --user="username#1234" .. neh
	// --userid="57435436543523345657"
	// --uid="57438745734657"
	// @mention, another user might be mentioned in a option or later.. first mention(?)

	if BanCommandArgs.User != "" {
		return BanCommandArgs.User[len("<@") : len(BanCommandArgs.User)-len(">")] // Extracts uid from eg. <@237293823443199860>
	} else if BanCommandArgs.UserID != "" {
		return BanCommandArgs.UserID
	} //else if BanCommandArgs.UserName != "" {
	//	return BanCommandArgs.UserName // TODO extract username+discriminator and find a user id
	//}

	// TODO: should users be able to mention someone without using a --user optional?

	return ""
}

func banCommandAction(ctx *unison.Context, m *discordgo.Message, request string) error {
	logrus.Info("[BanCommand] Executing")
	if !permissionRole(ctx, m) {
		return errors.New("Member " + m.Author.Username + " do not have kick+ban permissions")
	}
	logrus.Info("[BanCommand] Permissions OK")

	channel, err := ctx.Discord.Channel(m.ChannelID)
	if err != nil {
		return fmt.Errorf("Unable to find channel %s", m.ChannelID)
	}
	logrus.Info("[BanCommand] Detected channel id")
	guildID := channel.GuildID
	userID := retrieveTargetUserID()
	reason := BanCommandArgs.Reason
	labels := ""
	days := BanCommandArgs.Days // remove all messages from the last X days

	if userID == "" {
		return errors.New("Unable to find a user id")
	}
	logrus.Info("[BanCommand] Found user id")

	// extract the member data before he gets banned
	member, err := ctx.Discord.GuildMember(guildID, userID)
	if err != nil {
		return errors.New("Unable to find a member in the guild with id " + userID)
	}

	// get guild details
	_, err = ctx.Discord.Guild(guildID)
	if err != nil {
		return errors.New("Unable to get guild information using guildID: " + guildID)
	}

	// ban member from guild
	err = ctx.Discord.GuildBanCreateWithReason(guildID, userID, reason, days)
	if err != nil {
		return err
	}
	logrus.Info(fmt.Sprintf("Banned user %s{id:%s} from Guild %s with reason %s", member.User.String(), userID, guildID, reason))

	// let the moderator know what has been done
	msg := fmt.Sprintf("Banned user %s{id:%s} and removed messages within the last %d days", member.User.String(), userID, days)
	_, err = ctx.Bot.Discord.ChannelMessageSend(m.ChannelID, msg)
	guild, _ := ctx.Discord.Guild(guildID)

	// save ban to database
	// 	1. check if user exist in db
	userDoc := &document.UserDocument{
		DiscordID: member.User.ID,
	}
	err = userDoc.GetExisting()
	if err != nil { // user not found
		user := member.User
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
			return err
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
			return err
		}

		// saved
		guildDoc.ID = id
	}
	//  3. create new EventDocument with Type == ban
	eventDoc := &document.EventDocument{
		GuildID: guildDoc.ID,
		UserID:  userDoc.ID,
		Type:    document.EventTypeBan,
		Reason:  reason,
		Labels:  labels,
	}
	//  4. save
	_, err = eventDoc.Insert()
	if err != nil {
		return err
	}

	// publish news
	notification := notify.NewBanNotification(userID, member.User.Username, member.User.Discriminator, guild.Name, guildID, guild.MemberCount, reason)

	notify.Publish(ctx, notification)

	return err
}
