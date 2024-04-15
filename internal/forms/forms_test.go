package forms

import (
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	postedData := url.Values{}

	form := New(postedData)

	isValid := form.Valid()

	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	form = New(postedData)
	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("form shows invalid when required fields are present")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("name", "test")

	form := New(postedData)

	form.MinLength("x", 3)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have an error")
	}

	form = New(postedData)

	form.MinLength("name", 3)

	if !form.Valid() {
		t.Error("name property is greater than required but got error")
	}

	form.MinLength("name", 5)

	if form.Valid() {
		t.Error("name property is smaller than required but got valid")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("email", "test@test.com")
	form := New(postedData)

	form.IsEmail("email")

	if !form.Valid() {
		t.Error("email property value is valid but got error")
	}

	postedData = url.Values{}
	postedData.Add("email", "test@")
	form = New(postedData)

	form.IsEmail("email")

	if form.Valid() {
		t.Error("email property value is invalid but got valid")
	}

}

func TestForm_Has(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("email", "test@test.com")
	postedData.Add("name", "test")
	form := New(postedData)

	has := form.Has("name")

	if !has {
		t.Error("name property exists but hasnt found")
	}

	has = form.Has("error")

	if has {
		t.Error("error property not exists but has found")
	}

}
