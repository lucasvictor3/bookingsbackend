package utils

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/lucasvictor3/bookingsbackend/internal/config"
	"github.com/lucasvictor3/bookingsbackend/internal/models"
)

var app *config.AppConfig
var pathToTemplates = "./templates"

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, templateData *models.TemplateData) error {
	// create a template cache
	var templateCache map[string]*template.Template
	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	// get requested template from cache
	template, ok := templateCache[tmpl]
	if !ok {
		log.Println("Could not get template from template cache")
		return errors.New("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	templateData = AddDefaultData(templateData, r)
	err := template.Execute(buf, templateData)
	if err != nil {
		log.Println(err)
		return err
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all of the files named *page.tmpl
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, nil
	}

	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = templateSet
	}

	return myCache, nil
}
