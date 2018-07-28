package di

// HttpDef is an injectable go net/http handler definition
type HttpDef struct {
	// Handler is the handler for the http request. All parameters
	// for Handler will be injected
	Handler interface{}

	// Pattern is the URL pattern used to match the request. See go's
	// documentation for DefaultServeMux
	Pattern string
}
