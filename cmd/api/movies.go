package main

import (
	"fmt"
	"net/http"
	"yamda_go/internal/models"
	"yamda_go/internal/validator"

	"github.com/julienschmidt/httprouter"
)

const (
	tag                = "CreateMovieHandler:"
	TagUpdatehandler   = "UpdateMovieHandler:"
	ErrContentNotValid = "content of movie entity is not valid"
)

func (app *Application) CreateMovieHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//how we expect data from outside
	var input struct {
		Title   string         `json:"title"`
		Year    int32          `json:"year"`
		Runtime models.Runtime `json:"runtime"`
		Genres  []string       `json:"genres"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, err)
		return
	}

	//validate movie contents
	v := validator.New()
	// Copy the values from the input struct to a new Movie struct.
	movie := &models.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}
	movie.Validate(v)
	if !v.IsValid() {
		app.failedValidationResponse(w, v.Errors)
		return 
	}
	//create new movie
	if _, err := app.provider.Insert(movie); err != nil {
		problem := models.ErrorProblem{
			Title:  "movie not created",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
		app.log.Println(tag, err.Error())
		if err = app.writeError(w, http.StatusInternalServerError, problem, nil); err != nil {
			app.log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	headers := make(http.Header)
  headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))
	if err := app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers); err != nil {
		app.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Application) GetMovieHandler(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	num, err := app.ParseId(p)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	movie, err := app.provider.Get(num)
	if err != nil {
		problem := models.ErrorProblem{
			Title:  "movie not found",
			Status: http.StatusNotFound,
			Detail: fmt.Sprintf("movie with id %d not found", num),
		}
		if err = app.writeError(w, http.StatusNotFound, problem, nil); err != nil {
			app.log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil); err != nil {
		app.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Application) UpdateMovieHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var m struct {
		ID      int64          `json:"id"`
		Title   string         `json:"title"`
		Year    int32          `json:"year"`
		Runtime models.Runtime `json:"runtime"`
		Genres  []string       `json:"genres"`
	}
	if err := app.readJSON(w, r, &m); err != nil {
		app.badRequestResponse(w, err)
		return
	}

	//validate movie contents
	v := validator.New()
	// Copy the values from the input struct to a new Movie struct.
	movie := &models.Movie{
		ID:      m.ID,
		Title:   m.Title,
		Year:    m.Year,
		Runtime: m.Runtime,
		Genres:  m.Genres,
	}
	movie.ValidateWithId(v)
	if !v.IsValid() {
		app.failedValidationResponse(w, v.Errors)
		return 
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *Application) DeleteMovieHandler(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	num, err := app.ParseId(p)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	if err := app.provider.Delete(num); err != nil {
		notFndErr := models.ErrorProblem{
			Title:  "movie not found",
			Status: http.StatusNotFound,
			Detail: fmt.Sprintf("movie with id %d not found", num),
		}
		if err = app.writeError(w, http.StatusNotFound, notFndErr, nil); err != nil {
			app.log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func parsingError(err error) models.ErrorProblem {
	return models.ErrorProblem{
		Title:  "error trying to parse id from route:",
		Status: http.StatusBadRequest,
		Detail: err.Error(),
	}
}
