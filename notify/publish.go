package notify

import "github.com/andersfylling/unison"

var topics map[NotificationType]*Topic

func init() {
	topics = make(map[NotificationType]*Topic)
	topics[TypeBan] = NewBanTopic()
}

func Publish(ctx *unison.Context, nt *Notification) {
	if topic, ok := topics[nt.Type]; ok {
		topic.Publish(ctx, nt)
	}
}
