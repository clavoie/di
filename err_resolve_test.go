package di

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestErrResolve(t *testing.T) {
	type someDep struct{}
	var depInst someDep
	depType := reflect.TypeOf(depInst)
	depStr := fmt.Sprintf("%v", depType)

	t.Run("String", func(t *testing.T) {
		err := errors.New("some_err")

		resolveErr := newErrResolve(nil, err, depType)
		str := resolveErr.String()

		if strings.Contains(str, "[]") == false {
			t.Fatal("was expecting empty path")
		}

		if strings.Contains(str, err.Error()) == false {
			t.Fatal("was expecting err string")
		}

		if strings.Contains(str, depStr) == false {
			t.Fatal("was expecting type name")
		}

		resolveErr = newErrResolve([]reflect.Type{depType}, err, depType)
		str = resolveErr.String()
		expected := depStr + " => " + depStr

		if strings.Contains(str, expected) == false {
			t.Fatal("expecting dep chain path", str, expected)
		}
	})
}
