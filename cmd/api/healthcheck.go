package main

import (
	"net/http"
)

// Declare a handler which writes plain-text response with information about the
// application status, operating environment and vesion.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// USING Interpolated string
	// js := fmt.Sprintf(`{"status": "available", "environment": "development", "version": %q}`, version)

	// // Add the Content-Type header such that we return json
	// w.Header().Set("Content-Type", "application/json")

	// // Write the byte array out to the response stream
	// w.Write([]byte(js))

	// USING json.Marshal()

	// Create a map that represents our data (this could have been a struct)
	envelope := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, envelope, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
