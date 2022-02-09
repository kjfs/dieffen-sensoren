package render

import (
	"net/http"
	"testing"

	"github.com/kjfs/dieffe_sensor/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	//req, err :=	http.NewRequestWithContext(ctx context.Context, method string, url string, body io.Reader)
	// req, err := http.NewRequest("POST", "/some-URL", nil) // URL irrelevant, daher fake. Body einfach nil, da ebenfalls unnoetig.
	// if err != nil {
	// 	t.Error(err)
	// }

	req, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(req.Context(), "flash", "123")

	result := AddDefaultData(&td, req) // AddDefaultData returns '*models.TemplateData' which stores all HTMl5 messages like warning, flash, error.
	if result.Flash != "123" {
		t.Error("Flash value '123' not found in session")
	} // returns *models.TemplateData in result VAR

}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"

	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	} // This will create the templateCache in this function.
	app.TemplateCache = tc // Will add the TC into app VAR. This is a package level Variable = This is pakcage level = render and the app is in the same level (its in render.go)

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = Template(&ww, "home.page.tmpl", &models.TemplateData{}, r)
	if err != nil {
		t.Error("ERROR: Error writing template to browser.")
	}

	err = Template(&ww, "non-existent.page.tmpl", &models.TemplateData{}, r)
	if err == nil {
		t.Error("ERROR: Rendered template that does not exist.")
	}

}

func TestNewRenderer(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func getSession() (*http.Request, error) {

	req, err := http.NewRequest("POST", "/some-URL", nil) // URL irrelevant, daher fake. Body einfach nil, da ebenfalls unnoetig.
	if err != nil {
		return nil, err
	} // Neuer request wird erstellt.

	ctx := req.Context()                                      // Ein Contxt wird erstelly, der 'context.Background()' returned.
	ctx, err = session.Load(ctx, req.Header.Get("X-Session")) // Load retrieves the session data for the given token from the session store, and returns a new context.Context containing the session data. If no matching token is found then this will create a new session.
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx) // WithContext returns a shallow copy of r with its context changed to ctx. The provided ctx must be non-nil.

	return req, nil

}
