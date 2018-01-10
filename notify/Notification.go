package notify

import (
	"time"

	"github.com/andersfylling/unison"
	"gopkg.in/bwmarrin/Discordgo.v0"
)

type NotificationType string

const (
	TypeBan    NotificationType = "BAN"
	TypeUnban  NotificationType = "UNBAN"
	TypeKick   NotificationType = "KICK"
	TypeNotice NotificationType = "NOTICE" // general way of reporting behavior..
)

func ToStr(nt NotificationType) string {
	return string(nt)
}

type Notification struct {
	Type NotificationType

	UserID            string
	UserName          string
	UserDiscriminator string

	GuildName string
	GuildID   string
	GuildSize int // number of members

	Reason string

	Date time.Time
}

func NewNotification(t NotificationType, userid, username, userdiscriminator, guildname, guildid string, guildsize int, reason string) *Notification {
	return &Notification{
		Type: t,

		UserID:            userid,
		UserName:          username,
		UserDiscriminator: userdiscriminator,

		GuildID:   guildid,
		GuildName: guildname,
		GuildSize: guildsize,

		Reason: reason,

		Date: time.Now(),
	}
}

func NewNotification2(t NotificationType, user *discordgo.User, guild *discordgo.Guild, auditEntry *unison.AuditLogEntry) *Notification {
	return &Notification{
		Type: t,

		UserID:            user.ID,
		UserName:          user.Username,
		UserDiscriminator: user.Discriminator,

		GuildID:   guild.ID,
		GuildName: guild.Name,
		GuildSize: guild.MemberCount,

		Reason: auditEntry.Reason,

		Date: time.Now(),
	}
}

func NewBanNotification(userid, username, userdiscriminator, guildname, guildid string, guildsize int, reason string) *Notification {
	return NewNotification(TypeBan, userid, username, userdiscriminator, guildname, guildid, guildsize, reason)
}

func NewUnbanNotification(userid, username, userdiscriminator, guildname, guildid string, guildsize int, reason string) *Notification {
	return NewNotification(TypeUnban, userid, username, userdiscriminator, guildname, guildid, guildsize, reason)
}

func NewKickNotification(userid, username, userdiscriminator, guildname, guildid string, guildsize int, reason string) *Notification {
	return NewNotification(TypeKick, userid, username, userdiscriminator, guildname, guildid, guildsize, reason)
}

func NewNoticeNotification(userid, username, userdiscriminator, guildname, guildid string, guildsize int, reason string) *Notification {
	return NewNotification(TypeNotice, userid, username, userdiscriminator, guildname, guildid, guildsize, reason)
}
