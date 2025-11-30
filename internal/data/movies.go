package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/wendelfabianchinsamy/lets-go-further/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	// CreatedAt time.Time `json:"-"` // by using the hyphen directive we can hide the CreatedAt field from being sent in the request
	Title   string   `json:"title"`
	Year    int32    `json:"year,omitzero"`    // omitzero basically says don't encode this field if the value is the zero value (can and probably should be done with a pointer)
	Runtime Runtime  `json:"runtime,omitzero"` // use the custom Runtime type so we get the custom marhsalling logic
	Genres  []string `json:"genres,omitzero"`
	Version int32    `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(len(movie.Genres) >= 1, "genres", "must have at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must have no more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicates")
}

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type MovieModel struct {
	DB *sql.DB
}

// The insert method accepts a pointer to a movie struct which should contain the
// data for the new record.
func (m MovieModel) Insert(movie *Movie) error {
	// Define a sql query for inserting a new record in the movies table and returning
	// the system-generated data.
	const query = `
		INSERT INTO movies (
			title,
			year,
			runtime,
			genres)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
		RETURNING 
			id,
			created_at,
			version`

	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our sql query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	// Use the QueryRow() method to execute the sql query on our connection pool
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)

	// You did not have to create an args array of course. We could pass the
	// placeholder values like so.
	// return m.DB.QueryRow(query, movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	const query = `
		SELECT
			id,
			created_at,
			title,
			year,
			runtime,
			genres,
			version
		FROM
			movies
		WHERE
			id = $1;`

	var movie Movie

	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	// so we also check if the error is actually a no rows found error
	// this way we can send a record not found response
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}

		return nil, err
	}

	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	const query = `
		UPDATE 
			movies
		SET 
			title = $1,
			year = $2,
			runtime = $3,
			genres = $4,
			version = version + 1
		WHERE 
			id = $5
		RETURNING 
			version;`

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.Version,
	}

	return m.DB.QueryRow(query, args...).Scan(&movie.Version)
}

func (m MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	const query = `
		DELETE FROM
			movies
		WHERE
			id = $1;`

	sqlRes, err := m.DB.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := sqlRes.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
