package di

// IResolver is an object which knows how to resolve dependency chains
// and instantiate the dependencies according to their cache policies
type IResolver interface {
	// Curry takes a func, resolves all parameters of the func which
	// are known to the container, and returns a new func with those
	// parameters supplied by the container.
	//
	// Explicitly:
	//   func foo(i int, dep Dep) int { ... }
	//   curryFoo, err := container.Curry(foo)
	//   if err { ... }
	//   val := curryFoo.(func(int) int)(4)
	Curry(fn interface{}) (interface{}, error)

	// Invoke resolves all known dependencies of func fn, and then attempts to execute the func.
	// If an error is encountered while resolving the dependencies of fn an error is returned.
	// If fn can be resolved and fn returns a single value which is an error, that is returned.
	// Otherwise nil is returned
	Invoke(fn interface{}) error

	// Resolve attempts to resolve a known dependency. The parameter must be a pointer to an interface
	// type known to the resolver
	//
	// Example:
	//   var dep Dep
	//   err := container.Resolve(&dep)
	Resolve(ptrToIface interface{}) error
}
