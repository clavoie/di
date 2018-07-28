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
	&di.Def{NewHttpDep, di.PerHttpRequest},
	&di.Def{NewILogger, di.Singleton},
}

func DepHandler(dep HttpDep)                         {}
func HttpHandler(http.ResponseWriter, *http.Request) {}
func WriteToLog(resolveErr *di.ErrResolve)           { /* etc */ }

func ErrHandler(err *di.ErrResolve, w http.ResponseWriter, r *http.Request) {
	WriteToLog(err)
	w.WriteHeader(http.StatusInternalServerError)
}

var httpDefs = []*di.HttpDef{
	&di.HttpDef{DepHandler, "/"},
	&di.HttpDef{HttpHandler, "/some/pattern"},
}

func ExampleIHttpResolver_httpHandler() {
	resolver, err := di.NewResolver(ErrHandler, deps)

	if err != nil {
		panic(err)
	}

	err = resolver.SetDefaultServeMux(httpDefs)

	if err != nil {
		panic(err)
	}

	// listen and serve
}
