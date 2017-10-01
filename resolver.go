package di

import (
	"fmt"
	"net/http"
	"reflect"
)

var resolverNoCache = map[reflect.Type]*singleton{}

var requestType = reflect.TypeOf((**http.Request)(nil)).Elem()
var responseWriterType = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()

type resolver struct {
	container  *container
	perDep     map[reflect.Type]*depNode
	perHttp    map[reflect.Type]*singleton
	perResolve map[reflect.Type]*singleton
}

func newResolver(c *container) *resolver {
	return &resolver{
		container:  c,
		perDep:     c.deps,
		perHttp:    make(map[reflect.Type]*singleton),
		perResolve: make(map[reflect.Type]*singleton),
	}
}

func newHttpResolver(c *container, w http.ResponseWriter, r *http.Request *resolver) {
	resolver := newResolver(c)

	resolver.perHttp[requestType] = newSingletonValue(reflect.ValueOf(r))
	resolver.perHttp[responseWriterType] = newSingletonValue(reflect.ValueOf(w))

	return resolver
}

func (r *resolver) Resolve(rtype reflect.Type) (reflect.Value, error) {
	dep, hasDep := r.container.allDeps[rtype]
	if hasDep == false {
		return reflect.Value{}, fmt.Errorf("di: cannot resolve %v, no definition found", rtype)
	}

	cache := r.lifetimeToCache(dep.Lifetime)
	cacheValue, hasCacheValue := cache[rtype]
	if hasCacheValue == false {
		cacheValue = newSingleton(dep)

		if cache != resolverNoCache {
			cache[rtype] = cacheValue
		}
	}

	value, hasValue := cacheValue.Value()
	if hasValue {
		return value, nil
	}

	return r.resolveNoCache(dep, cacheValue)
}

func (r *resolver) lifetimeToCache(l Lifetime) map[reflect.Type]*singleton {
	switch l {
	case LifetimeSingleton:
		return r.container.singleton
	case LifetimePerDependency:
		return resolverNoCache
	case LifetimePerHttpRequest:
		return r.perHttp
	case LifetimePerResolution:
		return r.perResolve
	}
}

func (r *resolver) resolveNoCache(node *depNode, s *singleton) (reflect.Value, error) {
	if node.IsLeaf() {
		return s.SetValue([]reflect.Value{})
	}

	values := make([]reflect.Value, len(node.DependsOn))
	for index, dep := range node.DependsOn {
		value, err := r.Resolve(dep)

		if err != nil {
			return reflect.Value{}, err
		}

		values[index] = value
	}

	return s.SetValue(values)
}
