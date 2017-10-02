package di

import (
	"net/http"
	"reflect"
)

type resolver struct {
	container  *container
	perDep     map[reflect.Type]*depNode
	perHttp    *resolveCache
	perResolve *resolveCache
}

func newResolver(c *container) *resolver {
	return &resolver{
		container:  c,
		perDep:     c.deps,
		perHttp:    newResolveCache(),
		perResolve: newResolveCache(),
	}
}

func newHttpResolver(c *container, w http.ResponseWriter, r *http.Request) *resolver {
	resolver := newResolver(c)

	resolver.perHttp.Set(requestType, newSingletonValue(reflect.ValueOf(r)))
	resolver.perHttp.Set(responseWriterType, newSingletonValue(reflect.ValueOf(w)))

	return resolver
}

func (r *resolver) Resolve(rtype reflect.Type) (reflect.Value, error) {
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

func (r *resolver) lifetimeToCache(l Lifetime) *resolveCache {
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
