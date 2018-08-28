package di_test

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/clavoie/di"
)

func ExampleIResolver_curry() {
	type Dep interface{}
	newDep := func() Dep { return new(struct{}) }

	resolver, err := di.NewResolver(
		func(er *di.ErrResolve, w http.ResponseWriter, r *http.Request) { panic(er) },
		[]*di.Def{
			{newDep, di.PerDependency},
		})
	if err != nil {
		panic(err)
	}

	fn := func(msg string, dep Dep) string {
		return fmt.Sprintf("%v:%v %v", msg, reflect.TypeOf(dep), dep == nil)
	}

	ifn, resolveErr := resolver.Curry(fn)
	if resolveErr != nil {
		panic(resolveErr)
	}

	var newFn func(string) string
	newFn = ifn.(func(string) string)

	fmt.Println(newFn("type"))
	// Output: type:*struct {} false
}

func ExampleIResolver_invoke() {
	type Dep interface{}
	newDep := func() Dep { return new(struct{}) }
	resolver, err := di.NewResolver(
		func(er *di.ErrResolve, w http.ResponseWriter, r *http.Request) { panic(er) },
		[]*di.Def{
			{newDep, di.PerDependency},
		})
	if err != nil {
		panic(err)
	}

	resolveErr := resolver.Invoke(func(dep Dep) {
		fmt.Println(reflect.TypeOf(dep))
	})

	if resolveErr != nil {
		panic(resolveErr)
	}

	myErr := fmt.Errorf("my error")
	resolveErr = resolver.Invoke(func(dep Dep) error { return myErr })

	if resolveErr.Err != myErr {
		panic(myErr)
	}
	// Output: *struct {}
}

func ExampleIResolver_resolve() {
	type Dep interface{}
	newDep := func() Dep { return new(struct{}) }
	resolver, err := di.NewResolver(
		func(er *di.ErrResolve, w http.ResponseWriter, r *http.Request) { panic(er) },
		[]*di.Def{
			{newDep, di.PerDependency},
		})
	if err != nil {
		panic(err)
	}

	var dep Dep
	resolveErr := resolver.Resolve(&dep)
	if resolveErr != nil {
		panic(resolveErr)
	}

	fmt.Println(dep == nil, reflect.TypeOf(dep))

	var resolver2 di.IResolver
	resolveErr = resolver.Resolve(&resolver2)
	if resolveErr != nil {
		panic(resolveErr)
	}

	fmt.Println(resolver2 == nil)
	// Output: false *struct {}
	// false
}
