package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Yapcheekian/web-golang/internal/config"
	"github.com/Yapcheekian/web-golang/internal/driver"
	"github.com/Yapcheekian/web-golang/internal/handlers"
	"github.com/Yapcheekian/web-golang/internal/helpers"
	"github.com/Yapcheekian/web-golang/internal/models"
	"github.com/Yapcheekian/web-golang/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var (
	app      config.AppConfig
	session  *scs.SessionManager
	infoLog  *log.Logger
	errorLog *log.Logger
)

func main() {
	db, err := run()

	defer db.SQL.Close()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	gob.Register(models.Reservation{})

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app.InfoLog = infoLog
	app.ErrorLog = errorLog

	app.IsProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction
	app.Session = session

	log.Println("connecting to database...")
	db, err := driver.ConnectSQL("host=127.0.0.1 port=5432 dbname=bookings user=postgres password=mysecretpassword")
	if err != nil {
		log.Fatal("cannot connect to database! Dying...")
	}
	log.Println("connected to the database")

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = app.IsProduction

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewTemplate(&app)
	helpers.New(&app)

	return db, nil
}
