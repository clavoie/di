package di

// Lifetime indicates the caching policy for resolved types
type Lifetime int

const (
	// Singleton indicates only one instance of the type
	// should be created ever, and used for every dependency
	// encountered going forward
	Singleton Lifetime = iota

	// PerDependency indicates that a new instance of the
	// type should be created for each dependency encountered.
	//
	// Explicitly:
	//		func NewFoo(dep1, dep2 Dep) Foo
	//
	// dep1 and dep2 will be two separate instances of Dep
	PerDependency

	// PerHttpRequest indicates that a new instance of the type
	// should be created per http request.
	//
	// Explicitly:
	//		func NewFoo(dep1, dep2 Dep) Foo
	//
	// if Foo is resolved in an http request dep1 and dep2 will be
	// the exact same instance of Dep. If not resolved via an http
	// request PerHttprequest acts like PerResolve
	PerHttpRequest

	// PerResolve indicates that a new instance of the type should
	// be created per call to Resolve(), but that the same instance
	// of the type should be used throughout the Resolve() call
	//
	// Explicitly:
	//		func NewFoo(dep1, dep2 Dep) Foo
	//		err1 := container.Resolve(&foo1)
	//		err2 := container.Resolve(&foo2)
	//		foo1.dep1 == foo1.dep2
	//		foo1.dep1 != foo2.dep1
	PerResolve
)

// lifetimes is a collection of all known Lifetime values
var lifetimes = map[Lifetime]bool{
	Singleton:      true,
	PerDependency:  true,
	PerHttpRequest: true,
	PerResolve:     true,
}
