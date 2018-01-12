package document

import (
	"time"

	"github.com/andersfylling/jailbot/database/dbsession"
	"github.com/andersfylling/jailbot/notify"
	"gopkg.in/mgo.v2/bson"
)

type EventDocument struct {
	ID             bson.ObjectId           `json:"_id" bson:"_id,omitempty"`
	GuildID        string                  `json:"guildid"` // discord guild id
	UserID         string                  `json:"userid"`  // discord user id
	Type           notify.NotificationType `json:"type"`
	Reason         string                  `json:"reason"`
	BanRemoved     bool                    `json:"banremoved"`
	BanRemovedDate time.Time               `json:"banremoveddate"`
}

const EventDocumentCollection string = "event"

// Insert the document as a new one into the collection and returns the id
func (c *EventDocument) Insert() (id bson.ObjectId, err error) {
	id = ""
	err = nil

	ses, con, err := dbsession.GetCollection(EventDocumentCollection)
	if err != nil {
		return id, err
	}
	defer ses.Close()

	c.ID = bson.NewObjectId()
	err = con.Insert(c)

	return c.ID, err
}
