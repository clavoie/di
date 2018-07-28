package di

// Def represents a dependency definition
type Def struct {
	// Constructor is a func which instantiates the dependency
	// Must be a func of the signature:
	//		func Name(dependency*) (Interface, error?)
	// Examples:
	//    func Foo1() Dependency
	//    func Foo2(dep1 Dep1) Dependency
	//    func Foo3(dep1, dep2 Dep1) (Dependency, error)
	Constructor interface{}

	// Lifetime is the caching Lifetime of the dependency once
	// it has been resolved
	Lifetime Lifetime
}

// NewDef creates a new dependency definition which can be added to a Defs collection.
func NewDef(constructor interface{}, lifetime Lifetime) *Def {
	return &Def{
		Constructor: constructor,
		Lifetime:    lifetime,
	}
}
