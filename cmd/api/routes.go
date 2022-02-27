package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *Application) routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.HealthCheckHandler)
	router.Handle(http.MethodPost, "/v1/movies", app.CreateMovieHandler)
	router.Handle(http.MethodGet, "/v1/movies/:id", app.GetMovieHandler)
	router.Handle(http.MethodPatch, "/v1/movies", app.UpdateMovieHandler)
	router.Handle(http.MethodDelete, "/v1/movies/:id", app.DeleteMovieHandler)

	router.Handle(http.MethodPut, "/v1/movies/:id", app.PartialUpdateMovieHandler)

	return router
}
