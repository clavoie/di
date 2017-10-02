package di

import "reflect"

var resolverNoCache = newResolveCache()

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
