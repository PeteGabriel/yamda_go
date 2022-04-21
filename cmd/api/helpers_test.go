package main

import (
	"testing"
	"yamda_go/internal/validator"

	is2 "github.com/matryer/is"
)

func TestReadString_OK(t *testing.T) {

	is := is2.New(t)

	qs := map[string][]string{"title": {"Last Man Down"}}

	title := app.readString(qs, "title", "")

	is.Equal(title, "Last Man Down")
}

func TestReadString_WithoutKeyInQueryString_GetDefaultValue(t *testing.T) {
	is := is2.New(t)

	qs := map[string][]string{"not_title": {"Last Man Down"}}

	title := app.readString(qs, "title", "default title")

	is.Equal(title, "default title")
}

func TestReadCSV_OK(t *testing.T) {

	is := is2.New(t)

	qs := map[string][]string{"genres": {"Crime,Drama,Thriller"}}

	genres := app.readCSV(qs, "genres", []string{""})

	is.Equal(genres, []string{"Crime", "Drama", "Thriller"})
}

func TestReadCSV_WithoutKeyInQueryString_GetDefaultValue(t *testing.T) {

	is := is2.New(t)

	qs := map[string][]string{"not_genres": {"Crime,Drama,Thriller"}}

	genres := app.readCSV(qs, "genres", []string{""})

	is.Equal(genres, []string{""})
}

func TestReadInt_OK(t *testing.T) {

	is := is2.New(t)

	qs := map[string][]string{"runtime": {"123"}}

	v := validator.New()
	runtime := app.readInt(qs, "runtime", -1, v)

	is.True(v.IsValid())
	is.Equal(runtime, 123)
}

func TestReadInt_KeyIsNotValidInteger(t *testing.T) {

	is := is2.New(t)

	qs := map[string][]string{"runtime": {"NaN"}}

	v := validator.New()
	runtime := app.readInt(qs, "runtime", -1, v)

	is.True(!v.IsValid()) //has error
	is.Equal(runtime, -1)
}
