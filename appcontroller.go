package main

import "net/http"

// Action defines a standard function signature for us to use when creating
// controller actions. A controller action is basically just a method attached to
// a controller.
type Action func(rw http.ResponseWriter, r *http.Request) error

// AppController is a base controller for a web app
type AppController struct{}

// Action helps with error handling in a controller.
// for example, use c.Action(c.Index) to handle the Index page (c.Index needs to
// be defined).
func (c *AppController) Action(a Action) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := a(w, r); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})
}
