package main

import (
	"container/list"
	_ "database/sql"
	_ "flag"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	// _ "github.com/mattn/go-sqlite3"
	"github.com/unrolled/render"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Handles the homepage
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You've reached the recipebox hotline")
}

// Handles viewing and rendering recipes for web us3e
func RecipeHandler(recipedb *RecipeDatabase,
	renderer *render.Render) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
		recipe, err := recipedb.GetRecipe(id)

		if err != nil {
			ServeError(renderer, w, 404, err.Error())
			return
		}

		renderer.HTML(w, http.StatusOK, "recipes/recipe", recipe)
	})
}

// Handles retrieving a json for a particular recipe by id
func RecipeJSONHandler(recipedb *RecipeDatabase) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
		recipe, err := recipedb.GetRecipe(id)

		if err != nil {
			fmt.Fprintf(w, "%v", err.Error())
			return
		}

		fmt.Fprintf(w, "%v", recipe.ToJSON())
	})
}

//	Handles advanced JSON searches.
//	Searches are either strict or loose (by name)
//	and are done by season, mealtype, and cuisine.
func RecipeAdvancedJSONHandler(recipedb *RecipeDatabase) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			fmt.Fprintf(w, request)
		} else {
			fmt.Fprintf(w, "%v", err.Error())
		}
	})
}

//	Serves css files
func CSSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./webroot/css/pixyll.css")
}

//	Handles the /about route.
func AboutHandler(renderer *render.Render) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		renderer.HTML(w, http.StatusOK, "about", nil)
	})
}

func ServeError(renderer *render.Render, w http.ResponseWriter,
	code int, msg string) {

	m := make(map[string]string)
	m["error_code"] = fmt.Sprintf("%v", code)
	m["error_msg"] = msg
	renderer.HTML(w, http.StatusOK, "error", m)
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4747"
		fmt.Println("[-] No PORT environment variable detected. Setting to ", port)
	}
	return ":" + port
}

func main() {
	// Get command line arguments
	// portPtr := flag.Int("port", 8080, "the server port number")
	// flag.Parse()
	// portStr := fmt.Sprintf(":%v", *portPtr)

	// Read database, check to see if we can open
	// db_file := "testdb.sqlite"
	// db, _ := sqlx.Open("sqlite3", db_file)

	// Trying to figure out how to connect to heroku
	db_file := os.Getenv("DATABASE_URL")
	connection, _ := pq.ParseURL(db_file)
	connection += " sslmode=verify-full"
	db, _ := sqlx.Open("postgres", db_file)

	err := db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Unable to open database %v.  Error %v", connection, err.Error()))
	}
	recipedb := &RecipeDatabase{DB: db}

	// Set up renderer.  Default template is templates/layout.tmpl
	renderer := render.New(render.Options{
		Layout: "layout",
	})

	// Set up the router
	router := pat.New()
	router.Get("/css/pixyll.css", CSSHandler)
	router.Post("/recipes/jsonsearch", RecipeAdvancedJSONHandler(recipedb))
	router.Get("/recipes/{id:[0-9]+}/json", RecipeJSONHandler(recipedb))
	router.Get("/recipes/{id:[0-9]+}", RecipeHandler(recipedb, renderer))
	router.Get("/about", AboutHandler(renderer))
	router.Get("/", HomeHandler)

	// Setting up middleware (server, logging layer)
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(GetPort())
}
