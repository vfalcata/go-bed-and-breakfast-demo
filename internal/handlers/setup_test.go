package handlers

import (
	"encoding/gob"
	"fmt"
	"log"
	"myapp/internal/config"
	"myapp/internal/models"
	"myapp/internal/render"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplate = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {

	//////// FROM MAIN START //////////
	gob.Register(models.Reservation{}) // We MUST register the type for the session data or it will not know.
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

	// for testing, we hard coded our paths to each template. Those template paths are relative to the root of the project so they will break during testing because we are now relative to this root folder
	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create the template cache")
	}
	app.TemplateCache = tc
	app.UseCache = true // set this to true so template pathing does not break by using "render.go", recall that the cache is already generated in the mains

	repo := NewRepo(&app)
	NewHandlers(repo)

	render.NewRenderer(&app) //

	///////// FROM MAIN END ///////////

	///////// FROM ROUTES START /////////////
	mux := chi.NewRouter()

	// Middleware stuff is in the main pkg, so we cannot import those so instead we will just copy and paste the functions
	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf) // No need to test this.
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))

	// we now need to use the file server
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer)) //we strip out the path so we get the file path for the static file
	return mux
	//////////// FROM ROUTES END ////////////////

}

// ////////////// MIDDLEWARE FUNC START ///////////////
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",              // "/" is how you refer to the entire site for a cookie
		Secure:   app.InProduction, // good enough for dev mode we will change for prod
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// Webservers by their very nature are not state aware we need to add middleware that tells this web server (our web application) that it should remember state using sessions
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

///////////////////// MIDDLEWARE FUNC END /////////////////

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)
	// syntactically simpler way to do above
	myCache := map[string]*template.Template{} // instantiate empty map

	//create entire cache at once
	// One thing to bear in mind, when you are rendering a template that uses a layout, the FIRST thing MUST be the page you want to render, then after that the associated layouts and partials
	// this means we start parsing templates and adding to the cache we want to add anything that has suffix "*.pag.tmpl" first
	// thus we need to go to the folder that holds it all and everything that ends with "*.pag.tmpl" first into a data structure we can loop through

	// get all of the files with suffix "*.page.tmpl" from the "./templates" folder
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplate)) // go to some location and look for these files

	if err != nil {
		return myCache, err
	}

	// range through all the files ending with *.page.tmpl
	for _, page := range pages {
		// page variable will give the full path, whereas we only want the relative path to the root directory of our project
		name := filepath.Base(page) // Base will return last element of the path, which is our filename

		// 1. First we parse our templates for the currend page
		// ".New(name)" populates the template with a name. the parsed file gets stored in the template called "name"
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		// ts for template set, at this point it will only contain the page we are on

		if err != nil {
			return myCache, err
		}
		// 2. Then it looks for layouts and colleds the filepaths
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))

		if err != nil {
			return myCache, err
		}

		// 3. then it adds everything from 1. and 2. together with the ParseGlob function
		// if we have layouts, then we have to do something with them
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
			// ParseGlob
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts // store the template set to the cache
	}
	return myCache, nil

}
