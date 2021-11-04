package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *Application) CreateMovieHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

}

func (app *Application) GetMovieHandler(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	num, err := app.ParseId(p)
	if err != nil {
		//TODO handle error properly
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	app.log.Println("get movie with id", num)
	fmt.Fprint(w, "details of movie with id ", num)
}