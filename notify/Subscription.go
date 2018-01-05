package notify

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
