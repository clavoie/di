package di

import "reflect"

// resolverNoCache is an instance of resolveCache indicating no
// caching should take place
var resolverNoCache = newResolveCache()

// resolveCache is a cache of values instantiated along the
// dependency chain
type resolveCache struct {
	cache map[reflect.Type]*singleton
}

func newResolveCache() *resolveCache {
	return &resolveCache{
		cache: make(map[reflect.Type]*singleton),
	}
}

func (rc *resolveCache) Get(rtype reflect.Type) (*singleton, bool) {
	s, hasValue := rc.cache[rtype]
	return s, hasValue
}

func (rc *resolveCache) Set(rtype reflect.Type, value *singleton) {
	if rc == resolverNoCache {
		return
	}

	rc.cache[rtype] = value
}
