### recipebox-server

The RecipeBox web server, written in Go.  Get the code with

    $ go get github.com/harvardd4d/recipebox-server

or, if you are new to go, read `detailed-setup-instructions.md` in this
directory. View a live preview at 
[pc-recipebox.herokuapp.com](pc-recipebox.herokuapp.com).

### Setup

If you're new to Go, set up your GOPATH and Go workspace with
the instructions described in https://golang.org/doc/code.html.  Go
is rather picky about its workspace options but allows for some
very useful modularity as a result of this workspace organization.

- [github.com/codegansta/negroni](github.com/codegangsta/negroni)
- [github.com/gorilla/pat](github.com/gorilla/pat)
- [github.com/lib/pq](github.com/lib/pq)
- [github.com/unrolled/render](github.com/unrolled/render)
- [github.com/jmoiron/sqlx](github.com/jmoiron/sqlx)

The above dependencies are necessary to build this project. 
We are using `godep` to manage dependencies - you may choose to
install the dependencies locally using `go get` or by using
the `godep` tool.  Run 

    $ go get github.com/kr/godep

to install `godep` or a dependency.  

### Setting up a database

The recipebox-server app requires a recipes database. recipebox-server
will read in the `DATABASE_URL` environment variable and attempt to
connect to the database listed there.  The
server is set up to use a postgres database. Please install `psql` and look at
`sample_sql.txt` to setup a local postgres database.

### Building

Run `go build` in the recipebox-server folder to compile the code,
assuming you have the dependencies set up manually.  If you are using
`godep`, run `godep go build` to build the code without downloading the
dependencies.

This will create an executable named `recipebox-server`.  Run

    $ ./recipebox-server

The server will attach itself to port 8080 by default unless
the `PORT` environment variable is set.  In that case, the server
will listen on the value of the PORT environment variable.

The `DATABASE_URL` environment variable is essential to the
operation of the program and should be set to the location
of a database containing the recipes table.

### Testing

Run `go test` to test. Current, will test the server without the existence
of a postgres database.

### Expected behavior

The following routes are currently implemented.

1. `GET /recipes/:id` displays the contents of the recipe with specified id.
2. `GET /recipes/:id/json` displays a json string of the recipe with specified id.
3. `POST /recipes/jsonsearch ? strict=[0,1] name=<string> season=<int> mealtype=<int> cuisine=<int>`
searches for recipes that match name, season, mealtype, and
cuisine and returns them as a list of json strings seperated by newline characters.  
A search is either strict or loose.  Strict searches must 
have the name match exactly; weak searches can have the name be a substring.
4. `GET /about` displays about text.

### Code details

The directory is set up as so:

    recipebox-server/
    |
    |-- recipe.go (Recipe and RecipeDB types and functions)
    |-- appcontroller.go (Generic app controller type)
    |-- rbcontroller.go (RecipeBox app controller)
    |-- server.go (RecipeBox server)
    |
    +-- webroot
    |   +-- css
    |       |-- pixyll.css
    | 
    +-- templates
        +-- recipes
        |   |-- recipe.tmpl (recipe view template)
        |
        |-- layout.tmpl (layout template)
        |-- error.tmpl (error template)
        |-- about.tmpl (about template)

`server.go` handles routing.  `rbcontroller.go` is the RecipeBox
controller and handles rendering of html templates and RecipeDB querying.

### Todo

- More testing
- Move from gorilla/pat to gorilla/mux
- Recipes home page
- Recipes search by category
- User login (Google Authentication)
- Some way for users to keep track of recipes
- Saving recipes to database

### Thank you

Thank you to the following projects for your amazing tools! RecipeBox
wouldn't be here without you.

- [codegansta/negroni](github.com/codegangsta/negroni)
- [gorilla/pat](github.com/gorilla/pat)
- [lib/pq](github.com/lib/pq)
- [unrolled/render](github.com/unrolled/render)
- [jmoiron/sqlx](github.com/jmoiron/sqlx)
- [johnotander](github.com/johnotander/pixyll)