package di

import (
	"fmt"
	"net/http"
	"reflect"
)

type IContainer interface {
	Curry(fn interface{}) (interface{}, error)
	HttpHandler(fn interface{}, errFn func(error, http.ResponseWriter, *http.Request)) (func(http.ResponseWriter, *http.Request), error)
	Resolve(ptrToIface interface{}) error
}

type container struct {
	allDeps    map[reflect.Type]*depNode
	deps       map[reflect.Type]*depNode
	perHttp    map[reflect.Type]*depNode
	perResolve map[reflect.Type]*depNode
	singletons *resolveCache
}

func newContainer(cw *containerWriter) *container {
	deps := make(map[reflect.Type]*depNode)
	perHttp := make(map[reflect.Type]*depNode)
	perResolve := make(map[reflect.Type]*depNode)
	singletons := newResolveCache()

	for rtype, node := range cw.deps {
		switch node.Lifetime {
		case LifetimeSingleton:
			singletons.Set(rtype, newSingleton(node))
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

	return nil
}

func (c *container) Curry(fn interface{}) (interface{}, error) {
	fnValue := reflect.ValueOf(fn)
	err := c.verifyFn(fnValue)

	if err != nil {
		return nil, err
	}

	fnType := fnValue.Type()
	numIn := fnType.NumIn()
	isVariadic := fnType.IsVariadic()
	knowns := make([]bool, numIn)
	callTypes := make([]reflect.Type, 0, numIn)
	inVals := make([]reflect.Value, numIn)
	resolver := newResolver(c)

	for index := 0; index < numIn; index += 1 {
		inType := fnType.In(index)

		if index == numIn-1 && isVariadic {
			callTypes = append(callTypes, inType)
			continue
		}

		value, err := resolver.Resolve(inType)

		if err != nil {
			_, isErrDefMissing := err.(ErrDefMissing)

			if isErrDefMissing {
				callTypes = append(callTypes, inType)
				continue
			}

			return nil, err
		}

		knowns[index] = true
		inVals[index] = value
	}

	numOut := fnType.NumOut()
	outTypes := make([]reflect.Type, numOut)
	for index := 0; index < numOut; index += 1 {
		outTypes[index] = fnType.Out(index)
	}

	curryFnType := reflect.FuncOf(callTypes, outTypes, isVariadic)
	return reflect.MakeFunc(curryFnType, func(ins []reflect.Value) []reflect.Value {
		callVals := make([]reflect.Value, numIn)
		callIndex := 0
		for index := 0; index < numIn; index += 1 {
			if knowns[index] {
				callVals[index] = inVals[index]
			} else {
				callVals[index] = ins[callIndex]
				callIndex += 1
			}
		}

		return fnValue.Call(callVals)
	}).Interface(), nil
}

func (c *container) HttpHandler(fn interface{}, errFn func(error, http.ResponseWriter, *http.Request)) (func(http.ResponseWriter, *http.Request), error) {
	fnValue := reflect.ValueOf(fn)
	err := c.verifyFn(fnValue)

	if err != nil {
		return nil, err
	}

	fnType := fnValue.Type()
	numIn := fnType.NumIn()

	return func(w http.ResponseWriter, r *http.Request) {
		resolver := newHttpResolver(c, w, r)
		values := make([]reflect.Value, numIn)

		for index := range values {
			value, err := resolver.Resolve(fnType.In(index))

			if err != nil {
				errFn(err, w, r)
				return
			}

			values[index] = value
		}

		fnValue.Call(values)

		// TODO dispose of http stuff
	}, nil
}

func (c *container) Resolve(ptrToIface interface{}) error {
	ptrValue := reflect.ValueOf(ptrToIface)
	if ptrValue.Kind() != reflect.Ptr {
		return fmt.Errorf("di: ptrToIFace must be a *Interface type: %v", ptrValue.Type())
	}

	ifaceType := ptrValue.Type().Elem()
	if ifaceType.Kind() != reflect.Interface {
		return fmt.Errorf("di: ptrToIFace must be a *Interface type: %v", ptrValue.Type())
	}

	resolver := newResolver(c)
	value, err := resolver.Resolve(ifaceType)

	if err != nil {
		return err
	}

	reflect.Indirect(ptrValue).Set(value)
	return nil
}
