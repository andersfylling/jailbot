package dbsession

import (
	"errors"
	"log"
	"os"

	"sync"

	"github.com/andersfylling/jailbot/jailbotconst"
	mgo "gopkg.in/mgo.v2"
)

type singleton struct {
	session *mgo.Session
}

var instance *singleton
var database string
var connectionSetupIssue error
var once sync.Once

// https://console.bluemix.net/docs/services/ComposeForMongoDB/connecting-external.html#connecting-external-app
func setupConnection() {
	instance = &singleton{}
	uri := os.Getenv(jailbotconst.EnvMongoDBURL)
	if uri == "" {
		connectionSetupIssue = errors.New("No connection string provided - set env " + jailbotconst.EnvMongoDBURL)
		log.Fatal(connectionSetupIssue)
		return
	}

	// status
	//fmt.Println("Establishing connection to " + uri)
	ses, err := mgo.Dial(uri)
	if err != nil {
		connectionSetupIssue = err
		return
	}

	// check that its actually Connected
	// TODO

	// set synchronous session
	ses.SetMode(mgo.Monotonic, true)

	// store the session to the singleton instance
	instance.session = ses
}

// GetInstance singleton pattern
func GetInstance() (*mgo.Session, error) {
	once.Do(setupConnection)

	instance.session.Refresh() // make sure the session is alive

	return instance.session, connectionSetupIssue
}

func GetCopy() (*mgo.Session, error) {
	session, err := GetInstance()
	if err != nil {
		return nil, err
	}

	return session.Copy(), nil
}

// GetCollection Returns a collection based on a cloned session
func GetCollection(collection string) (*mgo.Session, *mgo.Collection, error) {
	session, err := GetInstance()
	if err != nil {
		return nil, nil, err
	}

	con := session.Clone()

	return con, con.DB(database).C(collection), err
}
