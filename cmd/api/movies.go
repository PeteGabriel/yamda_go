package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
	"yamda_go/internal/models"
)

func (app *Application) CreateMovieHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//how we expect data from outside
	var input struct {
		Title string `json:"title"`
		Year int32 `json:"year"`
		Runtime int32 `json:"runtime"`
		Genres []string `json:"genres"`
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() //do not accept unknown fields in input
	if err := dec.Decode(&input); err != nil {
		problem :=models.ErrorProblem{
			Title:  "input data not valid",
			Status: http.StatusBadRequest,
			Detail: "input data could not be decoded into expected structure",
		}
		app.log.Println("CreateMovieHandler:", err.Error())
		if err = app.writeError(w, http.StatusBadRequest, problem, nil); err != nil {
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