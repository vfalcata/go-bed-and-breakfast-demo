// This is a special filename that will run before any test run
// anytime you have a setup file, it MUST contain a function called 'TestMain'
package main

import (
	"net/http"
	"os"
	"testing"
)

// Before you start runnint the test run this
func TestMain(m *testing.M) {

	// Befor you exit run the actual tests
	os.Exit(m.Run())
}

// You can also store variables that you might need outside the test main function
// we need object to satisfy http.Handler interface
type myHandler struct{}

// this is the only function needed to satisfy the interface
func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
