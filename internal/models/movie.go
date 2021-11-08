package models

import "time"

type Movie struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Runtime   Runtime   `json:"runtime,omitempty,string"`
	Genres    []string  `json:"genres,omitempty"`
	Year      int32     `json:"year,omitempty"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"-"`
}
