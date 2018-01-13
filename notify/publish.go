package notify

import (
	"github.com/andersfylling/jailbot/common"
	"github.com/andersfylling/unison"
	"github.com/sirupsen/logrus"
)

var topics map[common.NotificationType]*Topic

func init() {
	topics = make(map[common.NotificationType]*Topic)
	topics[TypeBan] = NewBanTopic()
	topics[TypeKick] = NewTopic(TypeKick)
	topics[TypeUnban] = NewTopic(TypeUnban)
	topics[TypeNotice] = NewTopic(TypeNotice)
}

func Publish(ctx *unison.Context, nt *Notification) {
	if topic, ok := topics[nt.Type]; ok {
		topic.Publish(ctx, nt)
	} else {
		logrus.Error("[jailbot] Publish recieved invalid notification type")
	}
}
