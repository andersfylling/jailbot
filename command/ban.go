package command

import (
	"errors"
	"fmt"
	"regexp"

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

func retrieveTargetUserID(request string /*, args BanCommandArgs*/) string {
	args := BanCommandArgs // create a copy
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
	userID := retrieveTargetUserID(request) //, BanCommandArgs will be used here, copy on run
	reason := BanCommandArgs.Reason
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

	return err
}
