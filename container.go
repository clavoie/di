package di

import (
	"fmt"
	"net/http"
	"reflect"
)

type IContainer interface {
	//Curry(fn interface{}) (interface{}, error)
	//Invoke(fn interface{}) error
	HttpHandler(fn interface{}) (func(http.ResponseWriter, *http.Request), error)
}

type container struct {
	allDeps    map[reflect.Type]*depNode
	deps       map[reflect.Type]*depNode
	perHttp    map[reflect.Type]*depNode
	perResolve map[reflect.Type]*depNode
	singletons map[reflect.Type]*singleton
}

func newContainer(cw *containerWriter) *container {
	deps := make(map[reflect.Type]*depNode)
	perHttp := make(map[reflect.Type]*depNode)
	perResolve := make(map[reflect.Type]*depNode)
	singletons := make(map[reflect.Type]*singleton)

	for rtype, node := range cw.deps {
		switch node.Lifetime {
		case LifetimeSingleton:
			singleton[rtype] = newSingleton(node)
		case LifetimePerDependency:
			deps[rtype] = node
		case LifetimePerHttpRequest:
			perHttp[rtype] = node
		case LifetimePerResolution:
			perResolve[rtype] = node
		}
	}

	return &container{
		allDeps:    cw.deps,
		deps:       deps,
		perHttp:    perHttp,
		perResolve: perResolve,
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

func (c *Container) HttpHandler(fn interface{}) (func(http.ResponseWriter, *http.Request), error) {
	err := c.verifyFn(fn)
	if err != nil {
		return nil, err
	}
}
