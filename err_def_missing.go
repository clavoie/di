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
	// Type is the type of dependency for which there was no definition available
	// in the IResolver
	Type reflect.Type
}

func newErrDefMissing(t reflect.Type) ErrDefMissing {
	return ErrDefMissing{t}
}

func (edm ErrDefMissing) Error() string {
	return fmt.Sprintf("di: no definition found for %v in container", edm.Type)
}
