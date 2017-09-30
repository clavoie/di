package di

import (
	"fmt"
	"reflect"
)

type IContainer interface {
	Curry(fn interface{}) (interface{}, error)
	Invoke(fn interface{}) error
}

type container struct {
	deps       map[reflect.Type]*depNode
	singletons map[reflect.Type]*singleton
}

func newContainer(cw *containerWriter) *container {
	deps := make(map[reflect.Type]*depNode, len(cw.deps))
	singletons := make(map[reflect.Type]*singleton, len(cw.deps))

	for rtype, node := range cw.deps {
		switch node.Lifetime {
		case LifetimeSingleton:
			singleton[rtype] = newSingleton(node)
		case LifetimePerDependency:
			deps[rtype] = node
		case LifetimePerHttpRequest:
			panic("di: per http request not supported yet")
		}
	}

	return &container{
		deps:       deps,
		singletons: singletons,
	}
}

func (c *container) verifyFn(fnValue reflect.Value) error {
	if fnValue.Kind() != reflect.Func {
		return fmt.Errorf("di: constructor argument is not a function: %v", fnValue.Kind())

	}
}

func (c *Container) Curry(fn interface{}) (interface{}, error) {
	err := c.verifyFn(fn)

	if err != nil {
		return nil, err
	}
}

func (c *container) resolve(rtype reflect.Type) (reflect.Value, error) {

}

func (c *container) resolveCacheMiss(node *depNode) (reflect.Value, error) {
	if node.IsLeaf() {
		return node.NewValue([]reflect.Value{})
	}

	// TODO magic
}

func (c *container) resolveSingleton(rtype reflect.Type) (reflect.Value, error) {
	existingValue := c.singletons[rtype]

	if existingValue.Value.IsValid() {
		return existingValue.Value, nil
	}

	value, err := c.resolveCacheMiss(rtype)
	if err != nil {
		return value, err
	}

	existingValue.Value = value
	return value, nil
}
