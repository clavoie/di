package di_test

import (
	"net/http"

	"github.com/clavoie/di"
)

type HttpDep interface {
	// optional: only dependencies that perform cleanup need to implement
	// IHttpClosable
	di.IHttpClosable
}

type HttpImpl struct{}

func (hi *HttpImpl) Di_HttpClose() {
	// called after the request has ended
}

func NewHttpDep() HttpDep { return new(HttpImpl) }

func NewHttpResolver() di.IHttpResolver {
	defs := di.NewDefs()

	// one instance of HttpDep created per http request
	err := defs.Add(NewHttpDep, di.PerHttpRequest)

	if err != nil {
		panic(err)
	}

	resolver, err := di.NewResolver(defs)
	if err != nil {
		panic(err)
	}

	return resolver
}

var resolver = NewHttpResolver()

func DepHandler(dep HttpDep)                         {}
func HttpHandler(http.ResponseWriter, *http.Request) {}
func WriteToLog(err error)                           { /* etc */ }

func ErrHandler(err error, w http.ResponseWriter, r *http.Request) {
	WriteToLog(err)
	w.WriteHeader(http.StatusInternalServerError)
}

var urlDefs = []struct {
	url     string
	handler interface{}
}{
	{"/", DepHandler},
	{"/thing", HttpHandler},
}

func ExampleIHttpResolver_httpHandler() {
	for _, urlDef := range urlDefs {
		handleFunc, err := resolver.HttpHandler(urlDef.handler, ErrHandler)
		if err != nil {
			panic(err)
		}

		http.HandleFunc(urlDef.url, handleFunc)
	}

	// listen and serve
}
