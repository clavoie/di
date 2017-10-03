package di

import (
	"reflect"
	"strings"
	"testing"
)

func TestErrDefMissing(t *testing.T) {
	a := NewA()
	aType := reflect.TypeOf(a)
	edm := newErrDefMissing(aType)
	es := edm.Error()
	ts := aType.String()

	if strings.Contains(es, ts) == false {
		t.Fatal("was expecting the types String() in the error message", ts)
	}

	var err error
	err = edm
	if err == nil {
		t.Fatal(err)
	}
}
