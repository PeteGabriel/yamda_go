package main

import (
	"net/http"
	"time"
	"yamda_go/internal/models"
	"yamda_go/internal/validator"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) CreateMovieHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//how we expect data from outside
	var input struct {
		Title string `json:"title"`
		Year int32 `json:"year"`
		Runtime models.Runtime `json:"runtime"`
		Genres []string `json:"genres"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		problem :=models.ErrorProblem{
			Title:  "input data not valid",
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}
		app.log.Println("CreateMovieHandler:", err.Error())
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
	models.ValidateMovie(v, movie)
	if v.IsValid() {
		if err := app.writeJSON(w, http.StatusCreated, envelope{}, nil); err != nil {
			app.log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}else {
		problem :=models.ErrorProblem{
			Title:  "input data not valid",
			Status: http.StatusUnprocessableEntity,
			Detail: "content of movie entity is not valid",
			Errors: v.Errors,
		}
		app.log.Println("CreateMovieHandler:", problem)
		if err := app.writeError(w, http.StatusUnprocessableEntity, problem, nil); err != nil {
			app.log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
	app.log.Println("get movie with id", num)
	movie := models.Movie{
		ID:        1,
		Title:     "Casablanca",
		Runtime:   0,
		Genres:    []string{"drama", "war", "romance"},
		Year:      0,
		Version:   1,
		CreatedAt: time.Time{},
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil); err != nil {
		app.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}