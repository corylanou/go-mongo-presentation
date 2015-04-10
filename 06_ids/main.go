package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// START TYPE_ID OMIT
type TodoId string

// END TYPE_ID OMIT

// START TYPE OMIT
type (
	Todo struct {
		Id        TodoId    `bson:"_id"          json:"id"`
		Task      string    `bson:"t"            json:"task"`
		Created   time.Time `bson:"c"            json:"created"`
		Updated   time.Time `bson:"u,omitempty"  json:"updated"`
		Completed time.Time `bson:"cp,omitempty" json:"completed"`
	}
)

// END TYPE OMIT

// START BOILERPLATE OMIT
func newId() string {
	return bson.NewObjectId().Hex()
}

func NewTodoId() TodoId { return TodoId(newId()) }

func (id TodoId) Blank() bool   { return id == "" }
func (id TodoId) Present() bool { return id != "" }

func (id TodoId) Valid() bool   { return bson.IsObjectIdHex(string(id)) }
func (id TodoId) Invalid() bool { return !id.Valid() }

// END BOILERPLATE OMIT

// START BSON OMIT
func (id TodoId) GetBSON() (v interface{}, e error) {
	if id.Valid() {
		v = bson.ObjectIdHex(string(id))
	}
	return
}

func (id *TodoId) SetBSON(raw bson.Raw) error {
	var oid bson.ObjectId
	err := raw.Unmarshal(&oid)
	*id = TodoId(oid.Hex())
	return err
}

// END BSON OMIT
