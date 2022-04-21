package main

import (
	"net/url"
	"strconv"
	"strings"
	"yamda_go/internal/validator"
)

// readString extracts a string value from a key in the query string of a request.
// e.g: /v1/movies?title=godfather -> "godfather" for the key "title"
// readString(qs, "title", "")
func (app *Application) readString(qs url.Values, key string, defaultVal string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultVal
	}
	return s
}

// readCSV splits a string value on the comma character from a key in the query string of a request .
// e.g: /v1/movies?genres=Crime,Drama,Thriller -> ["Crime", "Drama", "Thriller"] for the key "genres"
// readCSV(qs, "genres")
func (app *Application) readCSV(qs url.Values, key string, defaultVal []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return defaultVal
	}
	return strings.Split(csv, ",")
}

func (app *Application) readInt(qs url.Values, key string, defaultVal int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultVal
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultVal
	}

	return i
}
