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
	innerErr error
	// Type is the type of dependency for which there was no
	// definition available in the IResolver
	Type reflect.Type
}

func newErrDefMissing(t reflect.Type) *ErrDefMissing {
	return &ErrDefMissing{
		innerErr: nil,
		Type:     t,
	}
}

func newErrDefMissingWrapper(innerErr error, t reflect.Type) *ErrDefMissing {
	return &ErrDefMissing{
		innerErr: innerErr,
		Type:     t,
	}
}

func (edm *ErrDefMissing) Error() string {
	if edm.innerErr == nil {
		return fmt.Sprintf("di: no definition found for %v in container", edm.Type)
	}

	return fmt.Sprintf("di: err resolving %v: %v", edm.Type, edm.innerErr)
}
