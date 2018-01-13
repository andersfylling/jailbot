package notify

import (
	"errors"

	"github.com/andersfylling/jailbot/common"
)

type Subscription struct {
	GuildID   string
	ChannelID string
}

func NewSubscription(gid, cid string) *Subscription {
	return &Subscription{
		GuildID:   gid,
		ChannelID: cid,
	}
}

func SubscribeToNotificationType(nt common.NotificationType, gid, cid string) error {
	err := errors.New("unable to find a topic with given notification type")
	if topic, ok := topics[nt]; ok {
		_, err = topic.Subscribe(gid, cid)
	}

	return err
}

func SubscribeToAll(gid, cid string) error {
	types := []common.NotificationType{
		TypeBan, TypeKick, TypeUnban, TypeNotice,
	}

	for _, t := range types {
		err := SubscribeToNotificationType(t, gid, cid)
		if err != nil {
			return err
		}
	}

	return nil
}
