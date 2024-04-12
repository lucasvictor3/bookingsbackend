package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	form := New(r.PostForm)

	isValid := form.Valid()

	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	form := New(r.PostForm)
	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r = httptest.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("form shows invalid when required fields are present")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	postedData := url.Values{}
	postedData.Add("name", "test")
	r.PostForm = postedData

	form := New(r.PostForm)
	form.MinLength("name", 3, r)

	if !form.Valid() {
		t.Error("name property is greater than required but got error")
	}

	form.MinLength("name", 5, r)

	if form.Valid() {
		t.Error("name property is smaller than required but got valid")
	}
}

func TestForm_IsEmail(t *testing.T) {

}

func TestForm_Has(t *testing.T) {

}
