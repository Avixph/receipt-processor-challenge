package main

import (
	"fmt"
	"net/http"
)

// logError() helps with logging error messages along with the current request method and URL as attributes in the entry.
func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.RequestURI
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

// errorResponse() helps with sending JSON-formatted error messages to the client with a
// given status code.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	jsnEnv := envelope{"error": message}
	err := app.writeJSON(w, status, jsnEnv, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// serverErrorResponse() is used when the application encounters an unexpected problem
// at runtime. It logs the detailed error message, then uses the errorResponse() helper
// to send a 500 Internal Server Error status code and JSON response (containing a
// generic error message) to the client.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// notFoundResponse() will be used to send a 404 Not Found status code and
// JSON response to the client.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowedResponse() method will be used to send a 405 Method Not Allowed
// status code and JSON response to the client.
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// badRequestResponse() method will be used to send a 400 Bad Request status code and
// JSON response to the client.
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// failedValidationResponse() method writes a 422 Unprocessable Entity and the contents of
// the errors map from our new Validator type as a JSON response body.
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
