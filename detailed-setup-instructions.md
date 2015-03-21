### recipebox-server and Go detailed setup instructions

This is a detailed setup guide for setting up recipebox-server with Go.

### Installing Go

First, you should install Go.  Download Go from here: 
[https://golang.org/doc/install](https://golang.org/doc/install).
You should select an installer package based on your operating system,
and it should set up most of the necessary environment variables for you.

Verify that Go is installed by typing `go` at the command line (terminal).

### Setting up your go workspace

All of your go code should reside in a workspace. Please
follow the instructions in this [document](https://golang.org/doc/code.html)
to setup your go workspace and the `$GOPATH` environment variable.
All of your go code will live in your `$GOPATH` directory.

### Get the recipebox-server code

Now, you can download the source code by using the go tool.  Use

    $ go get github.com/harvardd4d/recipebox-server

to get the code.  Note that you can run this command from anywhere;
go will automatically put the project into your workspace under
`$GOROOT/src/github.com/harvardd4d/recipebox-server`.

### Installing godep or other dependencies.

We use godep to manage our dependencies - you may choose to use godep
or manually install the dependencies.  Run

    $ go get github.com/kr/godep

to install godep.  To install a dependency named `github.com/gorilla/pat`, run

    $ go get github.com/gorilla/pat

### Building the program

If you wish to use godep to install dependencies, run

    $ godep go build

in the recipebox-server folder to build the code. This will
produce the `recipebox-server` executable.

If you installed dependencies manually, run

    $ go build

to build the same executable.

### Running the server

Type 

    $ ./recipebox-server

to run the recipebox server.  You'll notice that you need to setup
the DATABASE_URL environment variable.  RecipeBox requires a local
database of recipes.

### Setting up the local database

Currently, recipebox-server uses a postgres database.  You should install
psql for your system from [http://www.postgresql.org/](http://www.postgresql.org/).

After you have installed psql, you should create a database on your
local machine.  Do this with

    $ psql -U postgres

psql will ask you for the password you provided when you installed psql.

After you have logged in, create a local recipes table with (copy and paste)

    CREATE TABLE recipes (
      name text NOT NULL, 
      description text NOT NULL, 
      cuisine integer NOT NULL, 
      mealtype integer NOT NULL, 
      season integer NOT NULL, 
      ingredientlist text NOT NULL, 
      instructions text NOT NULL, 
      id integer PRIMARY KEY NOT NULL, 
      picture bytea
    );

and two sample recipes

    INSERT INTO recipes VALUES (
      'Chinese Broccoli',
      'Lightly flavored Broccoli from the East',
      1,
      1,
      1,
      'Broccoli; Sesame oil',
      'Steam the Broccoli.  Add sesame oil and serve.',
      1,
      NULL
    );

    INSERT INTO recipes VALUES (
      'Toasted Toast',
      'Toasty Toasted Toast',
      1,
      1,
      1,
      'Toast',
      'Toast toast',
      2,
      NULL
    );

### Set your DATABASE_URL environment variable

Now, after you've setup your database, you need to set the DATABASE_URL
environment variable so that recipebox-server will know where to find
the database.  Do this with

    $ export DATABASE_URL=postgres://<username>:<password>@localhost:5432/<database>

at the terminal.  For example, if your username is postgress, password is password,
and you set up a table not in a database, you would do

    $ export DATABASE_URL=postgres://postgres:password@localhost:5432/postgres

By default, any tables you create not in a database are under a database
that has the same name as your username.

### Run recipebox-server

Now that everything is set up, run

    $ ./recipebox-server  

and open your web browser to `localhost:8080/recipes/1` to see
recipe number 1.