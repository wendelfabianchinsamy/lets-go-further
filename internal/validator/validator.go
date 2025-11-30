package validator

import (
	"regexp"
	"slices"
)

// Regular expression that defines what is a valid email addresss.
var EmailRegEx = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

// Define a new Validator type which contains a map of validation errors
type Validator struct {
	Errors map[string]string
}

// New is a helper which creates a new Validator isntance with an empty errors map.
// This is basically a constructor.
func New() *Validator {
	return &Validator{Errors: map[string]string{}}
}

// Valid returns true if the errors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the map provided it
// does not already exist in the map.
func (v *Validator) AddError(key string, value string) {
	_, exists := v.Errors[key]

	if !exists {
		v.Errors[key] = value
	}
}

// Check adds an error message to the map only if a validation check is not ok.
func (v *Validator) Check(ok bool, key string, value string) {
	if !ok {
		v.AddError(key, value)
	}
}

// Generic function which returns true if a specific value is in a list of permitted values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches returns true if a string value matches a specific regexp pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Generic function which returns true if all values in a slice are unique.
func Unique[T comparable](values []T) bool {
	uniqueValues := map[T]bool{}

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
