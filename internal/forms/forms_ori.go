package forms

//
// import (
// 	"fmt"
// 	"net/http"
// 	"net/url"
// 	"strings"

// 	"github.com/asaskevich/govalidator"
// )

// // custom form, that has both url and errors
// type Form struct {
// 	url.Values
// 	Errors errors
// }

// // if there are any errors at all that we are getting from the form then retrun a bool (true if no errors, false if there is)
// func (f *Form) Valid() bool {
// 	return len(f.Errors) == 0
// }

// // initialize a form struct, a constructor
// func New(data url.Values) *Form {
// 	return &Form{
// 		data,                          // data from arg
// 		errors(map[string][]string{}), // empty map
// 	}
// }

// // checks that required fields have values
// func (f *Form) Required(fields ...string) {
// 	for _, field := range fields {
// 		value := f.Get(field) // get the value of all form fields
// 		if strings.TrimSpace(value) == "" {
// 			f.Errors.Add(field, "This field cannot be blank")
// 		}
// 	}
// }

// // ensures that the field is of a min length
// func (f *Form) MinLength(field string, length int, r *http.Request) bool {
// 	x := r.Form.Get(field)
// 	if len(x) < length {
// 		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
// 		return false
// 	}
// 	return true
// }

// // does the form, from the post request have some field ?
// // we no longer need this as the above covers it all
// func (f *Form) Has(field string, r *http.Request) bool {
// 	x := r.Form.Get(field)
// 	if x == "" {
// 		// f.Errors.Add(field, "This field cannot be blank")
// 		// We no longer want to push an error for this
// 		// we still keep it since it will become a useful function because there are certain fields, like checkbox that are actually handled differently than text input
// 		// For example if there is a mandatory checkbox, and we want it part of the post request, if it is not checked, it is not included at all, so we cannot really do anything with it
// 		// this is useful to check whether or not submitted form data includes certain fields
// 		return false
// 	}
// 	return true
// }

// // ensures email is valid
// func (f *Form) IsEmail(field string) {
// 	if !govalidator.IsEmail(f.Get(field)) {
// 		f.Errors.Add(field, "Invalid email address")
// 	}
// }
