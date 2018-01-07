package document

import (
	"errors"

	"github.com/andersfylling/jailbot/database/dbsession"
	"gopkg.in/mgo.v2/bson"
)

type UserDocument struct {
	ID            bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	DiscordID     string        `json:"discordid"`
	Avatar        string        `json:"avatar"`
	Bot           bool          `json:"bot"`
	Discriminator string        `json:"discriminator"`
	Email         string        `json:"email"`
	MFAEnabled    bool          `json:"mfaenabled"`
	Token         string        `json:"token"`
	Username      string        `json:"username"`
	Verified      bool          `json:"verified"`
}

const UserDocumentCollection string = "user"

// Insert the document as a new one into the collection and returns the id
func (c *UserDocument) Insert() (id bson.ObjectId, err error) {
	id = ""
	err = nil

	ses, col, err := dbsession.GetCollection(UserDocumentCollection)
	if err != nil {
		return id, err
	}
	defer ses.Close()

	c.ID = bson.NewObjectId()
	err = col.Insert(c)

	return c.ID, err
}

func (c *UserDocument) Exist() (bool, error) {
	if c.DiscordID == "" {
		return false, errors.New("no discord id set")
	}

	ses, collection, err := dbsession.GetCollection(UserDocumentCollection)
	if err != nil {
		return false, err
	}
	defer ses.Close()

	total, err := collection.Find(bson.M{"discordid": c.DiscordID}).Count()
	return total > 0, err
}

func (c *UserDocument) GetExisting() error {
	if c.DiscordID == "" {
		return errors.New("no discord id set")
	}

	ses, collection, err := dbsession.GetCollection(UserDocumentCollection)
	if err != nil {
		return err
	}
	defer ses.Close()

	collection.Find(bson.M{"discordid": c.DiscordID}).One(&c)

	// check if anything was found
	if c.ID == "" {
		return errors.New("no user with discord id " + c.DiscordID + " found")
	}

	return nil
}
