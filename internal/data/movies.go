package data

import "time"

type Movie struct {
	ID int64 `json:"id"`
	// CreatedAt time.Time `json:"created_at"`
	CreatedAt time.Time `json:"-"` // by using the hyphen directive we can hide the CreatedAt field from being sent in the request
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitzero"`    // omitzero basically says don't encode this field if the value is the zero value (can and probably should be done with a pointer)
	Runtime   Runtime   `json:"runtime,omitzero"` // use the custom Runtime type so we get the custom marhsalling logic
	Genres    []string  `json:"genres,omitzero"`
	Version   int32     `json:"version"`
}
