package di

import "net/http"

// IHttpResolver is an IResolver which can also generate http request
// handlers that resolve their dependencies
type IHttpResolver interface {
	IResolver

	// HttpHandler creates a new http request handler from a fn containing
	// dependencies. The ResponseWriter and *Request are supplied as
	// dependencies of the container, and will be resolved in the supplied
	// func or one of its dependencies. errFn is an error handling func
	// which will be called if there is an err while resolving one of the
	// dependencies
	HttpHandler(fn interface{}, errFn func(error, http.ResponseWriter, *http.Request)) (func(http.ResponseWriter, *http.Request), error)
}
