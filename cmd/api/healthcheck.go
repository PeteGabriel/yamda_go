package main

import (
	"fmt"
	"net/http"
)

// HealthCheckHandler handler writes a plain-text response with information about the
// application status, operating environment and version.
func (app *Application) HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.Env)
	fmt.Fprintf(w, "version: %s\n", version)
}
