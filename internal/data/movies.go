package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the QueryRow() method to execute the sql query on our connection pool
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)

	// You did not have to create an args array of course. We could pass the
	// placeholder values like so.
	// return m.DB.QueryRow(query, movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

// func (m MovieModel) Get(id int64) (*Movie, error) {
// 	if id < 1 {
// 		return nil, ErrRecordNotFound
// 	}

// 	const query = `
// 		SELECT
// 			id,
// 			created_at,
// 			title,
// 			year,
// 			runtime,
// 			genres,
// 			version
// 		FROM
// 			movies
// 		WHERE
// 			id = $1;`

// 	var movie Movie

// 	err := m.DB.QueryRow(query, id).Scan(
// 		&movie.ID,
// 		&movie.CreatedAt,
// 		&movie.Title,
// 		&movie.Year,
// 		&movie.Runtime,
// 		pq.Array(&movie.Genres),
// 		&movie.Version,
// 	)

// 	// so we also check if the error is actually a no rows found error
// 	// this way we can send a record not found response
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, ErrRecordNotFound
// 		}

// 		return nil, err
// 	}

// 	return &movie, nil
// }

// Example to show cancellation of long running sql queries using context.
func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// Add pg_sleep(8) to the select statement such that we wait 8 seconds
	// before getting the response back from the db.
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

	// Use context.WithTimeout() to create a context.Context which carries a
	// 3 second timeout deadline.
	// Note that we're using the empty context.Background() as the parent context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// Use defer to make sure we cancel the context before the Get()
	// method returns
	defer cancel()

	// Use the QueryRowContext() method to execute the query, passing in the context
	// with the deadline as the first argument.
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
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
	// Change the update query to include the version number to avoid data races.
	// We call this approach optimistic locking.
	// If we can't find a record with a matching id and version number we will return
	// an error. Since this means that the record has been updated since it was read.
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
		AND
			version = $6
		RETURNING 
			version;`

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)

	if err != nil {
		// We will return this error if no record was found.
		// This will mean that since we initially read the record
		// it was updated and has a new version number and therefore
		// a data race has occurred.
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}

		return err
	}

	return nil
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sqlRes, err := m.DB.ExecContext(ctx, query, id)

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

func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, *Metadata, error) {
	query := fmt.Sprintf(`
		SELECT
			COUNT(*) OVER(),
			id,
			created_at,
			title,
			year,
			runtime,
			genres,
			version
		FROM
			Movies
		WHERE 
			($1::text IS NULL OR $1::text = '' OR LOWER(title) = LOWER($1::text))
		AND 
			($2::text[] IS NULL OR array_length($2::text[], 1) = 0 OR genres @> $2::text[])
		ORDER BY %v %v
		LIMIT $3
		OFFSET $4;`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		title,
		pq.Array(genres),
		filters.limit(),
		filters.offset(),
	}

	rows, err := m.DB.QueryContext(ctx, query, args...)

	// You must check the error before defering rows.Close()
	if err != nil {
		return nil, &Metadata{}, err
	}

	defer rows.Close()

	// here we are creating a slice of pointers to movies
	// it is differnt from this movies := &[]Movie{}
	// which is a pointer to a movie slice i.e. the pointer
	// points to a slice that contains movies.
	movies := []*Movie{}

	totalRecords := 0

	// rows.Next() will return true if there is a row available to be read
	// if there is a row, it will read it and move the pointer to the next row.
	for rows.Next() {
		var movie Movie

		err := rows.Scan(
			&totalRecords,
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)

		if err != nil {
			return nil, &Metadata{}, err
		}

		// add the movie to the slice
		movies = append(movies, &movie)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, &Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return movies, &metadata, nil
}
