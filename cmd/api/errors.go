package main

import (
	"fmt"
	"net/http"
)

// the logError() method is a generic helper for logging an error message along
// with the current request method and URL as attributes in the log entry.
func (app *application) logError(r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

// The errorResponse method is a generic helper for sending json formatted error
// messages to the client with a given status code. Note that we're using the any
// type for the message parameter, rather than just a string type as this gives us
// more flexibility over the values that we can include in the response.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	envelope := envelope{"error": message}

	err := app.writeJSON(w, status, envelope, nil)

	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// The serverErrorRespons method will be used when our application encounters an
// unexpected prbolem at runtime. It logs the details error message, then uses the
// errorResponse helper to send a 500 status code and json response to the client.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "The server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// The notFoundResponse will be used to send 404 status codes and json responses
// to the client.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// the methodNotAllowed will be used to sned 405 status codes and json responses to the client.
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %v method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
