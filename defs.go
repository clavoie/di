package di

import (
	"fmt"
	"reflect"
)

// errType is the typeof(error)
var errType = reflect.TypeOf((*error)(nil)).Elem()

// Def represents a dependency definition
type Def struct {
	// Constructor must be a func of the signature:
	//		func Name(dependency*) (Interface, error?)
	// Examples:
	//    func Foo1() Dependency
	//    func Foo2(dep1 Dep1) Dependency
	//    func Foo3(dep1, dep2 Dep1) (Dependency, error)
	Constructor interface{}

	// Lifetime is the caching lifetime of the dependency once it has been
	// resolved
	Lifetime Lifetime
}

// NewDef creates a new dependency definition which can be added to a Defs collection. See Def.Constructor for
// the format of the constructor parameter
func NewDef(constructor interface{}, lifetime Lifetime) *Def {
	return &Def{
		Constructor: constructor,
		Lifetime:    lifetime,
	}
}

// Defs represents a collection of dependency definitions
type Defs struct {
	deps   map[reflect.Type]*depNode
	joined []*Defs
}

// NewDefs creates a new Defs collection
func NewDefs() *Defs {
	return &Defs{
		deps:   make(map[reflect.Type]*depNode),
		joined: make([]*Defs, 0),
	}
}

// Add adds a dependency definition to this Defs collection. The
// supplied constructor must have at least 1 return value,
// the first of which is an interface. The optional second return
// value may be an error
func (d *Defs) Add(constructor interface{}, lifetime Lifetime) error {
	constructorValue := reflect.ValueOf(constructor)
	arg1, err := d.verifyConstructor(constructorValue)

	if err != nil {
		return err
	}

	if lifetimes[lifetime] == false {
		return fmt.Errorf("di: unknown lifetime: %v", lifetime)
	}

	newNode := newDepNode(constructorValue, lifetime, d.deps)
	d.deps[arg1] = newNode
	for _, node := range d.deps {
		node.AddEdge(newNode)
	}

	return nil
}

// AddAll is a bulk version of Add
func (d *Defs) AddAll(defs []*Def) error {
	for _, def := range defs {
		err := d.Add(def.Constructor, def.Lifetime)

		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Defs) all() []*depNode {
	deps := make([]*depNode, 0, len(d.deps)+len(d.joined))

	for _, dep := range d.deps {
		deps = append(deps, dep)
	}

	for _, defs := range d.joined {
		deps = append(deps, defs.all()...)
	}

	return deps
}

func (d *Defs) build() (map[reflect.Type]*depNode, error) {
	allDeps := d.all()
	finalDeps := &Defs{
		deps: make(map[reflect.Type]*depNode, len(allDeps)),
	}

	for _, dep := range allDeps {
		err := finalDeps.Add(dep.Constructor.Interface(), dep.Lifetime)

		if err != nil {
			return nil, err
		}
	}

	checked := make(map[*depNode]bool, len(d.deps))

	for _, node := range finalDeps.deps {
		if checked[node] {
			continue
		}

		err := node.CheckForCycle([]*depNode{}, checked)

		if err != nil {
			return nil, err
		}
	}

	return finalDeps.deps, nil
}

// Join combines two Defs collections together into a new Defs
func (d1 *Defs) Join(d2 *Defs) *Defs {
	return &Defs{
		deps:   make(map[reflect.Type]*depNode),
		joined: []*Defs{d1, d2},
	}
}

func (d *Defs) verifyConstructor(constructorValue reflect.Value) (reflect.Type, error) {
	var arg1 reflect.Type

	if constructorValue.Kind() != reflect.Func {
		return arg1, fmt.Errorf("di: constructor argument is not a function: %v", constructorValue.Kind())
	}

	constructorType := constructorValue.Type()
	numOut := constructorType.NumOut()
	if numOut == 0 || numOut > 2 {
		return arg1, fmt.Errorf("di: constructor can return exactly 1 or 2 values")
	}

	arg1 = constructorType.Out(0)
	if arg1.Implements(errType) {
		return arg1, fmt.Errorf("di: return value 1 cannot be an error: %v", arg1)
	}

	if arg1.Kind() != reflect.Interface {
		return arg1, fmt.Errorf("di: return value 1 must be an interface: %v", arg1)
	}

	_, hasDep := d.deps[arg1]
	if hasDep {
		return arg1, fmt.Errorf("di: a dependency for %v already exists", arg1)
	}

	if numOut == 2 {
		arg2 := constructorType.Out(1)

		if arg2.Implements(errType) == false {
			return arg1, fmt.Errorf("di: return value 2, if provided, must be an error: %v", arg2)
		}
	}

	return arg1, nil
}
