package handler_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/corylanou/go-mongo-presentation/07_stubbing_and_testing"
)

// START STUB OMIT

type DatabaseStub struct {
	todoFindFunc   func(id string) (handler.Todo, error)
	todoUpsertFunc func(todo *handler.Todo, returnNew bool) error
}

func (d *DatabaseStub) TodoFind(id string) (handler.Todo, error) {
	return d.todoFindFunc(id)
}

func (d *DatabaseStub) TodoUpsert(todo *handler.Todo, returnNew bool) error {
	return d.todoUpsertFunc(todo, returnNew)
}

// END STUB OMIT

// START FOUND1 OMIT
func Test_TodoShow_Found(t *testing.T) {
	id := bson.NewObjectId().Hex()
	dm := &DatabaseStub{}
	dm.todoFindFunc = func(id string) (handler.Todo, error) {
		return handler.Todo{ID: handler.TodoID(id)}, nil
	}
	server := handler.Server{
		Database: dm,
	}
	// END FOUND1 OMIT
	// START FOUND2 OMIT

	req, err := http.NewRequest("GET", "http://example.com/todos?todoId="+id, nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	server.TodoShow(w, req)

	t.Logf("%d - %s", w.Code, w.Body.String())
	if w.Code != http.StatusOK {
		t.Fatalf("Unexpected error code.  exp: %d, got: %d", http.StatusOK, w.Code)
	}
	got := strings.TrimSpace(w.Body.String())
	exp := fmt.Sprintf(`{"id":"%s","task":"","created":"0001-01-01T00:00:00Z","updated":"0001-01-01T00:00:00Z","completed":"0001-01-01T00:00:00Z"}`, id)
	if exp != got {
		t.Fatalf("Unexpected body.  \nexp: %q, \ngot: %q", exp, got)
	}
}

// END FOUND2 OMIT

func Test_TodoShow_NotFound(t *testing.T) {
	dm := &DatabaseStub{}
	// START NOTFOUND1 OMIT
	dm.todoFindFunc = func(id string) (handler.Todo, error) {
		return handler.Todo{}, nil
	}
	// END NOTFOUND1 OMIT
	server := handler.Server{
		Database: dm,
	}

	req, err := http.NewRequest("GET", "http://example.com/todos?todoId=foo", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	server.TodoShow(w, req)

	t.Logf("%d - %s", w.Code, w.Body.String())

	// START NOTFOUND2 OMIT
	if w.Code != http.StatusNotFound {
		t.Fatalf("Unexpected error code.  exp: %d, got: %d", http.StatusNotFound, w.Code)
	}
	if got, exp := strings.TrimSpace(w.Body.String()), `{"code":404,"text":"Not Found"}`; exp != got {
		t.Fatalf("Unexpected body.  \nexp: %q, \ngot: %q", exp, got)
	}
	// END NOTFOUND2 OMIT
}

func Test_TodoShow_ServerError(t *testing.T) {
	dm := &DatabaseStub{}
	// START SERVERERROR1 OMIT
	dm.todoFindFunc = func(id string) (handler.Todo, error) {
		return handler.Todo{}, fmt.Errorf("I blew up!")
	}
	// END SERVERERROR1 OMIT
	server := handler.Server{
		Database: dm,
	}

	req, err := http.NewRequest("GET", "http://example.com/todos?todoId=foo", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	server.TodoShow(w, req)

	t.Logf("%d - %s", w.Code, w.Body.String())
	// START SERVERERROR2 OMIT
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Unexpected error code.  exp: %d, got: %d", http.StatusInternalServerError, w.Code)
	}
	if got, exp := strings.TrimSpace(w.Body.String()), `{"code":500,"text":"Internal Server Error"}`; exp != got {
		t.Fatalf("Unexpected body.  \nexp: %q, \ngot: %q", exp, got)
	}
	// END SERVERERROR2 OMIT
}
