package main

import (
	"fmt"
	"log"
	"net/http"
	"yamda_go/internal/config"
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
func (app *Application) HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.Env)
	fmt.Fprintf(w, "version: %s\n", version)
}