package di

import (
	"net/http"
	"reflect"
)

// resolverParent is a type which contains all the combined
// dependency definitions, and created new resolverChild
// types to handle Resolve() requests.
type resolverParent struct {
	allDeps    map[reflect.Type]*depNode
	deps       map[reflect.Type]*depNode
	perHttp    map[reflect.Type]*depNode
	perResolve map[reflect.Type]*depNode
	singletons *resolveCache
}

// NewResolver returns a new instance of IHttpResolver from
// a collection of dependency definitions
func NewResolver(d *Defs) (IHttpResolver, error) {
	allDeps, err := d.build()

	if err != nil {
		return nil, err
	}

	numDeps := len(allDeps)
	deps := make(map[reflect.Type]*depNode, numDeps/4)
	perHttp := make(map[reflect.Type]*depNode, numDeps/4)
	perResolve := make(map[reflect.Type]*depNode, numDeps/4)
	singletons := newResolveCache()

	for rtype, node := range allDeps {
		switch node.Lifetime {
		case Singleton:
			singletons.Set(rtype, newSingleton(node))
		case PerDependency:
			deps[rtype] = node
		case PerHttpRequest:
			perHttp[rtype] = node
		case PerResolve:
			perResolve[rtype] = node
		}
	}

	return &resolverParent{
		allDeps:    allDeps,
		deps:       deps,
		perHttp:    perHttp,
		perResolve: perResolve,
		singletons: singletons,
	}, nil
}

func (c *resolverParent) Curry(fn interface{}) (interface{}, error) {
	resolver := newResolverChild(c)
	return resolver.Curry(fn)
}

func (c *resolverParent) HttpHandler(fn interface{}, errFn func(error, http.ResponseWriter, *http.Request)) (func(http.ResponseWriter, *http.Request), error) {
	fnValue := reflect.ValueOf(fn)
	err := verifyFn(fnValue)

	if err != nil {
		return nil, err
	}

	fnType := fnValue.Type()
	numIn := fnType.NumIn()

	return func(w http.ResponseWriter, r *http.Request) {
		resolver := newHttpResolverChild(c, w, r)
		values := make([]reflect.Value, numIn)

		for index := range values {
			value, err := resolver.resolveCache(fnType.In(index))

			if err != nil {
				errFn(err, w, r)
				return
			}

			values[index] = value
		}

		fnValue.Call(values)

		for _, closable := range resolver.closables {
			closable.Di_HttpClose()
		}
	}, nil
}

func (c *resolverParent) Resolve(ptrToIface interface{}) error {
	resolver := newResolverChild(c)
	return resolver.Resolve(ptrToIface)
}
