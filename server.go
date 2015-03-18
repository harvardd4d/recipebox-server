package main

import (
	"container/list"
	_ "database/sql"
	"flag"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/unrolled/render"
	"net/http"
	"strconv"
	"strings"
)

// the main database.  TODO dependency injection
var recipedb *RecipeDatabase
var renderer *render.Render

// Handler for the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You've reached the recipebox hotline")
}

// Handler for viewing Recipes
func RecipeHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	recipe, err := recipedb.GetRecipe(id)

	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	renderer.HTML(w, http.StatusOK, "recipes/recipe", recipe)
}

// Handler for getting a JSON for a particular recipe
func RecipeJSONHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	recipe, err := recipedb.GetRecipe(id)

	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintf(w, "%v", recipe.ToJSON())
}

// Handler for advanced JSON searches
func RecipeAdvancedJSONHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	strict, err := strconv.Atoi(r.FormValue("strict"))
	name := r.FormValue("name")
	cuisine, _ := strconv.Atoi(r.FormValue("cuisine"))
	season, _ := strconv.Atoi(r.FormValue("season"))
	mealtype, _ := strconv.Atoi(r.FormValue("mealtype"))

	// get all the recipes that match
	var recipes *list.List
	if strict == 0 {
		recipes, err = recipedb.GetRecipesLoose(name, cuisine, mealtype, season)
	} else {
		recipes, err = recipedb.GetRecipesStrict(name, cuisine, mealtype, season)
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
	} else {
		request = err.Error()
	}
	fmt.Fprintf(w, request)
}

func CSSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./webroot/css/pixyll.css")
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	renderer.HTML(w, http.StatusOK, "about", nil)
}

// func ErrorHandler(w http.ResponseWriter, r *http.Request) {
// 	m := make(map[string]string)
// 	m["error_code"] = "404"
// 	m["error_msg"] = "Sorry, we can't seem to find your page."
// 	renderer.HTML(w, http.StatusOK, "error", m)
// }

func main() {
	// Get command line arguments
	portPtr := flag.Int("port", 8080, "the server port number")
	flag.Parse()
	portStr := fmt.Sprintf(":%v", *portPtr)

	// Read database, check to see if we can open
	db_file := "testdb.sqlite"
	db, _ := sqlx.Open("sqlite3", db_file)
	err := db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Unable to open database %v.", db_file))
	}
	recipedb = &RecipeDatabase{DB: db}

	// Set up renderer.  Default template is templates/layout.tmpl
	renderer = render.New(render.Options{
		Layout: "layout",
	})

	// Set up the router
	r := pat.New()
	r.Get("/css/pixyll.css", CSSHandler)
	r.Post("/recipes/jsonsearch", RecipeAdvancedJSONHandler)
	r.Get("/recipes/{id:[0-9]+}/json", RecipeJSONHandler)
	r.Get("/recipes/{id:[0-9]+}", RecipeHandler)
	r.Get("/about", AboutHandler)
	r.Get("/", HomeHandler)

	// Setting up middleware (server, logging layer)
	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(portStr)
}
