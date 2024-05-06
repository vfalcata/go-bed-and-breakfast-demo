package models

import "myapp/internal/forms"

// Holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} // In go when you are unsure what the type will be use an interface, its like a worse version of a generic
	CSRFToken string                 // stands for Cross Site Request Forgery Token. When you build a webpage with a form on it you have a hidden field in that form, which is a long string of random numbers and they change every single time somebody goes to a page
	Flash     string                 // appears once and is taken out of the session
	Warning   string
	Error     string
	Form      *forms.Form
}
