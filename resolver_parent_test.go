package di

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

type TestResponseWriter struct{}

func (trw *TestResponseWriter) Header() http.Header       { return nil }
func (trw *TestResponseWriter) Write([]byte) (int, error) { return 0, nil }
func (trw *TestResponseWriter) WriteHeader(int)           {}

type HttpCloser struct {
	isClosed bool
}

func (hc *HttpCloser) A() int        { return 1 }
func (hc *HttpCloser) Di_HttpClose() { hc.isClosed = true }

type Logger struct {
	isCalled bool
}

func (l *Logger) HttpDuration(time.Duration) { l.isCalled = true }

func resolverParentErr(er *ErrResolve, w http.ResponseWriter, r *http.Request) { panic(er) }

func TestResolverParent(t *testing.T) {
	t.Run("NewResolver", func(t *testing.T) {
		t.Run("InvalidDefs", func(t *testing.T) {
			_, err := NewResolver(resolverParentErr, []*Def{
				&Def{NewA, Singleton}, &Def{NewA, PerResolve},
			})
			if err == nil {
				t.Fatal("expecting NewResolver error")
			}
		})
	})
	t.Run("HttpHandler", func(t *testing.T) {
		w := (http.ResponseWriter)(new(TestResponseWriter))
		r := new(http.Request)

		var closer1 HttpCloser
		closer := &closer1
		logger := new(Logger)
		errHandler := func(err *ErrResolve, w http.ResponseWriter, r *http.Request) {
			t.Fatal(err)
		}

		resolver, err := NewResolver(errHandler, []*Def{
			&Def{func() A { return closer }, PerHttpRequest},
			&Def{func() ILogger { return logger }, PerHttpRequest},
		})

		if err != nil {
			t.Fatal(err)
		}

		handler := func(a A, innerW http.ResponseWriter, innerR *http.Request) {
			if a == nil {
				t.Fatal(a)
			}

			if w != innerW {
				t.Fatal(w, innerW)
			}

			if r != innerR {
				t.Fatal(r, innerR)
			}
		}

		handlerFn, err := resolver.HttpHandler(handler)
		if err != nil {
			t.Fatal(err)
		}

		handlerFn(w, r)
		if closer1.isClosed == false {
			t.Fatal("dependency not closed")
		}

		if logger.isCalled == false {
			t.Fatal("logger not called")
		}

		var closer2 HttpCloser
		closer = &closer2
		handlerFn(w, r)
		if closer2.isClosed == false {
			t.Fatal("per http request not closed")
		}
	})
	t.Run("ErrFn", func(t *testing.T) {
		w := (http.ResponseWriter)(new(TestResponseWriter))
		r := new(http.Request)

		isErrHandlerCalled := false
		errHandler := func(err *ErrResolve, w http.ResponseWriter, r *http.Request) {
			isErrHandlerCalled = true
		}

		resolver, err := NewResolver(errHandler, []*Def{
			&Def{func() (A, error) { return nil, fmt.Errorf("error") }, PerHttpRequest},
		})
		if err != nil {
			t.Fatal(err)
		}

		handler := func(a A, innerW http.ResponseWriter, innerR *http.Request) {}

		handlerFn, err := resolver.HttpHandler(handler)
		if err != nil {
			t.Fatal(err)
		}

		handlerFn(w, r)
		if isErrHandlerCalled == false {
			t.Fatal("err handler never called")
		}
	})
	t.Run("InvalidHandler", func(t *testing.T) {
		resolver, err := NewResolver(resolverParentErr, []*Def{})
		if err != nil {
			t.Fatal(err)
		}

		_, err = resolver.HttpHandler("hello")
		if err == nil {
			t.Fatal("expecting http handler to fail with invalid args")
		}
	})
}
