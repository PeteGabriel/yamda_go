package dto

import "yamda_go/internal/models"

//Input should be used to deserialize content coming from requests.
type Input struct {
	Title   *string         `json:"title"`
	Year    *int32          `json:"year"`
	Runtime *models.Runtime `json:"runtime"`
	Genres  []string        `json:"genres"`
}
