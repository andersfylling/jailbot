package command

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Sirupsen/logrus"
	"github.com/andersfylling/jailbot/hook"
	"github.com/andersfylling/unison"
	"gopkg.in/bwmarrin/Discordgo.v0"
)

var RecordsCommandArgs struct {
	User   string `arg:"-u" help:"mention targeted member"`
	UserID string `arg:"-i" help:"specify a member id"`
}

var RecordsCommand = &unison.Command{
	Name:        "records",
	Usage:       "Retrieve records and a summary of a user",
	Action:      recordsCommandAction,
	Deactivated: false,
	Flags:       &RecordsCommandArgs,
}

func recordsRetrieveTargetUserID(request string /*, args BanCommandArgs*/) string {
	args := RecordsCommandArgs // create a copy
	// --user="<@57435436543523345657>", --user="<@!57435436543523345657>", -u="<@57435436543523345657>"
	// --userid="57435436543523345657", --i="57438745734657"
	//
	// just mention: ban @user
	// 	Assumption: the first mention if --user nor --userid is set, is the target user

	if args.User != "" {
		var mention string
		if args.User[:3] == "<@!" { // nickname
			mention = "<@!"
		} else if args.User[:2] == "<@" {
			mention = "<@"
		} else {
			return ""
		}
		return args.User[len(mention) : len(args.User)-len(">")] // Extracts uid from eg. <@237293823443199860>
	} else if args.UserID != "" {
		return args.UserID
	} else {
		// mention
		if request == "" {
			return ""
		}
		userIDPattern := regexp.MustCompile("<@[^>]*>")
		mention := userIDPattern.FindString(request)

		userStr := mention[2 : len(mention)-1]
		// if a nickname is used, this will contain a prefix of `!`
		var userID string
		if userStr[:1] == "!" {
			userID = userStr[1:len(userStr)]
		} else {
			userID = userStr
		}

		return userID
	} //else if BanCommandArgs.UserName != "" {
	//	return BanCommandArgs.UserName // TODO extract username+discriminator and find a user id
	//}
}

func readAuditLogPermission(ctx *unison.Context, m *discordgo.Message) bool {
	authorPermissions, _ := ctx.Bot.Discord.UserChannelPermissions(m.Author.ID, m.ChannelID)
	required := discordgo.PermissionViewAuditLogs

	return (authorPermissions & required) == required
}

func recordsCommandAction(ctx *unison.Context, m *discordgo.Message, request string) error {
	if !readAuditLogPermission(ctx, m) {
		return errors.New("Member " + m.Author.Username + " do not have permissions to read audit logs")
	}

	channel, err := ctx.Discord.Channel(m.ChannelID)
	if err != nil {
		return fmt.Errorf("Unable to find channel %s", m.ChannelID)
	}

	guildID := channel.GuildID
	userID := recordsRetrieveTargetUserID(request)

	if userID == "" {
		return errors.New("Unable to find a user id")
	}

	// extract the targeted member data
	member, err := ctx.Discord.GuildMember(guildID, userID)
	if err != nil {
		return errors.New("Unable to find a member in the guild with id " + userID)
	}

	moderator, err := ctx.Discord.GuildMember(guildID, m.Author.ID)
	if err != nil {
		return err
	}

	guild, err := ctx.Discord.Guild(guildID)
	if err != nil {
		return err
	}

	report, err := hook.GetMemberRecords(ctx, member, guild)
	if err != nil {
		return err
	}

	_, err = ctx.Bot.Discord.ChannelMessageSend(m.ChannelID, report)

	// let the moderator know what has been done
	logEntry := fmt.Sprintf("Guild{%s, id:%s} member{%s, id:%s} requested a report about user{%s, id:%s}", guild.Name, guild.ID, moderator.User.String(), moderator.User.ID, member.User.String(), member.User.ID)
	logrus.Info(logEntry)

	return err
}
