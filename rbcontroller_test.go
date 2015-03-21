package main

import (
	"errors"
	"fmt"
	"github.com/unrolled/render"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockAction is a pretend action used to test Action
func (c *RBController) MockAction(err error) Action {
	return Action(func(w http.ResponseWriter, r *http.Request) error {
		if err == nil {
			fmt.Fprintf(w, "Hello, world")
		}
		return err
	})
}

// TestAction tests the Action method, which is responsible
// for producing a http.HandlerFunc based on the action
func TestAction(t *testing.T) {
	// setup controller and renderer
	renderer := render.New(render.Options{
		Layout: "layout",
	})
	c := &RBController{Render: renderer, RecipeDB: nil}

	// setup request, aboutHandler
	req, _ := http.NewRequest("GET", "", nil)

	// test with no error
	w := httptest.NewRecorder()
	actionHandler := c.Action(c.MockAction(nil))
	actionHandler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Action didn't return %v, expected %v", http.StatusOK,
			w.Code)
	}

	// test with error
	w = httptest.NewRecorder()
	actionHandler = c.Action(c.MockAction(errors.New("MockError")))
	actionHandler.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Action didn't return %v, got %v", http.StatusInternalServerError,
			w.Code)
	}
}

// TestAbout tests the About action, which should display the
// About webpage.
func TestAbout(t *testing.T) {
	// setup controller
	renderer := render.New(render.Options{
		Layout: "layout",
	})
	c := &RBController{Render: renderer, RecipeDB: nil}

	// setup request, aboutHandler
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	aboutHandler := c.Action(c.About)
	aboutHandler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("About page didn't return %v", http.StatusOK)
	}
}
