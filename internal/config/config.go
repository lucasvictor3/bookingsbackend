package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/lucasvictor3/bookingsbackend/internal/models"
)

// AppConfig holds the application config
type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
