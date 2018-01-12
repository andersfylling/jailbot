package hook

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/andersfylling/jailbot/database/dbsession"
	"github.com/andersfylling/jailbot/database/document"
	"github.com/andersfylling/jailbot/notify"
	"github.com/andersfylling/unison"
	"github.com/andersfylling/unison/events"
	"github.com/sirupsen/logrus"
	"gopkg.in/bwmarrin/Discordgo.v0"
	"gopkg.in/mgo.v2/bson"
)

var UserJoinedHook = &unison.EventHook{
	Name:    "userjoined",
	Usage:   "When a joins the guild, check his record and show the moderators",
	OnEvent: unison.EventHandlerFunc(userJoinedHookAction),
	Events: []events.EventType{
		events.GuildMemberAddEvent,
	},
}

func GetMemberRecords(ctx *unison.Context, member *discordgo.Member, guild *discordgo.Guild) (string, error) {
	// get user records
	// 1. does user exist in database ?
	userDoc := &document.UserDocument{
		DiscordID: member.User.ID,
	}
	err := userDoc.GetExisting()
	if err != nil {
		return "", errors.New("No info found of user") // not an error, user just simply didn't exist a user record
	}

	records := []*document.EventDocument{}
	ses, collection, err := dbsession.GetCollection(document.EventDocumentCollection)
	if err != nil {
		return "", err
	}
	defer ses.Close()

	// get all records
	err = collection.Find(bson.M{"userid": userDoc.DiscordID}).All(&records)
	if err != nil {
		return "", err
	} else if len(records) == 0 {
		return "", errors.New("No records found of user") // not an error, just no records exists
	}

	var nrOfBans int
	var nrOfUnbans int
	var nrOfKicks int
	var nrOfIncidents int
	for _, record := range records {
		switch record.Type {
		case notify.TypeBan:
			nrOfBans++
			break
		case notify.TypeKick:
			nrOfKicks++
			break
		case notify.TypeNotice:
			nrOfIncidents++
			break
		case notify.TypeUnban:
			nrOfUnbans++
			break
		}
	}

	// create a record string
	var recordsBuffer bytes.Buffer
	recordsLength := 0
	for ii, record := range records {
		i := ii + 1
		if recordsBuffer.Len() > 1000 {
			recordsBuffer.WriteString("...\n")
			break
		}
		guild, err := ctx.Discord.Guild(record.GuildID)
		guildName := "?"
		if err == nil {
			guildName = guild.Name // don't query the discord api due to potential rate limit issues
			// TODO: the name will be stored in the database
		}
		eventType := notify.ToStr(record.Type)
		if record.BanRemoved {
			// add number of days as a suffix
			// BAN(3days)
			start := record.ID.Time()
			end := record.BanRemovedDate
			diff := end.Sub(start)
			nrOfDays := int(diff.Hours() / 24)

			eventType += "(" + strconv.Itoa(nrOfDays) + "days)"
		}
		recordsBuffer.WriteString(fmt.Sprintf("%d. [%s] %s: --%s--\n", i, guildName, eventType, record.Reason))

		recordsLength++
	}

	// create user report
	var reportBuffer bytes.Buffer
	reportBuffer.WriteString("Report for user <@" + member.User.ID + ">\n")
	reportBuffer.WriteString("```markdown\n")
	reportBuffer.WriteString(fmt.Sprintf("# Summary about %s\n", member.User.String()))
	reportBuffer.WriteString(fmt.Sprintf("* Registered records: %d\n", len(records)))
	if nrOfBans > 0 {
		reportBuffer.WriteString(fmt.Sprintf("* Bans: %d, where %d were temporary\n", nrOfBans, nrOfUnbans))
	}
	if nrOfKicks > 0 {
		reportBuffer.WriteString(fmt.Sprintf("* Kicks: %d\n", nrOfKicks))
	}
	if nrOfIncidents > 0 {
		reportBuffer.WriteString(fmt.Sprintf("* Behavior reports: %d\n", nrOfIncidents))
	}
	reportBuffer.WriteString("\n")
	reportBuffer.WriteString(fmt.Sprintf("# Records %d/%d\n", recordsLength, len(records)))
	reportBuffer.WriteString(recordsBuffer.String())
	reportBuffer.WriteString("```")

	return reportBuffer.String(), nil
}

func userJoinedHookAction(ctx *unison.Context, event *events.DiscordEvent, self bool) (handled bool, err error) {
	logrus.Info("A user joined")
	var guildID string
	var member *discordgo.Member

	// Convert the interface into its correct type
	switch ev := event.Event.(type) {
	default:
		return true, errors.New("The event struct was not a *discordgo.GuildMemberAdd type")
	case *discordgo.GuildMemberAdd:
		member = ev.Member
		guildID = ev.GuildID
		break
	}

	guild, err := ctx.Discord.Guild(guildID)
	if err != nil {
		return true, err
	}

	// check that the guild has a notification channel before continuing
	notifyChanByteArr, err := ctx.Bot.GetGuildValue(guildID, notificationChannelKey)
	if len(notifyChanByteArr) == 0 || err != nil {
		return true, fmt.Errorf("[JailBot] Guild{name:%s, id:%s} did not have a notification channel", guild.Name, guild.ID)
	}

	report, err := GetMemberRecords(ctx, member, guild)
	if err != nil {
		return true, err
	}

	// notify server
	ctx.Discord.ChannelMessageSend(string(notifyChanByteArr), report)

	return true, nil
}
