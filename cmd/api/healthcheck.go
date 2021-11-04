package main

import (
	"net/http"
)

// HealthCheckHandler handler writes a plain-text response with information about the
// application status, operating environment and version.
func (app *Application) HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	data := map[string]string{
		"status": "available",
		"environment": app.config.Env,
		"version": version,
	}
	if err := app.writeJSON(w, 200, data, nil); err != nil {
		app.log.Println(err)
		http.Error(w, "an error happened at the server level", http.StatusInternalServerError)
	}
}
