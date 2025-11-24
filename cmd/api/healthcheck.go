package main

import (
	"fmt"
	"net/http"
)

// Declare a handler which writes plain-text response with information about the
// application status, operating environment and vesion.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %v\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}
