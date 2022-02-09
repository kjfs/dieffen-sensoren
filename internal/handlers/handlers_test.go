package handlers

import (
	"net/http"
)

var theTests = []struct {
	name   string // Name for the the Test.
	url    string // path which is matched by routes.
	method string
	// params             []postData // data which will be posted.
	expectedSTatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"temp", "/temperature", "GET", http.StatusOK},
	{"radio", "/radioactivity", "GET", http.StatusOK},
	{"search", "/search", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"no-existent-page", "/not_existent", "GET", http.StatusNotFound}, // This test if a route does not exist!
	// New Routes:
	{"Login", "/user/login", "GET", http.StatusOK},
	{"LogOut", "/user/logout", "GET", http.StatusOK},
	{"Dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"Calendar", "/admin/registrations-calendar", "GET", http.StatusOK},
	{"New Registration", "/admin/registrations-new", "GET", http.StatusOK},
	{"All Registration", "/admin/registrations-all", "GET", http.StatusOK},
	{"Show Registration", "/admin/registrations/new/1/show", "GET", http.StatusOK},
}
