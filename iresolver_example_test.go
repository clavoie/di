package di

import (
	"fmt"
	"reflect"
)

func ExampleIResolver_curry() {
	type Dep interface{}
	newDep := func() Dep { return new(struct{}) }
	defs := NewDefs()
	err := defs.Add(newDep, PerDependency)
	if err != nil {
		panic(err)
	}

	resolver, err := NewResolver(defs)
	if err != nil {
		panic(err)
	}

	fn := func(msg string, dep Dep) string {
		return fmt.Sprintf("%v:%v %v", msg, reflect.TypeOf(dep), dep == nil)
	}

	ifn, err := resolver.Curry(fn)
	if err != nil {
		panic(err)
	}

	var newFn func(string) string
	newFn = ifn.(func(string) string)

	fmt.Println(newFn("type"))
	// Output: type:*struct {} false
}

func ExampleIResolver_invoke() {
	type Dep interface{}
	newDep := func() Dep { return new(struct{}) }
	defs := NewDefs()
	err := defs.Add(newDep, PerDependency)
	if err != nil {
		panic(err)
	}

	resolver, err := NewResolver(defs)
	if err != nil {
		panic(err)
	}

	err = resolver.Invoke(func(dep Dep) {
		fmt.Println(reflect.TypeOf(dep))
		// Output: *struct {}
	})

	if err != nil {
		panic(err)
	}

	myErr := fmt.Errorf("my error")
	err = resolver.Invoke(func(dep Dep) error { return myErr })

	if err != myErr {
		panic(myErr)
	}
}

func ExampleIResolver_resolve() {
	type Dep interface{}
	newDep := func() Dep { return new(struct{}) }
	defs := NewDefs()
	err := defs.Add(newDep, PerDependency)
	if err != nil {
		panic(err)
	}

	resolver, err := NewResolver(defs)
	if err != nil {
		panic(err)
	}

	var dep Dep
	err = resolver.Resolve(&dep)
	if err != nil {
		panic(err)
	}

	fmt.Println(dep == nil, reflect.TypeOf(dep))

	var resolver2 IResolver
	err = resolver.Resolve(&resolver2)
	if err != nil {
		panic(err)
	}

	fmt.Println(resolver2 == nil)
	// Output: false *struct {}
	// false
}
