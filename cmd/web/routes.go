package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kjfs/dieffe_sensor/internal/config"
	"github.com/kjfs/dieffe_sensor/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/temperature", handlers.Repo.Temperature)
	mux.Get("/radioactivity", handlers.Repo.Radioactivity)
	mux.Get("/missing-permissions", handlers.Repo.MissingPermissions)

	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/registration", handlers.Repo.UserRegistration)
	mux.Post("/registration", handlers.Repo.UserRegistrationPost)
	mux.Get("/registration-summary", handlers.Repo.UserRegistrationSummary)

	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.ShowLoginPost)
	mux.Get("/user/logout", handlers.Repo.Logout)
	mux.Get("/user/temperature", handlers.Repo.Temperature)


	fileServer := http.FileServer(http.Dir("./static/"))             
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer)) 

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)

		mux.Get("/registrations-new", handlers.Repo.AdminRegistrationNew)
		mux.Get("/registrations-all", handlers.Repo.AdminRegistrationAll)

		mux.Get("/process-registrations/{src}/{id}/do", handlers.Repo.AdminProcessUser)
		mux.Get("/delete-registrations/{src}/{id}/do", handlers.Repo.AdminDeleteUser)
		mux.Get("/registrations/{src}/{id}/show", handlers.Repo.AdminShowRegistration)
		mux.Post("/registrations/{src}/{id}", handlers.Repo.AdminShowRegistrationPost)

	})
	return mux
}
