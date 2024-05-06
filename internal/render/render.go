package render

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"myapp/internal/config"
	"myapp/internal/models"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/justinas/nosurf"
)

// Renders template utilizing html templates
func RenderTemplateV1(w http.ResponseWriter, tmpl string) {

	// when you use a template file you need to parse it
	// this will re render the page each time, we need to build a template cache
	parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl,
		"./templates/base.layout.tmpl")
	// If we have more mandatory files, we just add that one line here
	// If the current page USES one of these files, then it will be parsed
	// If the page DOES NOT use one of these files that it will simply be ignored

	// "./" is the root directory

	err := parsedTemplate.Execute(w, nil) // w is the writer we want to write with

	if err != nil {
		fmt.Println("error parsing template: ", err)
	}

}

// This is a PACKAGE LEVEL VARIABLE that will exist for the life of our program
// need to access the template cache "tc" outside of this file to get performance savings
/*
var tc = make(map[string]*template.Template) // the value *template.Template is used here becuase theat is what template.ParseFiles returns
*/
func RenderTemplateV2(w http.ResponseWriter, t string) {
	var err error
	var tc = make(map[string]*template.Template) // added here so func doesnt throw error, as it was previously a pkg level variable

	// check first too see if the template is already is in our cache
	_, inMap := tc[t] // look in map tc, for key t, and popluate inMap, which is true it it is there, false if not

	if !inMap {
		// template is not in our cache so we need to read from disk and parse
		log.Println("creating template and adding to cache")
		err = createTemplateCacheV1(t)
		if err != nil {
			log.Println(err)
		}

	} else {
		// else template exists in the cache
		log.Println("Using cached template")
	}

	// at this point, we can assume we have the template in cache
	err = tc[t].Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func createTemplateCacheV1(t string) error {
	var tc = make(map[string]*template.Template) // added here so func doesnt throw error, as it was previously a pkg level variable

	// this slice hold each entry required to render the web page
	// recall that the base layout is mandatory
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}

	// now parse the templates
	tmpl, err := template.ParseFiles(templates...) // "..." means to expand ther array/slice

	if err != nil {
		return err
	}
	// add template to cache
	tc[t] = tmpl

	return nil
}

var app *config.AppConfig
var pathToTemplate = "./templates"
var functions = template.FuncMap{} // variable that holds all the functions that we want to put into and make available to our goland templates

// sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// here we can specify data that maybe used for every page
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// these will automatically get populated everytime we are rendering a page (flash, error, warning)
	td.Flash = app.Session.PopString(r.Context(), "flash") // appears once and is taken out of the session
	// this will put something in the session until the next time a page is displayed and then it's taken out automatically, so it is perfect to put messages we want to send to ou user
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// Renders template utilizing html templates
// was previously named "RenderTemplate" but it is convention not to name a func with the package name
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	// 1. create template cache, and put it to the app config
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	/* We no longer manually create the cache template here we load it from configs
	tc, err := createTemplateCache() // at this point each time the function is called the cache is rebuilt

	if err != nil {
		log.Fatal(err)
	}
	*/
	// we want to avoid loading the entire template cache every call
	// We need to have some settings configurations such that once this template set is poplated we do not load it again untill the application restarts
	// the normal approach to this is to use global variables but we want a better approach

	// 2. get the requested template from cache
	t, ok := tc[tmpl]
	// t is the template we want to render
	// ok is true if the template is in the cache, false if not
	if !ok {
		log.Println("Could not get template from template cache")
		return errors.New("can't get template from cache")
	}
	// arbitrary choice here
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	// execute on the buffer and then write it out instead of executing directly
	// This is useful for fine grain error checking
	_ = t.Execute(buf, td) // thus if we cannot execute it we can figure out where the error is coming from
	// NOTE that the writer we use here is the buffer, not the regular "w"

	//if err != nil {
	//	log.Println(err) // this will tell me if the error comes from the value that is stored in the map
	//}
	// 3. render the template
	_, err := buf.WriteTo(w) // now we write the buffer out to "w"

	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func CreateTemplateCache() (map[string]*template.Template, error) {
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
