package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Create a new custom error that will reprsent any runtime value
// that does not match the required format which is "<runtime> mins".
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

// Declare a custom runtime type which has the underlying type int32
type Runtime int32

// Implement a MarshalJSON() method on the Runtime tyep so that it satisfies the
// json.Marshaler interface. This should return the JSON-encoded value for the movie
// runtime.
func (r Runtime) MarshalJSON() ([]byte, error) {
	// Generate a string containing the movie runtime in the required format.
	jsonValue := fmt.Sprintf("%d mins", r)

	// Use strconv.Quote() on the string to wrap it in double quotes. It
	// needs to be surrounded by double quotes in order to be a valid *JSON string*
	quotedJSONValue := strconv.Quote(jsonValue)

	// Convert the quoted string value to a byte slice and return it.
	return []byte(quotedJSONValue), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	// We expect the incoming json value will be a string in the format
	// "<runtime> mins" and the first thing we need to do is remvoe the surrounding
	// double quotes from this string. If we can't unquote it then we return the
	// ErrInvalidRuntimeFormat error.
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))

	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Split the string to isolate teh part containing the number
	parts := strings.Split(unquotedJSONValue, " ")

	// Sanity check the parts of the string to make sure it was in the expected format.
	// If it isn't we return the ErrInvalidRuntimeFormat error again.
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	// Otherwise parse the string containing the number into an int32. Again, if this
	// fails return the ErrInvalidRuntimeFormat error.
	i, err := strconv.ParseInt(parts[0], 10, 32)

	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Convert the int32 to a Runtime type and assign this to the receiver. Note that
	// we use the * operator to dereference the receiver (which is a pointer to a Runtime type)
	// in order to set the underlying value to the pointer.
	*r = Runtime(i)

	return nil
}
