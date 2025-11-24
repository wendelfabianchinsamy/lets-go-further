# PROJECT STRUCTURE

## bin

- The bin directory will contain our compiled application binaries, ready for deployment to the production server.

## cmd/api

- The cmd/api directory will contain the application-specific code for our API application.
- This will include the code for running the server, reading and writing HTTP requests and managing authentication.

## internal

- The internal directory will contain various ancillary packages used by our API.
- It will contain code for interacting with our database, doing data validation, sending emails and so on.
- Basically, any code which isn't application-specific and can potentially be reused will live in here.
- Our Go code under cmd/api will import the packages in the internal directory but the internal directory code will never import code from the cmd/api directory
- Any code that resides in the internal directory should not be imported by any other application.

## migrations

- The migrations directory will contain the SQL migration files for out database.

## remote

- The remove directory will contain the configuration files and setup scripts for our production server.

## go.mod

- The go.mod file will declare our project dependencies, versions and module path.

## Makefile

- the Makefile will contain recipes for automating common administrative tasks - like auditing our God code, building binaries and executing database migrations
