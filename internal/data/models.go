package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error for when we do not find a record
// in the database.
var ErrRecordNotFound = errors.New("record not found")

// Define a custom ErrEditConflict error for when we have a data race when
// attempting to mutate a record (update or delete).
var ErrEditConflict = errors.New("edit conflict")

// Create a Models struct which wraps the MovieModel. We'll add other models to this
// like a UserModel and PermissionModel as our build progresses.
type Models struct {
	Movies MovieModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}
