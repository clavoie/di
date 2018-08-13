# di [![GoDoc Reference](https://img.shields.io/badge/GoDoc-Reference-blue.svg)](https://godoc.org/github.com/clavoie/di) ![Build Status](https://travis-ci.org/clavoie/di.svg?branch=master) [![codecov](https://codecov.io/gh/clavoie/di/branch/master/graph/badge.svg)](https://codecov.io/gh/clavoie/di) [![Go Report Card](https://goreportcard.com/badge/github.com/clavoie/di)](https://goreportcard.com/report/github.com/clavoie/di)

di is a dependency injection framework for Go. di supplies several dependency lifetime caching policies, provides dependency aware http handlers compatible with net/http, and provides a way to clean up dependencies instantiated during an http request.

di only resolves dependencies which are interfaces, the resolver itself, http.ResponseWriter, and *http.Request.

## Dependency Lifetimes
- Singleton
  - only one instance created throughout the life of the resolver
- PerDependency
  - a new instance created for each dependency encountered, every time Resolve() is called
- PerHttpRequest
  - a new instance is created for each HTTP request, and reused for the duration of the request. The dependency has an option to implement an interface which will be called back once the request is over
- PerResolve
  - a new instance is created for each call of Resolve(), and then re-used throughout that call. Subsequent calls to Resolve() create new instances
  
[More about lifetimes](https://godoc.org/github.com/clavoie/di#Lifetime)

## Http
```go
  dependencies := []*di.Def{
    &di.Def{SomeConstructor, di.PerHttpRequest},
    // etc...
  }

  httpDefs := []*di.HttpDef{
    // the http.ResponseWriter and *http.Request values are available as dependencies,
    // the resolver is also available as a dependency as an di.IResolver 
    // SomeHandler => func(dep1 Dep1, dep2 Dep2, etc) 
    &di.HttpDef{SomeHandler, "/some/pattern"},
    // etc...
  }

  resolver, err := di.NewResolver(errHandler, dependencies)
  // if err
  
  err = resolver.SetDefaultServeMux(httpDefs)
  // if err

  log.Fatal(http.ListenAndServe(":8080", nil))
```
[A more complete example is available here](https://godoc.org/github.com/clavoie/di#example-IHttpResolver)

## Types
di can resolve a dependency directly if known. The dependency instance follows the lifecycle caching rules of the
resolver
```go
  var someDependency Dependency
  resolveErr := resolver.Resolve(&someDependency)
```

## Curry Funcs
di can curry the parameters of funcs with dependencies known to a resolver, returning a new func that only contains
parameters the caller would like to supply themselves.
```go
  func normalFunc(name string, dep Dep) (int, string) {
    // DoThing(name) == 5
    return dep.DoThing(name), name + "!"
  }
  
  ifunc, resolveErr := resolver.Curry(normalFunc)
  // if resolveErr
  
  value, msg := ifunc.(func (string)(int, string))("hello")
  // 5, "hello!"
```

## Invoke
```go
  resolveErr := resolver.Invoke(func (dep Dep){
    // do()
  })
```

If the func returns an error, and no error is encountered while resolving the dependencies, that error is returned instead

```go
  resolveErr := resolver.Invoke(func (dep Dep) error {
     innerErr := dep1.Do()
     return innerErr
  })
  // resolveErr.Err == innerErr
```
