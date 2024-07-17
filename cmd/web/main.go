package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/lucasvictor3/bookingsbackend/driver"
	"github.com/lucasvictor3/bookingsbackend/internal/config"
	"github.com/lucasvictor3/bookingsbackend/internal/handlers"
	"github.com/lucasvictor3/bookingsbackend/internal/helpers"
	"github.com/lucasvictor3/bookingsbackend/internal/models"
	"github.com/lucasvictor3/bookingsbackend/internal/utils"
)

var port string

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port

	err := godotenv.Load("/etc/secrets/.env")
	if err != nil {
		err = godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	defer close(app.MailChan)
	listerForMail()

	fmt.Println(fmt.Sprintf("Starting applicatio at port: %s", port))

	serve := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {

	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})
	gob.Register(map[string]int{})

	// read flags
	inProduction, err := strconv.ParseBool(os.Getenv("PRODUCTION"))
	if err != nil {
		fmt.Println("production flag error!")
		os.Exit(1)
	}
	useCache, err := strconv.ParseBool(os.Getenv("CACHE"))
	if err != nil {
		fmt.Println("production flag error!")
		os.Exit(1)
	}
	dbName := os.Getenv("DBNAME")
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbPort := os.Getenv("DBPORT")
	dbHost := os.Getenv("DBHOST")
	dbSSL := os.Getenv("DBSSL")

	flag.Parse()

	if dbName == "" || dbUser == "" {
		fmt.Println("Missing required flags!")
		os.Exit(1)
	}

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// NOTE: change this to true when in prod
	app.InProduction = inProduction
	app.UseCache = useCache

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

	log.Println("Connecting to the database...")
	connectionString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", dbHost, dbPort, dbName, dbUser, dbPass, dbSSL,
	)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}

	log.Println("Connected to the database!")

	templateCache, err := utils.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache:", err)
		return nil, err
	}

	app.TemplateCache = templateCache

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)

	utils.NewRenderer(&app)

	return db, nil
}
