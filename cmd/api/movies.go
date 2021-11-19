package main

import (
	"fmt"
	"net/http"
	"yamda_go/internal/models"
	"yamda_go/internal/validator"

	"github.com/julienschmidt/httprouter"
)

const tag = "CreateMovieHandler:"

func (app *Application) CreateMovieHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//how we expect data from outside
	var input struct {
		Title   string         `json:"title"`
		Year    int32          `json:"year"`
		Runtime models.Runtime `json:"runtime"`
		Genres  []string       `json:"genres"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		problem := models.ErrorProblem{
			Title:  "input data not valid",
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}
		app.log.Println(tag, err.Error())
		if err = app.writeError(w, http.StatusBadRequest, problem, nil); err != nil {
			app.log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
	if v.IsValid() {
		//create new movie
		if _, err := app.movieSvc.CreateMovie(movie); err != nil {
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
		if err := app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, nil); err != nil {
			app.log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	problem := models.ErrorProblem{
		Title:  "input data not valid",
		Status: http.StatusUnprocessableEntity,
		Detail: "content of movie entity is not valid",
		Errors: v.Errors,
	}
	app.log.Println(tag, problem)
	if err := app.writeError(w, http.StatusUnprocessableEntity, problem, nil); err != nil {
		app.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (app *Application) GetMovieHandler(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	num, err := app.ParseId(p)
	if err != nil {
		problem := models.ErrorProblem{
			Title:  "error trying to parse id from route:",
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}
		app.log.Println(err)
		if err = app.writeError(w, http.StatusBadRequest, problem, nil); err != nil {
			app.log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	movie, err := app.movieSvc.GetMovie(num)
	if err != nil {
		problem := models.ErrorProblem{
			Title:  "movie not found",
			Status: http.StatusNotFound,
			Detail: fmt.Sprintf("movie with id %d not found", num),
		}
		app.log.Println(err)
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
	w.WriteHeader(http.StatusNoContent)
}

func (app *Application) DeleteMovieHandler(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {

}
