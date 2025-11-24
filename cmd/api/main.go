package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var config config

	// Read the value of the port and env command-line flags into the config struct.
	// We default to using the port number 4000 and environment 'development' if no
	// corresponding flags are provided.
	flag.IntVar(&config.port, "port", 4000, "API server port")
	flag.StringVar(&config.env, "env", "development", "Environment(development|staging|production)")
	flag.Parse()

	// Initialize a new structured logger which writes log entries to the standard out stream.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Declare an instance of the application struct containing the config struct and the logger.
	app := &application{
		config: config,
		logger: logger,
	}

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

	err := server.ListenAndServe()

	// If the server errors out send the error to the structured logger
	logger.Error(err.Error())

	// Exit gracefully
	os.Exit(1)
}
