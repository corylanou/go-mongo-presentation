package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// TodoID is the ID for a Todo task
// START TYPE_ID OMIT
type TodoID string

// END TYPE_ID OMIT

// Todo struct defines a task
// START TYPE OMIT
type Todo struct {
	ID        TodoID    `bson:"_id"          json:"id"`
	Task      string    `bson:"t"            json:"task"`
	Created   time.Time `bson:"c"            json:"created"`
	Updated   time.Time `bson:"u,omitempty"  json:"updated"`
	Completed time.Time `bson:"cp,omitempty" json:"completed"`
}

// END TYPE OMIT

// START BOILERPLATE OMIT
func newID() string {
	return bson.NewObjectId().Hex()
}

// NewTodoID creates a new ID
func NewTodoID() TodoID { return TodoID(newID()) }

// Blank determines if a Todo is blank
func (id TodoID) Blank() bool { return id == "" }

//Present determines if the Todo is present
func (id TodoID) Present() bool { return id != "" }

// Valid determines if a Todo is valid (has an valid ID)
func (id TodoID) Valid() bool { return bson.IsObjectIdHex(string(id)) }

// Invalid determines if the Todo is invalid
func (id TodoID) Invalid() bool { return !id.Valid() }

// END BOILERPLATE OMIT

// GetBSON defines a custom unmarshaller to convert the BSON ObjectID to a TodoID
// START BSON OMIT
func (id TodoID) GetBSON() (v interface{}, e error) {
	if id.Valid() {
		v = bson.ObjectIdHex(string(id))
	}
	return
}

// SetBSON defines a customer marshaller to convert at TodoID to a BSON ObjectID
func (id *TodoID) SetBSON(raw bson.Raw) error {
	var oid bson.ObjectId
	err := raw.Unmarshal(&oid)
	*id = TodoID(oid.Hex())
	return err
}

// END BSON OMIT
