package main

import (
	"container/list"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"net/http"
	"os"
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

// // RecipeDisplay is a display recipe object.
// // It's a recipe object + some conveniences, like a
// // parsed list of ingredients
// type RecipeDisplay struct {
// 	*Recipe
// }

// --------------------------------------------
//              HELPER FUNCTIONS
// --------------------------------------------

func PathExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// RenderError uses RBController's renderer to create an error
// based off of a template.
func (c *RBController) RenderError(w http.ResponseWriter, errorCode int, msg string) {
	m := make(map[string]string)
	m["error_code"] = fmt.Sprintf("%v", errorCode)
	m["error_msg"] = msg
	c.HTML(w, errorCode, "error", m)
}

// --------------------------------------------
//              ACTIONS AND HANDLERS
// --------------------------------------------

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

// Contact creates contact page
func (c *RBController) Contact(w http.ResponseWriter, r *http.Request) (err error) {
	c.HTML(w, http.StatusOK, "contact", nil)
	return nil
}

func (c *RBController) EditRecipe(w http.ResponseWriter, r *http.Request) (err error) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	recipe, err := c.GetRecipe(id)
	if err == nil {
		c.HTML(w, http.StatusOK, "recipes/edit", recipe)
	} else if err == sql.ErrNoRows {
		// this means that the recipe wasn't found, so we should return a 404 error
		c.RenderError(w, 404, "Sorry, your page wasn't found")
		err = nil
	}
	return
}

// Home creates the homepage
func (c *RBController) Home(w http.ResponseWriter, r *http.Request) (err error) {
	stats := map[string]string{
		"nRecipes":    "891",
		"nVolunteers": "200,000",
		"nCountries":  "100+",
		"nYears":      "54",
	}
	c.HTML(w, http.StatusOK, "home", stats)
	return nil
}

func (c *RBController) NewRecipe(w http.ResponseWriter, r *http.Request) (err error) {
	c.RenderError(w, 404, "New page coming soon!")
	return nil
}

// Recipe renders a recipe by id
func (c *RBController) Recipe(w http.ResponseWriter, r *http.Request) (err error) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
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
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
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
// TODO: use MUX.  this function currently doesn't work.
func (c *RBController) RecipeJSONAdvanced(w http.ResponseWriter, r *http.Request) (err error) {
	r.ParseForm()
	strict, err := strconv.Atoi(r.PostFormValue("strict"))
	name := r.PostFormValue("name")
	cuisine, _ := strconv.Atoi(r.PostFormValue("cuisine"))
	season, _ := strconv.Atoi(r.PostFormValue("season"))
	mealtype, _ := strconv.Atoi(r.PostFormValue("mealtype"))

	// get all the recipes that match
	var recipes *list.List
	if strict == 0 {
		recipes, err = c.GetRecipesLoose(name, cuisine, mealtype, season)
	} else {
		recipes, err = c.GetRecipesStrict(name, cuisine, mealtype, season)
	}

	// slice of jsons
	jsons := make([]string, recipes.Len())

	if err == nil {
		index := 0
		for e := recipes.Front(); e != nil; e = e.Next() {
			rec := e.Value.(*Recipe)
			jsons[index] = rec.ToJSON()
			index++
		}
		request := strings.Join(jsons, "\n")
		fmt.Fprintf(w, request)
	} else {
		fmt.Fprintf(w, "%v", err.Error())
	}
	return
}

// SaveRecipe takes a POST request from the /recipes/edit/ form
// and saves the recipe back into the database.
func (c *RBController) SaveRecipe(w http.ResponseWriter, r *http.Request) (err error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, _ := strconv.Atoi(idStr)

	name := r.PostFormValue(`name`)
	cuisine, err := strconv.Atoi(r.PostFormValue(`cuisine`))
	mealtype, err1 := strconv.Atoi(r.PostFormValue(`mealtype`))
	season, err2 := strconv.Atoi(r.PostFormValue(`season`))
	ingredients := r.PostFormValue(`ingredients`)
	instructions := r.PostFormValue(`instructions`)

	if err != nil || err1 != nil || err2 != nil {
		fmt.Println("[WARNING] Something went wrong in SaveRecipe")
		c.RenderError(w, 500, "Sorry, something went wrong.")
		return
	}

	// everything OK: build the recipe, and send it to the database
	recipe := Recipe{ID: id, Name: name, Cuisine: cuisine, Mealtype: mealtype,
		Season: season, Ingredientlist: ingredients, Instructions: instructions}
	err = c.RecipeDB.EditRecipe(&recipe)

	if err == nil {
		http.Redirect(w, r, "/recipes/"+idStr+"/", http.StatusFound)
	}
	return
}

// Static serves static pages
func (c *RBController) Static(w http.ResponseWriter, r *http.Request) (err error) {
	vars := mux.Vars(r)
	path := "./webroot/" + vars["path"]

	if PathExists(path) {
		http.ServeFile(w, r, path)
	} else {
		c.RenderError(w, 404, "Sorry, this page was not found.")
	}
	return nil
}
