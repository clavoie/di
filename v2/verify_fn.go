package di

import (
	"fmt"
	"reflect"
)

// verifyFn asserts that a value provided to one of the resolvers is
// actually a function type
func verifyFn(fnValue reflect.Value) error {
	if fnValue.Kind() != reflect.Func {
		return fmt.Errorf("di: constructor argument is not a function: %v", fnValue.Kind())

	}

	return nil
}
