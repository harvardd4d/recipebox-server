package main

import (
	_ "database/sql"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/unrolled/render"
	"html/template"
	"os"
)

// GetPort retrieves the port number set in the PORT environment variable.
// If PORT does not exist, GetPort will return 8080.
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("[recipebox] No PORT environment variable detected. Using port",
			port)
	}
	return ":" + port
}

// ConnectToDB connects to a postgres database.
// The database path should be stored in the DATABASE_URL environment var
func ConnectToDB() (recipedb *RecipeDB) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		panic("[recipebox] DATABASE_URL environment variable not set. Please see README.")
	}
	connection, _ := pq.ParseURL(dbURL)
	// connection += " sslmode=disable"

	// first, open database with sslmode=verify-full
	db, _ := sqlx.Open("postgres", connection+" sslmode=verify-full")
	err := db.Ping()
	if err != nil {
		fmt.Printf("[recipebox] Unable to open database, will retry with sslmode=disable.  Error %v\n",
			err.Error())

		// next, open postgres database with sslmode=disable. Only for local db.
		db, _ = sqlx.Open("postgres", connection+" sslmode=disable")
		err = db.Ping()
		if err != nil {
			panic(fmt.Sprintf("[recipebox] Unable to open database %v.  Error %v", connection, err.Error()))
		}
	}

	fmt.Println("[recipebox] Recipes database opened successfully.")
	recipedb = &RecipeDB{DB: db}
	return
}

func main() {
	// Connect to a database, get a *RecipeDB object
	recipedb := ConnectToDB()

	// Some helper functions for our renderer
	recipesHelper := template.FuncMap{
		"ParseIngredients": ParseIngredients,
		"ParseMeal":        ParseMealtype,
		"ParseSeason":      ParseSeason,
	}

	// Set up renderer.  Default template is templates/layout.tmpl
	renderer := render.New(render.Options{
		Layout: "layout",
		Funcs: []template.FuncMap{
			recipesHelper,
		},
	})

	// Set up the controller. The controller is responsible for
	// rendering, database queries, and handling requests
	c := &RBController{Render: renderer, RecipeDB: recipedb}

	// Set up the router and associate routes with the controller
	router := mux.NewRouter()
	router.HandleFunc("/recipes/jsonsearch/", c.Action(c.RecipeJSONAdvanced))
	router.HandleFunc("/recipes/{id:[0-9]+}/json/", c.Action(c.RecipeJSON))
	router.HandleFunc("/recipes/{id:[0-9]+}/edit/", c.Action(c.EditRecipe))
	router.HandleFunc("/recipes/{id:[0-9]+}/save/", c.Action(c.SaveRecipe))
	router.HandleFunc("/recipes/{id:[0-9]+}/", c.Action(c.Recipe))
	router.HandleFunc("/recipes/new/save/", c.Action(c.SaveRecipe))
	router.HandleFunc("/recipes/new/", c.Action(c.NewRecipe))
	router.HandleFunc("/about/", c.Action(c.About))
	router.HandleFunc("/contact/", c.Action(c.Contact))
	router.HandleFunc("/index/", c.Action(c.Home))
	router.HandleFunc("/", c.Action(c.Home))
	router.HandleFunc("/{path:.+}", c.Action(c.Static))

	// Setting up middleware (server, logging layer)
	n := negroni.Classic()
	n.UseHandler(router)

	// Run on specified port
	n.Run(GetPort())
}
