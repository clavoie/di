package di

import (
	"fmt"
	"reflect"
	"strings"
)

// ErrResolve is returned when an attempt is made to resolve a type
// but an error is encountered while resolving the dependency. The
// error could either be returned from the dependency constructor
// or be because no definition for the requested type exists
type ErrResolve struct {
	// DependencyChain is the chain of types leading up to the
	// type that could not be resolved. Does not contain Type
	DependencyChain []reflect.Type

	// Err represents the error that caused the resolution error. If
	// a definition was missing this will be an *ErrDefMissing,
	// otherwise it will be whatever error was returned by the
	// dependency constructor
	Err error

	// Type is the type of dependency which could not be resolved
	Type reflect.Type
}

// newErrResolve rcreates an returns a new ErrResolve
func newErrResolve(depChain []reflect.Type, err error, t reflect.Type) *ErrResolve {
	if depChain == nil {
		depChain = make([]reflect.Type, 0, 1)
	}

	return &ErrResolve{
		DependencyChain: depChain,
		Err:             err,
		Type:            t,
	}
}

// String returns an string describing the error encountered
func (er *ErrResolve) String() string {
	chain := append(er.DependencyChain, er.Type)
	depNames := make([]string, len(chain))

	for index, dep := range chain {
		depNames[index] = dep.String()
	}

	depPath := strings.Join(depNames, " => ")
	if len(chain) == 1 {
		depPath = "[]"
	}

	return fmt.Sprintf("di: could not resolve type %v in path: %v, err: %v", er.Type, depPath, er.Err)
}
