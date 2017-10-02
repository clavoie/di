package di

import (
	"fmt"
	"net/http"
	"reflect"
)

var iresolverType = reflect.TypeOf((*IResolver)(nil)).Elem()

type resolverChild struct {
	container  *resolverParent
	closables  []IHttpClosable
	perDep     map[reflect.Type]*depNode
	perHttp    *resolveCache
	perResolve *resolveCache
}

func newResolverChild(c *resolverParent) *resolverChild {
	resolver := &resolverChild{
		container:  c,
		closables:  make([]IHttpClosable, 0),
		perDep:     c.deps,
		perHttp:    newResolveCache(),
		perResolve: newResolveCache(),
	}

	resolver.perResolve.Set(iresolverType, newSingletonValue(reflect.ValueOf(resolver)))

	return resolver
}

func newHttpResolverChild(c *resolverParent, w http.ResponseWriter, r *http.Request) *resolverChild {
	resolver := newResolverChild(c)

	resolver.perHttp.Set(requestType, newSingletonValue(reflect.ValueOf(r)))
	resolver.perHttp.Set(responseWriterType, newSingletonValue(reflect.ValueOf(w)))

	return resolver
}

func (r *resolverChild) Curry(fn interface{}) (interface{}, error) {
	fnValue := reflect.ValueOf(fn)
	err := verifyFn(fnValue)

	if err != nil {
		return nil, err
	}

	fnType := fnValue.Type()
	numIn := fnType.NumIn()
	isVariadic := fnType.IsVariadic()
	knowns := make([]bool, numIn)
	callTypes := make([]reflect.Type, 0, numIn)
	inVals := make([]reflect.Value, numIn)

	for index := 0; index < numIn; index += 1 {
		inType := fnType.In(index)

		if index == numIn-1 && isVariadic {
			callTypes = append(callTypes, inType)
			continue
		}

		value, err := r.resolveCache(inType)

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

func (r *resolverChild) Resolve(ptrToIface interface{}) error {
	ptrValue := reflect.ValueOf(ptrToIface)
	if ptrValue.Kind() != reflect.Ptr {
		return fmt.Errorf("di: ptrToIFace must be a *Interface type: %v", ptrValue.Type())
	}

	ifaceType := ptrValue.Type().Elem()
	if ifaceType.Kind() != reflect.Interface {
		return fmt.Errorf("di: ptrToIFace must be a *Interface type: %v", ptrValue.Type())
	}

	value, err := r.resolveCache(ifaceType)

	if err != nil {
		return err
	}

	reflect.Indirect(ptrValue).Set(value)
	return nil
}

func (r *resolverChild) lifetimeToCache(l Lifetime) *resolveCache {
	switch l {
	case Singleton:
		return r.container.singletons
	case PerHttpRequest:
		return r.perHttp
	case PerResolution:
		return r.perResolve
	}

	return resolverNoCache
}

func (r *resolverChild) resolveCache(rtype reflect.Type) (reflect.Value, error) {
	dep, hasDep := r.container.allDeps[rtype]
	if hasDep == false {
		return reflect.Value{}, newErrDefMissing(rtype)
	}

	cache := r.lifetimeToCache(dep.Lifetime)
	cacheValue, hasCacheValue := cache.Get(rtype)
	if hasCacheValue == false {
		cacheValue = newSingleton(dep)
		cache.Set(rtype, cacheValue)
	}

	value, hasValue := cacheValue.Value()
	if hasValue {
		return value, nil
	}

	return r.resolveNoCache(dep, cacheValue)
}

func (r *resolverChild) resolveNoCache(node *depNode, s *singleton) (reflect.Value, error) {
	if node.IsLeaf() {
		return s.SetValue([]reflect.Value{}, &r.closables)
	}

	values := make([]reflect.Value, len(node.DependsOn))
	for index, dep := range node.DependsOn {
		value, err := r.resolveCache(dep)

		if err != nil {
			return reflect.Value{}, err
		}

		values[index] = value
	}

	return s.SetValue(values, &r.closables)
}
