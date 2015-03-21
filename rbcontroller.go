package main

import (
	"container/list"
	"database/sql"
	"fmt"
	"github.com/unrolled/render"
	"net/http"
	"strconv"
	"strings"
)

// RBController is the RecipeBox controller object.
// The controller is responsible for querying the database
// as well as providing http.HandleFunc to handle URL requests.
type RBController struct {
	AppController
	*RecipeDB
	*render.Render
}

// RenderError uses RBController's renderer to create an error
// based off of a template.
func (c *RBController) RenderError(w http.ResponseWriter, errorCode int, msg string) {
	m := make(map[string]string)
	m["error_code"] = fmt.Sprintf("%v", errorCode)
	m["error_msg"] = msg
	c.HTML(w, errorCode, "error", m)
}

// About creates the about page
func (c *RBController) About(w http.ResponseWriter, r *http.Request) (err error) {
	c.HTML(w, http.StatusOK, "about", nil)
	return nil
}

// Action helps with error handling in a controller.
// Overriding the AppController errors to make use of the renderer
func (c *RBController) Action(a Action) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := a(w, r); err != nil {
			c.RenderError(w, http.StatusInternalServerError,
				"Internal server error\n"+err.Error())
		}
	})
}

// CSS serves css files
func (c *RBController) CSS(w http.ResponseWriter, r *http.Request) (err error) {
	http.ServeFile(w, r, "./webroot/css/pixyll.css")
	return nil
}

// Home creates the homepage
func (c *RBController) Home(w http.ResponseWriter, r *http.Request) (err error) {
	fmt.Fprintf(w, "You've reached the recipebox hotline")
	return nil
}

// Recipe renders a recipe by id
func (c *RBController) Recipe(w http.ResponseWriter, r *http.Request) (err error) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	recipe, err := c.GetRecipe(id)
	if err == nil {
		c.HTML(w, http.StatusOK, "recipes/recipe", recipe)
	} else if err == sql.ErrNoRows {
		// this means that the recipe wasn't found, so we should return a 404 error
		c.RenderError(w, 404, "Sorry, your page wasn't found")
		err = nil
	}
	return
}

// RecipeJSON renders a raw JSON string of a recipe selected by id
func (c *RBController) RecipeJSON(w http.ResponseWriter, r *http.Request) (err error) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	recipe, err := c.GetRecipe(id)
	if err == nil {
		c.JSON(w, http.StatusOK, recipe)
	} else if err == sql.ErrNoRows {
		c.RenderError(w, 404, "Sorry, your page wasn't found")
		err = nil
	}
	return
}

// RecipeJSONAdvanced handles advanced JSON searches.
// Searches are either strict or loose (by name)
// and are done by season, mealtype, and cuisine.
func (c *RBController) RecipeJSONAdvanced(w http.ResponseWriter, r *http.Request) (err error) {
	r.ParseForm()
	strict, err := strconv.Atoi(r.FormValue("strict"))
	name := r.FormValue("name")
	cuisine, _ := strconv.Atoi(r.FormValue("cuisine"))
	season, _ := strconv.Atoi(r.FormValue("season"))
	mealtype, _ := strconv.Atoi(r.FormValue("mealtype"))

	// get all the recipes that match
	var recipes *list.List
	if strict == 0 {
		recipes, err = c.GetRecipesLoose(name, cuisine, mealtype, season)
	} else {
		recipes, err = c.GetRecipesStrict(name, cuisine, mealtype, season)
	}

	// slice of jsons
	jsons := make([]string, recipes.Len())
	request := ""

	if err == nil {
		index := 0
		for e := recipes.Front(); e != nil; e = e.Next() {
			rec := e.Value.(*Recipe)
			jsons[index] = rec.ToJSON()
			index++
		}
		request = strings.Join(jsons, "\n")
		fmt.Fprintf(w, request)
	} else {
		fmt.Fprintf(w, "%v", err.Error())
	}
	return
}
