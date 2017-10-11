package di_test

import (
	"fmt"
	"net/http"
	"time"

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

// PrintLogger is an implementation of di.ILogger
type PrintLogger struct{}

func NewILogger() di.ILogger { return new(PrintLogger) }

func (pl *PrintLogger) HttpDuration(resolveDuration time.Duration) {
	fmt.Println("time to resolve is: ", resolveDuration)
}

var deps = []*di.Def{
	di.NewDef(NewHttpDep, di.PerHttpRequest),
	di.NewDef(NewILogger, di.Singleton),
}

func NewHttpResolver() di.IHttpResolver {
	defs := di.NewDefs()
	err := defs.AddAll(deps)

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
