package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

//TODO generate this automatically at build time
const version = "0.0.1"

//Application type contains all dependencies for the top layer of
//the API.
type Application struct {
	log *log.Logger
	config *config.Settings
}


// HealthCheckHandler handler writes a plain-text response with information about the
// application status, operating environment and version.
func (app *Application) HealthCheckHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.Env)
	fmt.Fprintf(w, "version: %s\n", version)
}

func (app *Application) CreateMovieHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

}

func (app *Application) GetMovieHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

}