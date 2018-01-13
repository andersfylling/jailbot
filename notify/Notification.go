package notify

import (
	"github.com/andersfylling/jailbot/common"
	"github.com/andersfylling/jailbot/database/document"
	"gopkg.in/bwmarrin/Discordgo.v0"
)

const (
	TypeBan    common.NotificationType = "BAN"
	TypeUnban  common.NotificationType = "UNBAN"
	TypeKick   common.NotificationType = "KICK"
	TypeNotice common.NotificationType = "NOTICE" // general way of reporting behavior..
)

type Notification struct {
	Type common.NotificationType

	Record *document.EventDocument
	User   *discordgo.User
	Guild  *discordgo.Guild
}

func NewNotification(t common.NotificationType, record *document.EventDocument, user *discordgo.User, guild *discordgo.Guild) *Notification {
	return &Notification{
		Type: t,

		Record: record,
		User:   user,
		Guild:  guild,
	}
}
