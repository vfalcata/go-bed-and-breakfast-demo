package main // packages always live on the same directory

// To run:
// go run .\cmd\web\.

import (
	"encoding/gob"
	"fmt"
	"log"
	"myapp/internal/config"
	"myapp/internal/driver"
	"myapp/internal/handlers"
	"myapp/internal/helpers"
	"myapp/internal/models"
	"myapp/internal/render"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080" // we dont want it to change so we make it const, so it cannot ever be changed

// IMPORTANT QUICK NOTE
// When a function name or a variable name begins with a CAPITAL letter it is ACCESSIBLE outside a given package
// When a func starts with a lower case letter it is private in an oop sense, thus it is not visible outside this package
func addValues(x, y int) int {
	return x + y
}

var app config.AppConfig // no available to all main package files
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close() // we need this here and NOT in the run() function because we want it to run AFTER the main is done, aka the application stops running

	//// OLD CODE START ////
	//requires argument pathname url that we want to listen to
	/*
		http.HandleFunc("/", handlers.Repo.Home)
		http.HandleFunc("/about", handlers.Repo.About)
	*/
	//// OLD CODE END ////
	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	//// OLD CODE START ////
	// This one liner starts the server
	/*
		_ = http.ListenAndServe(portNumber, nil) //syntax to listen to a port Dont forget the COLON
	*/
	// we dont need a handler, because we already defined it above

	// This is the basis of every web app
	//// OLD CODE END ////

	// first we specify the server
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	// now we start the server
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	gob.Register(models.Reservation{}) // We MUST register the type for the session data or it will not know.
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	app.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	/* session := scs.New() */ //variable Shadowing, this session is not the same as our global outside scope one
	session = scs.New()        // this is how we assign the outside scope variable
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true                  // should the cookie persist after the user closes the window or their web browser
	session.Cookie.SameSite = http.SameSiteLaxMode // how strict do you want to be about what site this cookie applies to lax mode is the default
	session.Cookie.Secure = app.InProduction       // this will insis that the cookies be encrypted and that the connection is from https and not http, we set it to false right now for dev mode, but in prod we want it true

	app.Session = session

	// Connect to postgres DB
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=password")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying")
	}
	log.Println("Successfully connected to database!")

	// we want to load our template cache here instead of the render
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create the template cache")
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false // will rebuild the page on every request

	repo := handlers.NewRepo(&app, db) // now a repo is created that can be used in the handlers
	// NOTE db is NOT a database connection pool tied to a specific database itself, instead it is a pointer to a driver that can, at this moment only handle postgres, but can be easily
	handlers.NewHandlers(repo)

	render.NewRenderer(&app) //
	helpers.NewHelpers(&app) //give access helpers access to app
	return db, nil
}
