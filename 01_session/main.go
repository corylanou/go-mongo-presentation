//START IMPORT OMIT
package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

//END IMPORT OMIT

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
		err          error
	)
	//END IMPORT OMIT

	// START DIAL OMIT
	/*
		Format of string for dialing is as follows:

			mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb

	*/
	addr := fmt.Sprintf("mongodb://%s:%s@%s", mongoUser, mongoPassword, mongoUrl)

	// Dial up the server and establish a session
	if mongoSession, err = mgo.Dial(addr); err != nil {
		log.Fatal(err)
	}

	// Make sure the connection closes
	defer mongoSession.Close()
	// END DIAL OMIT

	// START FINAL OMIT
	// This will get the "default" database that the connection string specified
	database = mongoSession.DB("")

	// Get our collection
	collection = database.C(mongoCollection)

	// For debuging, print out the collection we found
	fmt.Printf("Collection: %+v", collection)

	// Close main function
}

// END FINAL OMIT
