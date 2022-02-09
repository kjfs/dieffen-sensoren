package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/kjfs/dieffe_sensor/helpers"
	"github.com/kjfs/dieffe_sensor/internal/config"
	"github.com/kjfs/dieffe_sensor/internal/driver"
	"github.com/kjfs/dieffe_sensor/internal/handlers"
	"github.com/kjfs/dieffe_sensor/internal/models"
	"github.com/kjfs/dieffe_sensor/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(app.MailChan)

	log.Println("Starting Mail listener")

	listenForMail()

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {

	gob.Register(models.User{})
	gob.Register(models.SensorSearch{})
	gob.Register(models.Sensor{})
	gob.Register(models.MeasurementKPI{})
	gob.Register(models.MeasurementValue{})
	gob.Register(map[string]int{})

	inProduction := flag.Bool("production", true, "Note: Application is in production")
	useCache := flag.Bool("cache", true, "Note: use template cache")
	dbName := flag.String("dbname", "", "Note: db name")
	dbHost := flag.String("dbhost", "localhost", "Note: db host")
	dbUser := flag.String("dbuser", "", "Note: db user")
	dbPass := flag.String("dbpassword", "", "Note: db password")
	dbPort := flag.String("dbport", "5432", "Note: db port")
	dbSSL := flag.String("dbSSL", "disable", "Note: db SSL settings (disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		log.Println("Missing required flags!")
		os.Exit(1)
	}

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	app.InProduction = *inProduction

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	log.Println("Connecting to db ...")

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)

	if err != nil {
		log.Fatal("Can't connect to db.", err)
	}

	log.Print("Connected to db \n")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache", err)
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = *useCache

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
