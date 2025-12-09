package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/wendelfabianchinsamy/lets-go-further/internal/validator"
)

// A type for representing our envelope data
// We will not send data directly to the client but instead wrap the json in another json object
// Giving it a descriptive name.
type envelope map[string]any

func (app *application) readIdParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)

	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	// js, err := json.Marshal(data)

	if err != nil {
		return err
	}

	// Header() returns the header map
	// So we are simply appending headers to the existing header map
	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(r *http.Request, destination any) error {
	// err := json.NewDecoder(r.Body).Decode(destination)

	// If we only want to allow the fields defined in the anonymous struct
	// we defined below as part of the request body, we can call DisallowUnknownFields()
	// on the decoder.
	// BUT IS THIS REALLY THAT MUCH OF A REQUIREMENT???
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(destination)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does, then return a plain english error message
		// which includes the location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// In some circumstancs Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax errors in the JSON. So we check for this using errors.Is() and
		// return a generic error message.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains badly-formed JSON")
		// Likewise, catch any json.UnmarshalTypeErrors. These occur when the
		// JSON value is the wrong type for the target destination. If the error relates
		// to a specific field, then we include that in our error message to make it
		// easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}

			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		// An io.EOF error will be returned by Decode() if the request body is empty. We
		// check for this with errors.Is() and return a plain english error message instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		// A json.InvalidUnmarshalError error will be returned if we pass something
		// that is not a non-nil pointer to Decode(). We catch this and panic rather
		// than returning an error to our handler.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		// Use the error.As() function to check whether the error has th type
		// *http.MaxBytesError. If it does, then it means the request boyd exceeded our
		// size limit of 1mb and we return a clear error message.
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		// For anything else, return the error message as is.
		default:
			app.logger.Info("inside readJSON")
			return err
		}
	}

	return nil
}

// The readString() helper returns a string value from the query string or provided default value
// if no matching key could be found.
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

// The readCSV() helper reads a string value from the query string and then splits it
// into a slice on the comma character. If no matching key could be found it returns
// the provided default value
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

// The readInt() helper reads a string value from the query string and converts that value
// into an int. If the conversion fails we add a validation error to the pointer validation
// instance and we return the provided default value. If parsing was successful then we return
// the parsed int.
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)

	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}
