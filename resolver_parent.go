package di

import (
	"net/http"
	"reflect"
	"time"
)

// iloggerType is typeof(ILogger)
var iloggerType = reflect.TypeOf((*ILogger)(nil)).Elem()

// resolverParent is a type which contains all the combined
// dependency definitions, and created new resolverChild
// types to handle Resolve() requests.
type resolverParent struct {
	allDeps    map[reflect.Type]*depNode
	deps       map[reflect.Type]*depNode
	hasLogger  bool
	perHttp    map[reflect.Type]*depNode
	perResolve map[reflect.Type]*depNode
	singletons *resolveCache
}

// NewResolver returns a new instance of IHttpResolver from
// a collection of dependency definitions. An IResolver definition is added
// to the collection, and can be included as a dependency for resolved types
// and funcs
func NewResolver(d *Defs) (IHttpResolver, error) {
	allDeps, err := d.build()

	if err != nil {
		return nil, err
	}

	numDeps := len(allDeps)
	deps := make(map[reflect.Type]*depNode, numDeps/4)
	hasLogger := false
	perHttp := make(map[reflect.Type]*depNode, numDeps/4)
	perResolve := make(map[reflect.Type]*depNode, numDeps/4)
	singletons := newResolveCache()

	for rtype, node := range allDeps {
		if rtype == iloggerType {
			hasLogger = true
		}

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
		hasLogger:  hasLogger,
		perHttp:    perHttp,
		perResolve: perResolve,
		singletons: singletons,
	}, nil
}

func (c *resolverParent) Curry(fn interface{}) (interface{}, *ErrResolve) {
	resolver := newResolverChild(c)
	return resolver.Curry(fn)
}

func (c *resolverParent) HttpHandler(fn interface{}, errFn func(*ErrResolve, http.ResponseWriter, *http.Request)) (func(http.ResponseWriter, *http.Request), error) {
	fnValue := reflect.ValueOf(fn)
	err := verifyFn(fnValue)

	if err != nil {
		return nil, err
	}

	fnType := fnValue.Type()
	numIn := fnType.NumIn()

	return func(w http.ResponseWriter, r *http.Request) {
		var epoch time.Time

		if c.hasLogger {
			epoch = time.Now()
		}

		resolver := newHttpResolverChild(c, w, r)
		values := make([]reflect.Value, numIn)

		for index := range values {
			value, err := resolver.resolveUsingCache(nil, fnType.In(index))

			if err != nil {
				errFn(err, w, r)
				return
			}

			values[index] = value
		}

		for _, closable := range resolver.closables {
			defer closable.Di_HttpClose()
		}

		if c.hasLogger {
			var logger ILogger
			err := resolver.Resolve(&logger)

			if err != nil {
				errFn(err, w, r)
				return
			}

			logger.HttpDuration(time.Since(epoch))
		}

		fnValue.Call(values)
	}, nil
}

func (c *resolverParent) Invoke(fn interface{}) *ErrResolve {
	resolver := newResolverChild(c)
	return resolver.Invoke(fn)
}

func (c *resolverParent) Resolve(ptrToIface interface{}) *ErrResolve {
	resolver := newResolverChild(c)
	return resolver.Resolve(ptrToIface)
}
