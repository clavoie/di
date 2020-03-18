package di

// IHttpClosable is an interface a dependency can implement if they
// would like a callback executed when an http request finishes
type IHttpClosable interface {
	// Di_HttpClose is called when an http request, in which the
	// implementing object was instantiated, completes
	Di_HttpClose()
}
