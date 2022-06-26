package main

import (
	"net/http"
)

// HealthCheckHandler handler writes a plain-text response with information about the
// application status, operating environment and version.
func (app *Application) HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.Env,
			"version":     version,
		},
	}
	if err := app.writeJSON(w, http.StatusOK, data, nil); err != nil {
		app.logger.PrintError(err, nil)
		http.Error(w, "an error happened at the server level", http.StatusInternalServerError)
	}
}
