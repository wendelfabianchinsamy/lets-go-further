package data

import (
	"fmt"
	"strconv"
)

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
