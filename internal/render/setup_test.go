package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/kjfs/dieffe_sensor/internal/config"
	"github.com/kjfs/dieffe_sensor/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	// What am i going to put in the session. We need tell our app what we want to add in the session.
	gob.Register(models.User{})

	// Change this to true when you are in production!
	testApp.InProduction = false

	session = scs.New()                            // It does not matter were to create the session Variable.
	session.Lifetime = 24 * time.Hour              // Need to declare how long the sessino should exist.
	session.Cookie.Persist = true                  // Need to create a store for the Coockie. User quits browser window, should the data persist?
	session.Cookie.SameSite = http.SameSiteLaxMode // Tells how strict the cookie shpuld apply to one site.
	session.Cookie.Secure = testApp.InProduction   // Pease use in Production use true! We need a encryption for this.
	testApp.Session = session                      // Session is declared, now: push the session into app struct app.Session. Gives access to the session data.
	app = &testApp
	os.Exit(m.Run())
}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {

	var h http.Header
	return h

}
func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
func (tw *myWriter) WriteHeader(statusCode int) {}
