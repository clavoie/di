package di

type IResolver interface {
	// Curry takes a func, resolves all parameters of the func which
	// are known to the container, and returns a new func with those
	// parameters supplied by the container.
	//
	// Explicitly:
	//   func foo(i int, dep Dep) int { ... }
	//   curryFoo, err := container.Curry(foo)
	//   if err { ... }
	//   val := (func(int) int)(curryFoo)(4)
	Curry(fn interface{}) (interface{}, error)

	// Resolve attempts to resolve a known dependency.
	//
	// Example:
	//   var dep Dep
	//   err := container.Resolve(&dep)
	Resolve(ptrToIface interface{}) error
}
