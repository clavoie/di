# di [![GoDoc Reference](https://img.shields.io/badge/GoDoc-Reference-blue.svg)](https://godoc.org/github.com/clavoie/di) ![Build Status](https://travis-ci.org/clavoie/di.svg?branch=master) [![codecov](https://codecov.io/gh/clavoie/di/branch/master/graph/badge.svg)](https://codecov.io/gh/clavoie/di) [![Go Report Card](https://goreportcard.com/badge/github.com/clavoie/di)](https://goreportcard.com/report/github.com/clavoie/di)

Dependency injection for Go

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
  defs := di.NewDefs()
  // ...
  
  resolver, err := di.NewResolver(defs)
  // if err
  
  for _, handler := range handlers {
    // ResponseWriter and *Request are available as dependencies
    // handler.fn => func(dep1 Dep1, dep2 Dep2, etc)
    httpFn, err := resolver.HttpHandler(handler.fn, errFn)
    // if err
    
    http.HandleFunc(handler.url, httpFn)
  }
```

## Types
```go
  var resolver di.Resolver
  // ...
  
  var someDependency Dependency
  err := resolver.Resolve(&someDependency)
```

## Curry Funcs
```go
  var resolver di.IResolver
  // ...
  
  func normalFunc(name string, someDependency Dependency) (int, string) {
    return someDependency.DoThing(name), name + "!"
  }
  
  ifunc, err := resolver.Curry(normalFunc)
  if err != nil {
  return err
  }
  
  value, msg := ((func (string)(int, string))(ifunc))("hello")
```
