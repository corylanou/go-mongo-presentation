// START IMPORT OMIT
package main

import (
	"time"

	"github.com/davecgh/go-spew/spew"
)

// END IMPORT OMIT

// START TYPE OMIT
type (
	Todo struct {
		Task      string    `bson:"t"`
		Created   time.Time `bson:"c"`
		Updated   time.Time `bson:"u,omitempty"`
		Completed time.Time `bson:"cp,omitempty"`
	}
)

// END TYPE OMIT

// START MAIN OMIT
func main() {
	var todo = Todo{
		Task:    "Demo mgo",
		Created: time.Now(),
	}
	spew.Dump(todo)
}

// END MAIN OMIT
