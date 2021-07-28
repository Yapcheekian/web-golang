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
	isValid := form.Valid()

	if isValid {
		t.Error("form is valid when required fields is missing")
	}

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "b")
	postData.Add("c", "c")

	r = httptest.NewRequest("POST", "/whatever", nil)
	r.PostForm = postData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	isValid = form.Valid()

	if !isValid {
		t.Error("show does not have required when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postData := url.Values{}
	form := New(postData)
	has := form.Has("a")
	if has {
		t.Error("form show has field when it does not")
	}

	postData = url.Values{}
	postData.Add("a", "a")
	form = New(postData)

	has = form.Has("a")

	if !has {
		t.Error("form show does not has field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	postData := url.Values{}
	form := New(postData)

	form.MinLength("a", 3)

	if form.Valid() {
		t.Error("form pass min length check for non-existent field")
	}

	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	postData = url.Values{}
	postData.Add("a", "aaaaa")
	form = New(postData)

	form.MinLength("a", 3)

	if !form.Valid() {
		t.Error("form does not pass min length check when it should")
	}

	isError = form.Errors.Get("a")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}

	postData = url.Values{}
	postData.Add("a", "a")
	form = New(postData)

	form.MinLength("a", 3)

	if form.Valid() {
		t.Error("form pass min length check when it should not")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postData := url.Values{}
	form := New(postData)

	form.IsEmail("a")

	if form.Valid() {
		t.Error("form show valid email for non-existent field")
	}

	postData = url.Values{}
	postData.Add("a", "a@test.com")
	form = New(postData)

	form.IsEmail("a")

	if !form.Valid() {
		t.Error("got an invalid email when we should not have")
	}

	postData = url.Values{}
	postData.Add("a", "a@test")
	form = New(postData)

	form.IsEmail("a")

	if form.Valid() {
		t.Error("got valid for invalid email")
	}
}
