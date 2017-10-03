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

func ExampleIHttpResolver_httpHandler() {
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

	logger := func(error) {}
	errHandler := func(err error, w http.ResponseWriter, r *http.Request) {
		logger(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	depHandler := func(dep HttpDep) {}
	normalHandler := func(http.ResponseWriter, *http.Request) {}
	urlDefs := []struct {
		url     string
		handler interface{}
	}{
		{"/", depHandler},
		{"/thing", normalHandler},
	}

	for _, urlDef := range urlDefs {
		handleFunc, err := resolver.HttpHandler(urlDef.handler, errHandler)
		if err != nil {
			panic(err)
		}

		http.HandleFunc(urlDef.url, handleFunc)
	}

	// listen and serve
}
