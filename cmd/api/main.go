package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	// Import the pq driver so that it can register itself with the database/sql
	// package. Note that we alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	_ "github.com/lib/pq"
	"github.com/wendelfabianchinsamy/lets-go-further/internal/data"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
		// the number of open connections the database can have.
		// by default the max number of open connections is unlimited.
		maxOpenConns int

		// the number of idle connections in the pool.
		// by default this value is 2.
		maxIdleConns int

		// the duration a connection can be idle.
		maxIdleTime time.Duration
	}
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
}

func main() {
	var config config

	// Read the value of the port and env command-line flags into the config struct.
	// We default to using the port number 4000 and environment 'development' if no
	// corresponding flags are provided.
	flag.IntVar(&config.port, "port", 4000, "API server port")
	flag.StringVar(&config.env, "env", "development", "Environment(development|staging|production)")

	// Read the dsn value from the db-dsn command line flag into the config struct. We
	// default to using our development DSN if no flag is provided.
	flag.StringVar(&config.db.dsn, "db-dsn", "postgres://wendel:password@localhost/letsgofurther?sslmode=disable", "PostgresSQL DSN")

	// Read the connection pool settings from the command line flags and config struct.
	// Notice that the default values we're using are the ones we discussed above.
	flag.IntVar(&config.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&config.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&config.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()

	// Initialize a new structured logger which writes log entries to the standard out stream.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(config)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Defer a call to db.Close() so that the connection pool is closed before the
	// main function exits.
	defer db.Close()

	// Declare an instance of the application struct containing the config struct and the logger.
	app := &application{
		config: config,
		logger: logger,
		models: data.NewModels(db),
	}

	app.logger.Info("database connection pool established")

	// Declare a HTTP server which listens on the port provided in the config struct, uses
	// the servemux we created above as the handler, has some sensible timeout settings
	// and writes any log messages to the structured logger at Error level.
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Start the HTTP server
	logger.Info(fmt.Sprintf("Starting server %v %v", server.Addr, config.env))

	err = server.ListenAndServe()

	// If the server errors out send the error to the structured logger
	logger.Error(err.Error())

	// Exit gracefully
	os.Exit(1)
}

func openDB(config config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool using the dsn from the config struct
	db, err := sql.Open("postgres", config.db.dsn)

	if err != nil {
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(config.db.maxOpenConns)

	// Set the maximum number of idle connections in the pool. Again passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(config.db.maxIdleConns)

	// Set the maximum idle timeout for connections in the pool. Passing a duration less
	// than or equal to 0 will mean that connections are not closed due to their idle time.
	db.SetConnMaxIdleTime(config.db.maxIdleTime)

	// Create a context with a 5 second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// an example of running the application while setting the flags
	// go run ./cmd/api -db-max-open-conns=50 -db-max-idle-conns=50 -db-max-idle-time=2h30m

	// Use PingContext() to establish a new connection to the database pasing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within a 5 second deadline, then this will return an
	// error. If we get this error, or any other, we close the connection pool and
	// return the error.
	err = db.PingContext(ctx)

	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
