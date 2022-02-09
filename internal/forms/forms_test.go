package forms

import (
	"net/url"
	"testing"
)

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("Form shows valid email for non-existing field")
	}

	postedData = url.Values{}
	postedData.Add("Email", "k.studnik@gmx.de")
	form = New(postedData)

	form.IsEmail("Email")
	if !form.Valid() {
		t.Error("Got an invalid Email when we should not have")
	}

	postedData = url.Values{}
	postedData.Add("Email", "k.studnik_gmx.de")
	form = New(postedData)

	form.IsEmail("Email")
	if form.Valid() {
		t.Error("Got a valid, but an invalid Email address")
	}

}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.MinLength("a", 3)
	if form.Valid() {
		t.Error("form shows min length for non - existing field")
	}

	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	postedValues := url.Values{}
	postedValues.Add("WS", "0123")
	form = New(postedValues)

	form.MinLength("WS", 1)
	if !form.Valid() {
		t.Error("shows minlength of 1 is not met, when it is")
	}

	isError = form.Errors.Get("WS")
	if isError != "" {
		t.Error("Should not have an error, but got one")
	}

	postedValues = url.Values{}
	postedValues.Add("test", "abc123")
	form = New(postedValues)

	form.MinLength("test", 100)
	if form.Valid() {
		t.Error("form shows wrong min length: Data is shorter")
	}

}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	Has := form.Has("v")

	if Has {
		t.Error("form is nil, but it shows data")
	}

	postedData = url.Values{}
	postedData.Add("a", "1")

	form = New(postedData)

	Has = form.Has("a")

	if !Has {
		t.Error("form shows invalid altough there is data")
	}

}

func TestForm_Valid(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when we should valid")
	}
}

func TestForm_Required(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when requried fields missing")
	}

	postedData = url.Values{}
	postedData.Add("a", "1")
	postedData.Add("b", "2")
	postedData.Add("c", "3")

	form = New(postedData)

	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}
