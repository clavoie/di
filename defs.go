package di

import (
	"fmt"
	"reflect"
)

// errType is the typeof(error)
var errType = reflect.TypeOf((*error)(nil)).Elem()

// Defs represents a collection of dependency type definitions
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

func (cw *Defs) build() (map[reflect.Type]*depNode, error) {
	allDeps := cw.all()
	finalDeps := &Defs{
		deps: make(map[reflect.Type]*depNode, len(allDeps)),
	}

	for _, dep := range allDeps {
		err := finalDeps.Add(dep.Constructor.Interface(), dep.Lifetime)

		if err != nil {
			return nil, err
		}
	}

	checked := make(map[*depNode]bool, len(cw.deps))

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
func (cw *Defs) Join(cw2 *Defs) *Defs {
	return &Defs{
		deps:   make(map[reflect.Type]*depNode),
		joined: []*Defs{cw, cw2},
	}
}

func (cw *Defs) verifyConstructor(constructorValue reflect.Value) (reflect.Type, error) {
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

	_, hasDep := cw.deps[arg1]
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
