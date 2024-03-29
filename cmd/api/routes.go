package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.HealthCheckHandler)
	router.Handle(http.MethodPost, "/v1/movies", app.CreateMovieHandler)
	router.Handle(http.MethodGet, "/v1/movies/:id", app.GetMovieHandler)
	router.Handle(http.MethodPatch, "/v1/movies", app.UpdateMovieHandler)
	router.Handle(http.MethodDelete, "/v1/movies/:id", app.DeleteMovieHandler)

	router.Handle(http.MethodPatch, "/v1/movies/:id", app.PartialUpdateMovieHandler)

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.ListMoviesHandler)

	router.Handle(http.MethodPost, "/v1/users", app.RegisterUserHandler)

	//ensure middleware is always called last
	return app.recoverPanic(app.rateLimit(router))
}
