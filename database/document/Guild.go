package document

import (
	"errors"

	"github.com/andersfylling/jailbot/database/dbsession"
	"gopkg.in/mgo.v2/bson"
)

type GuildDocument struct {
	ID        bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	DiscordID string        `json:"discordid"`
	Icon      string        `json:"icon"`
	OwnerID   string        `json:"ownerid"`
	Name      string        `json:"name"`
	Region    string        `json:"region"`
}

const GuildDocumentCollection string = "guild"

// Insert the document as a new one into the collection and returns the id
func (c *GuildDocument) Insert() (id bson.ObjectId, err error) {
	id = ""
	err = nil

	ses, con, err := dbsession.GetCollection(GuildDocumentCollection)
	if err != nil {
		return id, err
	}
	defer ses.Close()

	c.ID = bson.NewObjectId()
	err = con.Insert(c)

	return c.ID, err
}

func (c *GuildDocument) GetExisting() error {
	if c.DiscordID == "" {
		return errors.New("no discord id set")
	}

	ses, collection, err := dbsession.GetCollection(GuildDocumentCollection)
	if err != nil {
		return err
	}
	defer ses.Close()

	collection.Find(bson.M{"discordid": c.DiscordID}).One(&c)

	// check if anything was found
	if c.ID == "" {
		return errors.New("no guild with discord id " + c.DiscordID + " found")
	}

	return nil
}
