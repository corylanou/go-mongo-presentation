package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//START TYPE OMIT
type (
	Todo struct {
		Id        bson.ObjectId `bson:"_id"`
		Task      string        `bson:"t"`
		Created   time.Time     `bson:"c"`
		Updated   time.Time     `bson:"u,omitempty"`
		Due       time.Time     `bson:"d,omitempty"` // Add due field for example
		Completed time.Time     `bson:"cp,omitempty"`
	}

	// Define a structure to return our aggretation results to
	TodoDueCounts struct {
		Id    time.Time `bson:"_id"`
		Count int       `bson:"count"`
	}
)

//END TYPE OMIT

func main() {
	var (
		mongoURL        = os.Getenv("mongo_url")
		mongoUser       = os.Getenv("mongo_user")
		mongoPassword   = os.Getenv("mongo_password")
		mongoCollection = os.Getenv("mongo_collection")
	)

	var (
		mongoSession *mgo.Session
		database     *mgo.Database
		collection   *mgo.Collection
		err          error
	)

	addr := fmt.Sprintf(
		"mongodb://%s:%s@%s",
		mongoUser,
		mongoPassword,
		mongoURL,
	)

	if mongoSession, err = mgo.Dial(addr); err != nil {
		log.Fatal(err)
	}

	// This will get the "default" database that the connection string specified
	database = mongoSession.DB("")

	// Get our collection
	collection = database.C(mongoCollection)

	// START SEED OMIT
	// declare a slice of todos to insert to
	var todos []Todo

	// Define our times
	now := time.Now()
	tomorrow := now.Add(time.Hour * 24)

	//Create a quick helper for readability
	newId := bson.NewObjectId

	// Create seed data
	todos = append(todos, Todo{Id: newId(), Task: "First task", Created: now, Due: tomorrow})
	todos = append(todos, Todo{Id: newId(), Task: "Second task", Created: now, Due: now})
	todos = append(todos, Todo{Id: newId(), Task: "Third task", Created: now, Due: now})
	todos = append(todos, Todo{Id: newId(), Task: "Fourth task", Created: now, Due: now})
	todos = append(todos, Todo{Id: newId(), Task: "Fifth task", Created: now, Due: now})

	// Upsert seed data
	for _, todo := range todos {
		if _, err = collection.UpsertId(todo.Id, &todo); err != nil {
			log.Fatal(err)
		}
	}
	// END SEED OMIT

	// START PIPELINE OMIT
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":   bson.M{"$dayOfYear": "$d"},
			"count": bson.M{"$sum": 1},
		}},
	}
	// END PIPELINE OMIT

	// START RESULTS OMIT
	var (
		result  TodoDueCounts
		results []TodoDueCounts
	)
	// END RESULTS OMIT

	// START ITER OMIT
	iter := collection.Pipe(pipeline).Iter()
	for {
		if iter.Next(&result) {
			results = append(results, result)
		} else {
			break
		}
	}
	err = iter.Err()
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(results)

	// END ITER OMIT
}
