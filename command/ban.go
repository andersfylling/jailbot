package jailbotcmd


import (
	"github.com/bwmarrin/discordgo"
	"github.com/s1kx/unison"
	"strconv"
  "github.com/andersfylling/jailbot/jailbotque"
)

var BanCommand = &unison.Command{
	Name:        "ban",
	Description: "Ban a member",
	Action:      banCommandAction,
	Deactivated: false,
	Permission:  unison.NewCommandPermission(),
}
// PermissionRole Check if user has permission to deal with roles
func permissionRole(ctx *unison.Context, m *discordgo.Message) bool {
	authorPermissions, _ := ctx.Bot.Discord.UserChannelPermissions(m.Author.ID, m.ChannelID)
  required := discordgo.PermissionBanMembers | discordgo.PermissionKickMembers

	return (authorPermissions & required) == required
}

func retrieveTargetedMemberFromStr(input string) string {
  // --user="username#1234"
  // --userid="57435436543523345657"
  // --uid="57438745734657"
  // @mention, another user might be mentioned in a option or later.. first mention(?)

  
}

func banCommandAction(ctx *unison.Context, m *discordgo.Message, content string) error {
	if !permissionRole(ctx, m) {
		return errors.New("Member " + m.Author.Username + " do not have kick+ban permissions")
  }

  guildID := ctx.Discord.Channel(m.ChannelID).GuildID
  userID := ""
  reason := ""
  days := int(0) // remove all messages from the last X days

  // extract the member data before he gets banned
  member, err := ctx.Discord.GuildMember(guildID, userID)
  if err != nil {
    return errors.New("Unable to find a member in the guild with id " + userID)
  }

  // get guild details
  guild, err := ctx.Discord.Guild(guildID)
  if err != nil {
    return errors.New("Unable to get guild information using guildID: " + guildID)
  }

  // ban member from guild
  ctx.Discord.GuildBanCreate(guildID, userID, days)

  // create the ban alert
  banAlert := jailbotque.NewAlert(guild, member, reason)

  // add the ban alert to stack
  alertQue := jailbotque.GetInstance()
  alertQue.Alerts.Add(alert)

  // let the moderator know what has been done
  msg := fmt.Sprintf("Banned user %s and removed messages within the last %d days", userID, days)
	_, err := ctx.Bot.Discord.ChannelMessageSend(m.ChannelID, msg)

	return err
}
