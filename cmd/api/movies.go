package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
	"yamda_go/internal/models"
)

func (app *Application) CreateMovieHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

}

func (app *Application) GetMovieHandler(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	num, err := app.ParseId(p)
	if err != nil {
		problem := envelope{
			"error" : models.ErrorProblem{
				Title:  "error trying to parse id from route:",
				Status: http.StatusBadRequest,
				Detail: err.Error(),
			},
		}
		app.log.Println(problem)
		app.writeJSON(w, http.StatusBadRequest, problem, nil)
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
		http.Error(w, "an error happened at the server level", http.StatusInternalServerError)
	}
}