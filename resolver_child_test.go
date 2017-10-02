package di

import (
	"reflect"
	"testing"
)

func TestResolverChild(t *testing.T) {
	t.Run("Implements_IResolver", func(t *testing.T) {
		child := new(resolverChild)

		var resolver IResolver
		resolver = child
		_, err := resolver.Curry(func() {})

		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("InjectsSelf", func(t *testing.T) {
		defs := NewDefs()
		resolver, err := NewResolver(defs)

		if err != nil {
			t.Fatal(err)
		}

		parent, isParent := resolver.(*resolverParent)
		if isParent == false {
			t.Fatal(reflect.TypeOf(resolver))
		}

		child := newResolverChild(parent)
		var self IResolver
		err = child.Resolve(&self)

		if err != nil {
			t.Fatal(err)
		}

		if child != self {
			t.Fatal(child, self)
		}
	})
	t.Run("Curry", func(t *testing.T) {
		deps := NewDefs()
		err := deps.Add(NewA, PerDependency)

		if err != nil {
			t.Fatal(err)
		}

		resolver, err := NewResolver(deps)
		if err != nil {
			t.Fatal(err)
		}

		expectedA := aCounter + 1
		fn := func(s string, a A) (int, string) {
			return a.A() + 1, s
		}
		ifn, err := resolver.Curry(fn)

		if err != nil {
			t.Fatal(err)
		}

		newFn, isCorrectType := ifn.(func(string) (int, string))
		if isCorrectType == false {
			t.Fatal(reflect.TypeOf(ifn))
		}

		const sval = "test"
		i, s := newFn(sval)

		if sval != s {
			t.Fatal(sval, s)
		}

		if i != expectedA+1 {
			t.Fatal(i, expectedA+1)
		}
	})
}
