package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance
	router := httprouter.New()

	// Convert the notfoundresponse to a http.Handler using
	// http.HandlerFunc adapter and then set it as the custom error handler for 404
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// Convert the methodNotAllowedResponse to a http.Handler using
	// http.HandlerFunc adapter and then set it as the custom error handler for 405
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.getMovieByIdHandler)
	// router.HandlerFunc(http.MethodPut, "/v1/movies/:id", app.updateMovieHandler)

	// We change the allowed HTTP verb to patch since we are performing a partial update
	// i.e. we may not necessarily update the entire record but only parts of it.
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMoviesHandler)

	// We wrap our router with the panic recovery middleware.
	// This will ensure that the middleware runs for every one of our API endpoints.
	return app.recoverPanic(router)
	// return router
}
