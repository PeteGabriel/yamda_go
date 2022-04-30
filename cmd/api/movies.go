package main

import (
	"errors"
	"fmt"
	"net/http"
	"yamda_go/cmd/api/dto"
	"yamda_go/internal/data"
	"yamda_go/internal/data/provider"
	"yamda_go/internal/models"
	"yamda_go/internal/validator"

	"github.com/julienschmidt/httprouter"
)

/**
* POST /v1/movies -> 201 CREATED with JSON content
*
* LOCATION header contains the uri to the newly created resource.
**/
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
		app.serverErrorResponse(w, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))
	if err := app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers); err != nil {
		app.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/**
* GET /v1/movies/:id -> 200 OK with JSON content
**/
func (app *Application) GetMovieHandler(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	num, err := app.ParseId(p)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	movie, err := app.provider.Get(num)
	if err != nil {
		switch {
		case errors.Is(err, provider.ErrRecordNotFound):
			app.resourceNotFoundResponse(w, fmt.Errorf("movie with id %d not found", num))
			return
		default:
			app.serverErrorResponse(w, err)
			return
		}
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil); err != nil {
		app.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/**
* PATCH /v1/movies -> 200 OK with JSON content
**/
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

/**
* PATCH /v1/movies/:id -> 200 OK with JSON content
**/
func (app *Application) PartialUpdateMovieHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := app.ParseId(p)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	movie, err := app.provider.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, provider.ErrRecordNotFound):
			app.resourceNotFoundResponse(w, fmt.Errorf("movie with id %d not found", id))
			return
		default:
			app.serverErrorResponse(w, err)
			return
		}
	}

	i := dto.Input{}
	if err := app.readJSON(w, r, &i); err != nil {
		app.badRequestResponse(w, err)
		return
	}

	if i.Title != nil {
		movie.Title = *i.Title
	}
	if i.Year != nil {
		movie.Year = *i.Year
	}
	if i.Runtime != nil {
		movie.Runtime = *i.Runtime
	}
	if i.Genres != nil {
		movie.Genres = i.Genres
	}

	v := validator.New()
	movie.ValidateWithId(v)
	if !v.IsValid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	err = app.provider.Update(*movie)
	if err != nil {
		switch {
		case errors.Is(err, provider.ErrEditConflict):
			app.resourceEditConflictResponse(w)
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

/**
* DELETE /v1/movies/:id -> 200 OK with JSON content
**/
func (app *Application) DeleteMovieHandler(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	num, err := app.ParseId(p)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	if err := app.provider.Delete(num); err != nil {
		switch {
		case errors.Is(err, provider.ErrRecordNotFound):
			app.resourceNotFoundResponse(w, fmt.Errorf("movie with id %d not found", num))
			return
		default:
			app.log.Println(err)
			app.serverErrorResponse(w, err)
			return
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

/**
* GET /v1/movies -> 200 OK with JSON content
**/
func (app *Application) ListMoviesHandler(w http.ResponseWriter, r *http.Request) {

	v := validator.New()
	qs := r.URL.Query()

	input := data.Search{
		Title:  app.readString(qs, "title", ""),
		Genres: app.readString(qs, "genres", ""),
		Filters: data.Filter{
			Page:         app.readInt(qs, "page", 1, v),       //default page is 1
			PageSize:     app.readInt(qs, "page_size", 20, v), //default size is 20
			Sort:         app.readString(qs, "sort", "id"),    //default sort if by 'id' field
			SortSafelist: []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"},
		},
	}

	if input.Filters.Validate(v); !v.IsValid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	movies, meta, err := app.provider.GetAll(input)
	if err != nil {
		app.serverErrorResponse(w, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": meta, "movies": movies}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
