package notify

import (
	"errors"
	"fmt"
	"sync"

	"github.com/andersfylling/unison"
)

type Topic struct {
	Type          NotificationType
	Subscribers   []*Subscription
	freePositions []int
	index         int
	sync.RWMutex
}

func NewTopic(nt NotificationType) *Topic {
	t := &Topic{}

	t.Type = nt
	t.Subscribers = []*Subscription{}

	t.freePositions = []int{}
	t.index = 0

	return t
}

func NewBanTopic() *Topic {
	return NewTopic(TypeBan)
}

func (t *Topic) Subscribe(gid, cid string) (*Subscription, error) {
	t.Lock()
	defer t.Unlock()

	sub := NewSubscription(gid, cid)

	// check if there are any available spots
	index, err := t.popAvailablePosition()
	if err != nil { // none available
		t.Subscribers = append(t.Subscribers, sub)
	} else {
		t.Subscribers[index] = sub
	}

	return sub, nil
}

func (t *Topic) Unsubscribe(id string) error {
	t.Lock()
	defer t.Unlock()

	var subscriber *Subscription
	var index int
	for i, sub := range t.Subscribers {
		if sub.GuildID == id {
			subscriber = sub
			index = i
			break
		}
	}

	if subscriber == nil {
		return errors.New("No subscriber with given ID exists: " + id)
	}

	// remove entry
	t.Subscribers[index] = nil
	t.addAvailablePosition(index)

	return nil
}

func (t *Topic) addAvailablePosition(i int) {
	if len(t.freePositions) > t.index {
		t.freePositions[t.index] = i
		t.index++
	} else {
		t.freePositions = append(t.freePositions, i)
		t.index++
	}
}

func (t *Topic) popAvailablePosition() (int, error) {
	if t.index > 0 {
		t.index--
		return t.freePositions[t.index], nil
	}

	return 0, errors.New("No available position to be reused")
}

func (t *Topic) Delete() error {
	// Implementation

	return nil
}

func (t *Topic) Publish(ctx *unison.Context, n *Notification) error {
	t.Lock()
	defer t.Unlock()

	for _, subscriber := range t.Subscribers {

		// make sure the member exist in the subscriber guild
		member, err := ctx.Discord.GuildMember(subscriber.GuildID, n.UserID)
		if err != nil || member == nil {
			continue
		}

		userInfo := fmt.Sprintf("<@%s>, %s#%s{id:%s}", member.User.ID, member.User.Username, member.User.Discriminator, member.User.ID)
		msg := fmt.Sprintf("[%s] By Guild %s. User %s for reason `%s`", n.Type, n.GuildName, userInfo, n.Reason)
		ctx.Discord.ChannelMessageSend(subscriber.ChannelID, msg)
	}

	return nil
}
