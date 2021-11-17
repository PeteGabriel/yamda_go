package models

import (
	"strings"
	"time"
	"yamda_go/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Runtime   Runtime   `json:"runtime,omitempty,string"`
	Genres    []string  `json:"genres,omitempty"`
	Year      int32     `json:"year,omitempty"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"-"`
}

//Validate uses a validator interface to validate the contents of a given movie.
func (m *Movie) Validate(v *validator.Validator) {
	v.Check(m.Title != "", "title", "must be provided")
	title := strings.TrimSpace(m.Title)
	v.Check(len(title) >= 1 && len(title) <= 500, "title", "must not be empty or more than 500 bytes long")
	v.Check(m.Year != 0, "year", "must be provided")
	v.Check(m.Year >= 1888, "year", "must be greater than 1888")
	v.Check(m.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(m.Runtime != 0, "runtime", "must be provided")
	v.Check(m.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(m.Genres != nil, "genres", "must be provided")
	v.Check(len(m.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(m.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(m.Genres), "genres", "must not contain duplicate values")
}
