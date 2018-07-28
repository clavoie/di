package di

// Def represents a dependency definition
type Def struct {
	// constructor is a func which instantiates the dependency
	constructor interface{}

	// lifetime is the caching lifetime of the dependency
	lifetime Lifetime
}

// NewDef creates a new dependency definition which can be added to a Defs collection.
//
// constructor must be a func of the signature:
//		func Name(dependency*) (Interface, error?)
// Examples:
//    func Foo1() Dependency
//    func Foo2(dep1 Dep1) Dependency
//    func Foo3(dep1, dep2 Dep1) (Dependency, error)
//
// lifetime is the caching lifetime of the dependency once it has been
// resolved
func NewDef(constructor interface{}, lifetime Lifetime) *Def {
	return &Def{
		constructor: constructor,
		lifetime:    lifetime,
	}
}
