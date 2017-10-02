package di

import (
	"fmt"
	"reflect"
)

func verifyFn(fnValue reflect.Value) error {
	if fnValue.Kind() != reflect.Func {
		return fmt.Errorf("di: constructor argument is not a function: %v", fnValue.Kind())

	}

	return nil
}
