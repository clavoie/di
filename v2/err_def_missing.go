package di

import (
	"fmt"
	"reflect"
)

// ErrDefMissing is returned when an attempt is made to resolve
// a type but no definition for the type was found in the resolver.
//
// Implements the error interface
type ErrDefMissing struct {
	// Type is the type of dependency which could not be resolved
	Type reflect.Type
}

// newErrDefMissing creates and returns a new ErrDefMissing struct
func newErrDefMissing(t reflect.Type) *ErrDefMissing {
	return &ErrDefMissing{
		Type: t,
	}
}

// Error returns an error string describing the error encountered
func (edm *ErrDefMissing) Error() string {
	return fmt.Sprintf("di: definition missing for type: %v", edm.Type)
}
