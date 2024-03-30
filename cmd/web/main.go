package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/lucasvictor3/bookingsbackend/pkg/config"
	"github.com/lucasvictor3/bookingsbackend/pkg/handlers"
	"github.com/lucasvictor3/bookingsbackend/pkg/utils"
)

const port = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	// change this to true when in prod
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	templateCache, err := utils.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache:", err)
	}

	app.TemplateCache = templateCache
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	utils.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting applicatio at port: %s", port))

	serve := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()
	log.Fatal(err)
}
