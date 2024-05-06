package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"make-res", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"post-search-avail", "/search-availability", "Post", []postData{
		{key: "start", value: "2024-05-04"},
		{key: "end", value: "2024-05-05"},
	}, http.StatusOK},
	{"post-search-avail-json", "/search-availability-json", "Post", []postData{
		{key: "start", value: "2024-05-04"},
		{key: "end", value: "2024-05-05"},
	}, http.StatusOK},
	{"make-reservation-post", "/make-reservation", "Post", []postData{
		{key: "first_name", value: "Lok"},
		{key: "last_name", value: "Cas"},
		{key: "email", value: "lok@my.email"},
		{key: "phone", value: "333-333-3333"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close() // defer doesn't get executed until after the current function is finished. So when our test handlers are finished it will close the test server.

	for _, e := range theTests {
		if e.method == "GET" {
			// ts allows us to make client calls
			resp, err := ts.Client().Get(ts.URL + e.url) // we do not know the port or anything for the test server, but "ts.URL" does, so we use this. The "e.url" gives us the path without the port.
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			// create  empty variable that is required on our forms for the method on our server
			values := url.Values{} // holds information as a post request
			for _, x := range e.params {
				values.Add(x.key, x.value) // now values has everything it needs to make the post request
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
