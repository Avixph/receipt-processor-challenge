package main

import (
	"net/http"
)

// Handler that writes a json response with the application status,
// operating environment and version number.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	jsnEnv := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}
	err := app.writeJSON(w, http.StatusOK, jsnEnv, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
