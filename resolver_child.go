package di

import (
	"fmt"
	"net/http"
	"reflect"
)

// errorType is typeof(error)
var errorType = reflect.TypeOf((*error)(nil)).Elem()

// iresolverType is typeof(IResolver)
var iresolverType = reflect.TypeOf((*IResolver)(nil)).Elem()

// resolverChild is an object which resolves a dependency chain for
// one call of Resolve().
//
// resolverChild injects itself into the IResolver, and can be resolved
// by any dependencies as IResolver
type resolverChild struct {
	parent     *resolverParent
	closables  []IHttpClosable
	perDep     map[reflect.Type]*depNode
	perHttp    *resolveCache
	perResolve *resolveCache
}

// newResolverChild returns a new resolverChild. IResolver is mapped
// to the instance of this object
func newResolverChild(c *resolverParent) *resolverChild {
	resolver := &resolverChild{
		parent:     c,
		closables:  make([]IHttpClosable, 0),
		perDep:     c.deps,
		perHttp:    newResolveCache(),
		perResolve: newResolveCache(),
	}

	resolver.perResolve.Set(iresolverType, newSingletonValue(reflect.ValueOf(resolver)))

	return resolver
}

// newHttpResolverChild returns a new resolverChild which has the
// http.ResponseWriter and *http.Request mapped for injection into
// dependencies
func newHttpResolverChild(c *resolverParent, w http.ResponseWriter, r *http.Request) *resolverChild {
	resolver := newResolverChild(c)

	resolver.perHttp.Set(requestType, newSingletonValue(reflect.ValueOf(r)))
	resolver.perHttp.Set(responseWriterType, newSingletonValue(reflect.ValueOf(w)))

	return resolver
}

func (r *resolverChild) Curry(fn interface{}) (interface{}, *ErrResolve) {
	fnValue := reflect.ValueOf(fn)
	err := verifyFn(fnValue)

	if err != nil {
		return nil, newErrResolve(nil, err, reflect.TypeOf(fn))
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

		value, err := r.resolveUsingCache(nil, inType)

		if err != nil {
			_, isErrDefMissing := err.Err.(*ErrDefMissing)

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

		if isVariadic {
			lastVal := callVals[numIn-1]
			lastValLen := lastVal.Len()
			callVals = callVals[:numIn-1]

			for index := 0; index < lastValLen; index += 1 {
				callVals = append(callVals, lastVal.Index(index))
			}
		}

		return fnValue.Call(callVals)
	}).Interface(), nil
}

func (r *resolverChild) Invoke(fn interface{}) *ErrResolve {
	newFn, err := r.Curry(fn)

	if err != nil {
		return err
	}

	fnValue := reflect.ValueOf(newFn)
	fnType := reflect.TypeOf(newFn)
	if fnType.NumIn() > 0 {
		return newErrResolve(nil, fmt.Errorf("di: Invoke: cannot invoke a func with input parameters: %v", fnType.NumIn()), fnType)
	}

	outValues := fnValue.Call([]reflect.Value{})
	if fnType.NumOut() == 1 && fnType.Out(0) == errorType {
		errVal := outValues[0].Interface()

		if errVal != nil {
			return newErrResolve(nil, errVal.(error), fnType)
		}
	}

	return nil
}

func (r *resolverChild) Resolve(ptrToIface interface{}) *ErrResolve {
	ptrValue := reflect.ValueOf(ptrToIface)
	if ptrValue.Kind() != reflect.Ptr {
		return newErrResolve(nil, fmt.Errorf("di: ptrToIFace must be a *Interface type: %v", ptrValue.Type()), ptrValue.Type())
	}

	ifaceType := ptrValue.Type().Elem()
	if ifaceType.Kind() != reflect.Interface {
		return newErrResolve(nil, fmt.Errorf("di: ptrToIFace must be a *Interface type: %v", ptrValue.Type()), ptrValue.Type())
	}

	value, err := r.resolveUsingCache(nil, ifaceType)

	if err != nil {
		return err
	}

	reflect.Indirect(ptrValue).Set(value)
	return nil
}

// lifetimeToCache maps a Lifetime to one of the variosu caches of the
// resolver, returning the cache
func (r *resolverChild) lifetimeToCache(l Lifetime) *resolveCache {
	switch l {
	case Singleton:
		return r.parent.singletons
	case PerHttpRequest:
		return r.perHttp
	case PerResolve:
		return r.perResolve
	}

	return resolverNoCache
}

// resolveUsingCache attempts to resolve a value for a type using this
// resolver's cache. ErrDefMissing is returned if there is no
// definition in this resolver for the specified type
func (r *resolverChild) resolveUsingCache(depChain []reflect.Type, rtype reflect.Type) (reflect.Value, *ErrResolve) {
	if rtype == iresolverType {
		return reflect.ValueOf(r), nil
	}

	if rtype == requestType || rtype == responseWriterType {
		httpValue, hasValue := r.perHttp.Get(rtype)

		if hasValue == false {
			return reflect.Value{}, newErrResolve(depChain, newErrDefMissing(rtype), rtype)
		}

		value, hasValue := httpValue.Value()
		if hasValue == false {
			return reflect.Value{}, newErrResolve(depChain, newErrDefMissing(rtype), rtype)
		}

		return value, nil
	}

	dep, hasDep := r.parent.allDeps[rtype]
	if hasDep == false {
		return reflect.Value{}, newErrResolve(depChain, newErrDefMissing(rtype), rtype)
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

	return r.resolveIgnoringCache(depChain, dep, cacheValue)
}

// resolveIgnoringCache is called on a resolve cache miss. It attempts to
// resolve the missing type, and set the cache of the type which is
// missing with the instantiated value
func (r *resolverChild) resolveIgnoringCache(depChain []reflect.Type, node *depNode, s *singleton) (reflect.Value, *ErrResolve) {
	if node.IsLeaf() {
		value, err := s.SetValue([]reflect.Value{}, &r.closables)

		if err != nil {
			return reflect.Value{}, newErrResolve(depChain, err, node.Type)
		}

		return value, nil
	}

	values := make([]reflect.Value, len(node.DependsOn))
	childDepChain := append(depChain, node.Type)
	for index, dep := range node.DependsOn {
		value, err := r.resolveUsingCache(childDepChain, dep)

		if err != nil {
			return reflect.Value{}, err
		}

		values[index] = value
	}

	value, err := s.SetValue(values, &r.closables)
	if err != nil {
		return reflect.Value{}, newErrResolve(depChain, err, node.Type)
	}

	return value, nil
}
