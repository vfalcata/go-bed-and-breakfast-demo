package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// its pretty common to name your parameter next for your middleware
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("You Hit the PAGE!!!")

		// After the above we must move on to the "next" thing which maybe another middleware or another page
		next.ServeHTTP(w, r)
	})
}

// CSRF protection middleware

// nosurf requires that you have a hidden field, or at least a field that doesn't have to be hidden on the actual form that's doing the post and it needs to have a certain name
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
