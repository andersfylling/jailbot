package commonutil

import (
	"fmt"
	"strconv"
	"time"

	"github.com/andersfylling/jailbot/database/document"
	"gopkg.in/bwmarrin/Discordgo.v0"
)

// FmtEventRecordToStr formats the ban/kick/unban/report record into human readable string
func FmtEventRecordToStr(record *document.EventDocument, guild *discordgo.Guild, user *discordgo.User) string {
	if record == nil {
		return ""
	}

	eventType := NotificationTypeToStr(record.Type)
	if record.BanRemoved {
		// add number of days as a suffix
		// BAN(3days)
		start := record.ID.Time()
		end := record.BanRemovedDate
		diff := end.Sub(start)
		nrOfDays := int(diff.Hours() / 24)

		eventType += "(" + strconv.Itoa(nrOfDays) + "days)"
	}

	guildName := "?"
	if guild != nil {
		guildName = guild.Name
	}

	var userInfo string
	if user != nil {
		userInfo = fmt.Sprintf("[<@%s>]", user.ID)
	}

	var reason string
	if record.Reason != "" {
		reason = ", reason: " + record.Reason
	}

	// [GUILD_NAME] BAN(3days): --some reason for why the user was banned--
	// [GUILD_NAME] BAN: --some reason for why the user was banned--
	date := record.ID.Time()
	return fmt.Sprintf("[%s, %s] %s%s%s\n", guildName, date.Format(time.RFC822), eventType, userInfo, reason)
}
