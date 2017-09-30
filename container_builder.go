package di

import (
	"fmt"
	"reflect"
)

var errType = reflect.TypeOf((*error)(nil)).Elem()

type IContainerWriter interface {
	Add(constructor interface{}, lifetime Lifetime) error
	Build() (IContainer, error)
}

func NewContainerWriter() IContainerWriter {
	return &containerWriter{
		deps: make(map[reflect.Type]*depNode, 10),
	}
}

type containerWriter struct {
	deps map[reflect.Type][]reflect.Type
}

func (cw *containerWriter) Add(constructor interface{}, lifetime Lifetime) error {
	constructorValue := reflect.ValueOf(constructor)
	arg1, err := cw.verifyConstructor(constructorValue)

	if lifetimes[lifetime] == false {
		return fmt.Errorf("di: unknown lifetime: %v", lifetime)
	}

	newNode := newDepNode(constructorValue, lifetime)
	cw.deps[arg1] = newNode
	for _, node := range cw.deps {
		node.AddEdge(newNode)
	}

	return nil
}

func (cw *containerWriter) Build() (IContainer, error) {
	checked := make(map[*depNode]bool, len(cw.deps))

	for _, node := range cw.deps {
		if checked[node] {
			continue
		}

		err := node.CheckForCycle([]*depNode{}, checked)

		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (cw *containerWriter) join(cw2 *containerWriter) (*containerWriter, error) {
	cw3 := &containerWriter{
		deps: make(map[reflect.Type]*depNode, len(cw.deps)+len(cw2.deps)),
	}

	for node := range cw.deps {
		err := cw3.Add(node.Constructor, node.Lifetime)

		if err != nil {
			return nil, err
		}
	}

	for node := range cw2.deps {
		err := cw3.Add(node.Constructor, node.Lifetime)

		if err != nil {
			return nil, err
		}
	}

	return cw3, nil
}

func (cw *containerWriter) verifyConstructor(constructorValue reflect.Value) (reflect.Type, error) {
	var arg1, arg2 reflect.Type

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
