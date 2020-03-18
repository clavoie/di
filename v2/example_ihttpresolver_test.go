package di_test

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/clavoie/di/v2"
)

func ErrHandler(err *di.ErrResolve, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(os.Stderr, err)
	w.WriteHeader(http.StatusInternalServerError)
}

type SomeDependency interface {
	// optional: only dependencies that perform cleanup need to implement
	// IHttpClosable
	di.IHttpClosable
}
type DependencyImpl struct{ r *http.Request }

func (impl *DependencyImpl) Di_HttpClose() {
	// called after the request has ended
}
func NewDependency(r *http.Request) SomeDependency { return &DependencyImpl{r} }

// optional: only needed if you're interested in the time it takes
// di to resolve dependencies
type PrintLogger struct{}

func NewILogger() di.ILogger { return new(PrintLogger) }
func (pl *PrintLogger) HttpDuration(resolveDuration time.Duration) {
	fmt.Println("time to resolve is: ", resolveDuration)
}

var deps = []*di.Def{
	{NewDependency, di.PerHttpRequest},
	{NewILogger, di.Singleton},
}

func Handler(dep SomeDependency) { /* handle request */ }

var httpDefs = []*di.HttpDef{
	{Handler, "/some/pattern"},
}

func ExampleIHttpResolver() {
	resolver, err := di.NewResolver(ErrHandler, deps)

	if err != nil {
		panic(err)
	}

	err = resolver.SetDefaultServeMux(httpDefs)

	if err != nil {
		panic(err)
	}

	// http.ListenAndServe(":8080", nil)

	//
	// if you'd like to use something other than the DefaultServeMux
	// di can create net/http compatible functions that you can feed
	// to another serveMux or http.Handle
	//
	serveMux := http.NewServeMux()

	for _, httpDef := range httpDefs {
		httpHandler, err := resolver.HttpHandler(httpDef.Handler)

		if err != nil {
			panic(err)
		}

		serveMux.HandleFunc(httpDef.Pattern, httpHandler)
	}

	/*
		server := &http.Server{
			Handler: serveMux,
			// etc...
		}

		server.ListenAndServe()
	*/
}
