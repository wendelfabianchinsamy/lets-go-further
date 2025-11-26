package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/wendelfabianchinsamy/lets-go-further/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

func (app *application) getMovieByIdHandler(w http.ResponseWriter, r *http.Request) {
	// panic("foobar") use this to test recoverPanic() middleware
	id, err := app.readIdParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	envolope := envelope{
		"movie": data.Movie{
			ID:        id,
			CreatedAt: time.Now(),
			Title:     "Casablanca",
			Runtime:   102,
			Genres:    []string{"drama", "romance", "war"},
			Version:   1,
		},
	}

	err = app.writeJSON(w, http.StatusOK, envolope, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
