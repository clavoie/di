package di

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func resolverChildNew(defs []*Def) (IHttpResolver, error) {
	return NewResolver(func(er *ErrResolve, w http.ResponseWriter, r *http.Request) { panic(er) }, defs)
}

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
		resolver, err := resolverChildNew(nil)

		if err != nil {
			t.Fatal(err)
		}

		parent, isParent := resolver.(*resolverParent)
		if isParent == false {
			t.Fatal(reflect.TypeOf(resolver))
		}

		child := newResolverChild(parent)
		var self IResolver
		resolveErr := child.Resolve(&self)

		if resolveErr != nil {
			t.Fatal(resolveErr)
		}

		if child != self {
			t.Fatal(child, self)
		}
	})
	t.Run("Resolve", func(t *testing.T) {
		resolver, err := resolverChildNew([]*Def{
			&Def{NewA, PerDependency}, &Def{NewDependsOnHttp, PerDependency},
			&Def{NewSubDepNotFound, PerDependency},
		})

		if err != nil {
			t.Fatal(err)
		}

		t.Run("InvalidArg", func(t *testing.T) {
			err := resolver.Resolve("error")

			if err == nil {
				t.Fatal("expecting error when non interface ptr")
			}

			s := "error"
			err = resolver.Resolve(&s)

			if err == nil {
				t.Fatal("expecting error when non interface ptr")
			}

			type Invalid interface{}
			var invalid Invalid
			err = resolver.Resolve(&invalid)

			if err == nil {
				t.Fatal("expecting no def found err")
			}

			var doh DependsOnHttp
			err = resolver.Resolve(&doh)

			if err == nil {
				t.Fatal("expecting not being able to resolve http deps")
			}

			var sdnf SubDepNotFound
			err = resolver.Resolve(&sdnf)

			if err == nil {
				t.Fatal("expecting sub dep not found err")
			}
		})
	})
	t.Run("Curry", func(t *testing.T) {
		resolver, err := resolverChildNew([]*Def{&Def{NewA, PerDependency}})
		if err != nil {
			t.Fatal(err)
		}

		t.Run("Success", func(t *testing.T) {
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
		t.Run("InvalidFunc", func(t *testing.T) {
			_, err := resolver.Curry("hello")

			if err == nil {
				t.Fatal("expecting invalid fn err")
			}
		})
		t.Run("VariadicSuccess", func(t *testing.T) {
			expectedA := aCounter + 1
			fn := func(j string, a A, is ...string) (string, int) { return strings.Join(is, j), a.A() }
			ifn, err := resolver.Curry(fn)

			if err != nil {
				t.Fatal(err)
			}

			newFn, isCorrectType := ifn.(func(string, ...string) (string, int))

			if isCorrectType == false {
				t.Fatal(reflect.TypeOf(newFn))
			}

			s, i := newFn(",", "a", "b", "c")
			expectedS := "a,b,c"

			if s != expectedS {
				t.Fatal(s, expectedS)
			}

			if i != expectedA {
				t.Fatal(i, expectedA)
			}

			s, _ = newFn(",")
			if s != "" {
				t.Fatal(s)
			}
		})
	})
	t.Run("Invoke", func(t *testing.T) {
		resolver, err := resolverChildNew([]*Def{&Def{NewA, PerDependency}})

		if err != nil {
			t.Fatal(err)
		}

		t.Run("InvalidType", func(t *testing.T) {
			resolveErr := resolver.Invoke("hello")

			if resolveErr == nil {
				t.Fatal("expecting err")
			}
		})

		t.Run("InputParams", func(t *testing.T) {
			resolveErr := resolver.Invoke(func(s string) {})

			if resolveErr == nil {
				t.Fatal("expecting err")
			}
		})

		t.Run("ErrNil_NoErrResolve", func(t *testing.T) {
			resolveErr := resolver.Invoke(func() error { return nil })

			if resolveErr != nil {
				t.Fatal(resolveErr)
			}
		})
	})
}
