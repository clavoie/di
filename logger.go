package di

import "time"

// ILogger is a logging interface that a client can implement to
// capture output and metrics on the performance of a di resolver.
//
// When a resolver starts up it will look for a definition of
// ILogger. If present the hooks defined in ILogger will be called back
// on the implementation.
type ILogger interface {
	// HttpDuration is called after all the dependencies for an http
	// handler have been resolved, but before the http handler runs.
	// The parameter value is the time taken to resolve all the handlers
	// dependencies.
	HttpDuration(time.Duration)
}
