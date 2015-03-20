package main

import (
	"container/list"
	"fmt"
	"github.com/unrolled/render"
	"net/http"
	"strconv"
	"strings"
)

// MyController is the RecipeBox controller object.
// The controller is responsible for querying the database
// as well as providing http.HandleFunc to handle URL requests.
type MyController struct {
	AppController
	*RecipeDB
	*render.Render
}

// About creates the about page
func (c *MyController) About(w http.ResponseWriter, r *http.Request) (err error) {
	c.HTML(w, http.StatusOK, "about", nil)
	return nil
}

// Action helps with error handling in a controller.
// Overriding the AppController errors to make use of the renderer
func (c *MyController) Action(a Action) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := a(w, r); err != nil {
			m := make(map[string]string)
			m["error_code"] = "500"
			m["error_msg"] = err.Error()
			c.HTML(w, http.StatusOK, "error", m)
		}
	})
}

// CSS serves css files
func (c *MyController) CSS(w http.ResponseWriter, r *http.Request) (err error) {
	http.ServeFile(w, r, "./webroot/css/pixyll.css")
	return nil
}

// Home creates the homepage
func (c *MyController) Home(w http.ResponseWriter, r *http.Request) (err error) {
	fmt.Fprintf(w, "You've reached the recipebox hotline")
	return nil
}

// Recipe renders a recipe by id
func (c *MyController) Recipe(w http.ResponseWriter, r *http.Request) (err error) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	recipe, err := c.GetRecipe(id)
	if err == nil {
		c.HTML(w, http.StatusOK, "recipes/recipe", recipe)
	}
	return
}

// RecipeJSON renders a raw JSON string of a recipe selected by id
func (c *MyController) RecipeJSON(w http.ResponseWriter, r *http.Request) (err error) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	recipe, err := c.GetRecipe(id)
	if err != nil {
		fmt.Fprintf(w, "%v", err.Error())
	} else {
		fmt.Fprintf(w, "%v", recipe.ToJSON())
	}
	return
}

// RecipeJSONAdvanced handles advanced JSON searches.
// Searches are either strict or loose (by name)
// and are done by season, mealtype, and cuisine.
func (c *MyController) RecipeJSONAdvanced(w http.ResponseWriter, r *http.Request) (err error) {
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
