package main

import (
	_ "database/sql"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/unrolled/render"
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

	// Set up renderer.  Default template is templates/layout.tmpl
	renderer := render.New(render.Options{
		Layout: "layout",
	})

	// Set up the controller. The controller is responsible for
	// rendering, database queries, and handling requests
	c := &MyController{Render: renderer, RecipeDB: recipedb}

	// Set up the router and associate routes with the controller
	router := pat.New()
	router.Get("/css/pixyll.css", c.Action(c.CSS))
	router.Post("/recipes/jsonsearch", c.Action(c.RecipeJSONAdvanced))
	router.Get("/recipes/{id:[0-9]+}/json", c.Action(c.RecipeJSON))
	router.Get("/recipes/{id:[0-9]+}", c.Action(c.Recipe))
	router.Get("/about", c.Action(c.About))
	router.Get("/", c.Action(c.Home))

	// Setting up middleware (server, logging layer)
	n := negroni.Classic()
	n.UseHandler(router)

	// Run on specified port
	n.Run(GetPort())
}
