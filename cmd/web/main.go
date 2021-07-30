package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Yapcheekian/web-golang/internal/config"
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
	err := run()

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

func run() error {
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

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal(err)
		return err
	}

	app.TemplateCache = tc
	app.UseCache = app.IsProduction

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplate(&app)
	helpers.New(&app)

	return nil
}
