package render

import (
	"encoding/gob"
	"log"
	"myapp/internal/config"
	"myapp/internal/models"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
)

var session *scs.SessionManager // need session data to make render request
var testApp config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{}) // We MUST register the type for the session data or it will not know.
	testApp.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	/* session := scs.New() */ //variable Shadowing, this session is not the same as our global outside scope one
	session = scs.New()        // this is how we assign the outside scope variable
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true                  // should the cookie persist after the user closes the window or their web browser
	session.Cookie.SameSite = http.SameSiteLaxMode // how strict do you want to be about what site this cookie applies to lax mode is the default
	session.Cookie.Secure = false                  // this will insis that the cookies be encrypted and that the connection is from https and not http, we set it to false right now for dev mode, but in prod we want it true

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct {
}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}
func (tw *myWriter) WriteHeader(i int) {

}
func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
