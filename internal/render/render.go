package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
	"github.com/kjfs/dieffe_sensor/internal/config"
	"github.com/kjfs/dieffe_sensor/internal/models"
)

var functions = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
} // a map of functions, which will be used later: to pass functions to templates. e.g. format a date. return the current year.

var app *config.AppConfig
var pathToTemplates = "./templates"

//Add adds 1 to val
func Add(a, b int) int {
	return a + b
}

// Iterate returns a []ints. Starting at 1, going to count
func Iterate(count int) []int {
	var i int
	var items []int
	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

//NewRenderer sets the config for the template package.
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, req *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(req.Context(), "flash")
	td.Error = app.Session.PopString(req.Context(), "error")
	td.Warning = app.Session.PopString(req.Context(), "warning")
	td.CSRFToken = nosurf.Token(req)
	if app.Session.Exists(req.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	if app.Session.GetInt(req.Context(), "access_lvl") == 2 {
		td.IsAdmin = 2
	}
	return td
}

// HumanDate returns time in YYYY-MM-DD
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatData returns time in YYYY-MM-DD
func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

// Template Renders templates using html/template
func Template(w http.ResponseWriter, tmpl string, td *models.TemplateData, r *http.Request) error {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache // Get the templatecache from the app config.
	} else {
		tc, _ = CreateTemplateCache()
	}
	t, ok := tc[tmpl] // Pull the template out of the map. Wir muessen das template aus dem templatcache heraus abstrahieren, das wird wollen bzw der User anschaut.
	if !ok {
		return errors.New("can't get template from cache")
	} // We need to check if the template does / -> does not exist in our map. Des weiteren muessen wir schauen, ob das Template uberhaupt existiert. Daher: ,ok
	/*
		We have not read from disk. "t, ok := tc[tmpl]" does only exist in memory. We need to stored this in a bytes buffer.
	*/
	buf := new(bytes.Buffer)         // buf will hold bytes (for our templates).
	td = AddDefaultData(td, r)       // Adds default data from above.
	_ = t.Execute(buf, td)           // Put template, which is hold in t, into the buf variable and add extra data into it (nil, so no data this time).
	_, errWriteToW := buf.WriteTo(w) // Write the data to the response write.
	if errWriteToW != nil {
		fmt.Println("Error writing template to browser", errWriteToW)
		return errWriteToW
	}
	return nil
}

// CreateTemplateCache creates a TemplateCache as a map.
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}                                  // myCache stores the templates. Its a go function.
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates)) // Need to find and save all templates "*pages.tmpl", which are available.
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)                                     // We need to get the base name of the page. Stored is the full path of the *page.tmpl.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page) // 1. Lets create a template set with the data form above. 2. Add possiblity to add functions on top of those templates + 3. Parse the files.
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*layout.tmpl", pathToTemplates)) // Does this template match any layouts?
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		} // If there is one thing, then length will be greater than zero.
		myCache[name] = ts // Adds the template set to the Cache.
	} // Range through the pages! indices and value will be available.
	return myCache, nil // because we have a fucntino here. "return nil" is needed!
}
