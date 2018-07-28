package di

import (
	"errors"
	"fmt"
	"reflect"
)

// duplicateDefErr is returned from verifyConstructor when a duplicate
// dependency definition is found. It signals Add to skip adding the
// definition to the Defs collection as it already exists
var duplicateDefErr = errors.New("di: duplicate definition")

// errType is the typeof(error)
var errType = reflect.TypeOf((*error)(nil)).Elem()

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

// Add adds a dependency definition to this Defs collection. See Def.Constructor
// for the format of the constructor parameter
func (d *Defs) Add(constructor interface{}, lifetime Lifetime) error {
	constructorValue := reflect.ValueOf(constructor)
	arg1, err := d.verifyConstructor(constructorValue, lifetime)

	if err != nil {
		if err == duplicateDefErr {
			return nil
		}

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
		err := d.Add(def.constructor, def.lifetime)

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
func Join(ds ...*Defs) *Defs {
	return &Defs{
		deps:   make(map[reflect.Type]*depNode),
		joined: ds,
	}
}

func (d *Defs) verifyConstructor(constructorValue reflect.Value, lifetime Lifetime) (reflect.Type, error) {
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

	existingDep, hasDep := d.deps[arg1]
	if hasDep {
		existing := fmt.Sprintf("%#v", existingDep.Constructor)
		newConstructor := fmt.Sprintf("%#v", constructorValue)
		if existing != newConstructor {
			return arg1, fmt.Errorf("di: a dependency for %v already exists with a different constructor:  %v, %v", arg1, existing, newConstructor)
		}

		if existingDep.Lifetime != lifetime {
			return arg1, fmt.Errorf("di: a dependency for %v already exists with a different lifetime: %v, %v", arg1, existingDep.Lifetime, lifetime)
		}

		return arg1, duplicateDefErr
	}

	if numOut == 2 {
		arg2 := constructorType.Out(1)

		if arg2.Implements(errType) == false {
			return arg1, fmt.Errorf("di: return value 2, if provided, must be an error: %v", arg2)
		}
	}

	return arg1, nil
}
