package utils

import (
	"net/http"
	"testing"

	"github.com/lucasvictor3/bookingsbackend/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	data := AddDefaultData(&td, r)

	if data.Flash != "123" {
		t.Error("Flash value not found or different")
	}

}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	templateCache, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = templateCache

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = RenderTemplate(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("got template that does not exist")
	}

	app.UseCache = true
	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error using cache for templates")
	}
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)

	if err != nil {
		return nil, err

	}

	ctx := r.Context()

	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))

	r = r.WithContext(ctx)

	return r, nil
}
