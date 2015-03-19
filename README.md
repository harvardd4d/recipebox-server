### recipebox-server-go

The RecipeBox web server, written in Go

### Setup

If you're new to Go, set up your GOPATH and Go workspace with
the instructions described in https://golang.org/doc/code.html.  Go
is rather picky about its workspace options but allows for some
very useful modularity as a result of this workspace organization.

`github.com/codegangsta/negroni`

`github.com/gorilla/pat`

`github.com/lib/pq`

`github.com/unrolled/render`

`github.com/jmoiron/sqlx`

The above dependencies are necessary to build this project. 
We are using `godep` to manage dependencies - you may choose to
install the dependencies locally using `go get` or by using
the `godep` tool.  Run 

    $ go get github.com/kr/godep

to install `godep` or a dependency.  

### Building

Run `go build` in the recipebox-go-server folder to compile the code,
assuming you have the dependencies set up manually.  If you are using
`godep`, run `godep go build` to build the code without downloading the
dependencies.

This will create an executable named `recipebox-go-server`.  Run

    $ ./recipebox-go-server

The server will attach itself to port 8080 by default unless
the `PORT` environment variable is set.  In that case, the server
will listen on the value of the PORT environment variable.

The `DATABASE_URL` environment variable is essential to the
operation of the program and should be set to the location
of a database containing the recipes table.

### Expected behavior

The following routes are currently implemented.

1. `GET /recipes/:id` displays the contents of the recipe with specified id.
2. `GET /recipes/:id/json` displays a json string of the recipe with specified id.
3. `POST /recipes/jsonsearch ? strict=[0,1] name=<string> season=<int> mealtype=<int> cuisine=<int>`
searches for recipes that match name, season, mealtype, and
cuisine and returns them as a list of json strings seperated by newline characters.  
A search is either strict or loose.  Strict searches must 
have the name match exactly; weak searches can have the name be a substring.

### Thank you

Thank you to the following projects for your amazing tools! RecipeBox
wouldn't be here without you.

`github.com/codegangsta/negroni`

`github.com/gorilla/pat`

`github.com/lib/pq`

`github.com/unrolled/render`

`github.com/jmoiron/sqlx`

`github.com/johnotander/pixyll`