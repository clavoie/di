package di

// Def represents a dependency definition
type IDef interface {
	// Constructor returns the constructor value for this IDef
	Constructor() interface{}

	// Lifetime returns the lifetime value of this IDef
	Lifetime() Lifetime
}

// Def represents a dependency definition
type Def struct {
	// constructor must be a func of the signature:
	//		func Name(dependency*) (Interface, error?)
	// Examples:
	//    func Foo1() Dependency
	//    func Foo2(dep1 Dep1) Dependency
	//    func Foo3(dep1, dep2 Dep1) (Dependency, error)
	constructor interface{}

	// lifetime is the caching lifetime of the dependency once it has been
	// resolved
	lifetime Lifetime
}

// NewDef creates a new dependency definition which can be added to a Defs collection. See Def.Constructor for
// the format of the constructor parameter
func NewDef(constructor interface{}, lifetime Lifetime) *Def {
	return &Def{
		constructor: constructor,
		lifetime:    lifetime,
	}
}

func (d *Def) Constructor() interface{} { return d.constructor }
func (d *Def) Lifetime() Lifetime       { return d.lifetime }
