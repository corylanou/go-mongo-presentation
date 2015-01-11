//START IMPORT OMIT
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//END IMPORT OMIT

//START TYPE OMIT
type (
	Todo struct {
		Id        bson.ObjectId `bson:"_id"`
		Task      string        `bson:"t"`
		Created   time.Time     `bson:"c"`
		Updated   time.Time     `bson:"u,omitempty"`
		Completed time.Time     `bson:"cp,omitempty"`
	}
)

//END TYPE OMIT

//START MAIN OMIT
func main() {
	// END MAIN OMIT

	// START ENV OMIT
	var (
		mongoUrl        = os.Getenv("mongo_url")
		mongoUser       = os.Getenv("mongo_user")
		mongoPassword   = os.Getenv("mongo_password")
		mongoCollection = os.Getenv("mongo_collection")
	)
	// END ENV OMIT

	var (
		mongoSession *mgo.Session
		database     *mgo.Database
		collection   *mgo.Collection
		changeInfo   *mgo.ChangeInfo
		err          error
	)
	//END IMPORT OMIT

	// START DIAL OMIT
	/*
		Format of string for dialing is as follows:

			mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb
	*/
	addr := fmt.Sprintf(
		"mongodb://%s:%s@%s",
		mongoUser,
		mongoPassword,
		mongoUrl,
	)

	if mongoSession, err = mgo.Dial(addr); err != nil {
		panic(err)
	}
	defer mongoSession.Close()
	// END DIAL OMIT

	// START FINAL OMIT
	// This will get the "default" database that the connection string specified
	database = mongoSession.DB("")

	// Get our collection
	collection = database.C(mongoCollection)

	var todo = Todo{
		Id:      bson.NewObjectId(),
		Task:    "Demo mgo",
		Created: time.Now(),
	}

	// This is a shortcut to collection.Upsert(bson.M{"_id": todo.id}, &todo)
	if changeInfo, err = collection.UpsertId(todo.Id, &todo); err != nil {
		panic(err)
	}

	spew.Dump(todo)
	spew.Dump(changeInfo)
	// END FINAL OMIT

	// START CHANGE OMIT
	var change = mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": bson.M{
				"cp": time.Now(),
			}}}
	if changeInfo, err = collection.FindId(todo.Id).Apply(change, &todo); err != nil {
		panic(err)
	}
	// END CHANGE OMIT

	spew.Dump(todo)
	spew.Dump(changeInfo)

	// Close main function
}

// END FINAL OMIT
