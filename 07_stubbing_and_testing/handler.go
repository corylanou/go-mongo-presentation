package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// START TYPE OMIT
type Server struct {
	Database Database
}

type Database interface {
	TodoFind(id string) (Todo, error)
	TodoUpsert(todo *Todo, returnNew bool) error
}

// END TYPE OMIT

// START HANDLER1 OMIT
func (s *Server) TodoShow(w http.ResponseWriter, r *http.Request) {
	todoId := r.FormValue("todoId")
	todo, err := s.Database.TodoFind(todoId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusInternalServerError, Text: http.StatusText(http.StatusInternalServerError)}); err != nil {
			panic(err)
		}
		return
	}

	if !todo.ID.Blank() {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			panic(err)
		}
		return
	}

	// END HANDLER1 OMIT
	// START HANDLER2 OMIT
	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: http.StatusText(http.StatusNotFound)}); err != nil {
		panic(err)
	}
}

// END HANDLER2 OMIT

// TodoID is the ID for a Todo task
type TodoID string

// Todo struct defines a task
type Todo struct {
	ID        TodoID    `bson:"_id"          json:"id"`
	Task      string    `bson:"t"            json:"task"`
	Created   time.Time `bson:"c"            json:"created"`
	Updated   time.Time `bson:"u,omitempty"  json:"updated"`
	Completed time.Time `bson:"cp,omitempty" json:"completed"`
}

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

// GetBSON defines a custom unmarshaller to convert the BSON ObjectID to a TodoID
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

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}
