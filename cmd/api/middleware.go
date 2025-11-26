package main

import (
	"fmt"
	"net/http"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deffered function which will alway be run in the event of a panic
		defer func() {
			// Use the builtin recover function to check if there has been a panic
			err := recover()

			if err != nil {
				// If there was a panic, set the connection: close header on the
				// response. This acts as a trigger to make Go's http server
				// automatically close the current connection after a response has been
				// sent.
				w.Header().Set("Connection", "close")

				// The value returned by recover has th type any so we use
				// fmt.ErrorF to normalize it into an error and call our
				// serverErrorResponse helper. In turn, this will log the error using
				// our custom logger type at the ERROR level and send the client a 500
				// internal server error response.
				app.serverErrorResponse(w, r, fmt.Errorf("%v", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
