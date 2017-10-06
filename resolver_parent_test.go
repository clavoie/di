package di

import (
	"fmt"
	"net/http"
	"testing"
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

func TestResolverParent(t *testing.T) {
	t.Run("NewResolver", func(t *testing.T) {
		t.Run("InvalidDefs", func(t *testing.T) {
			defs1 := NewDefs()
			defs2 := NewDefs()
			err := defs1.Add(NewA, Singleton)

			if err != nil {
				t.Fatal(err)
			}

			err = defs2.Add(NewA, PerResolve)

			if err != nil {
				t.Fatal(err)
			}

			_, err = NewResolver(Join(defs1, defs2))
			if err == nil {
				t.Fatal("expecting NewResolver error")
			}
		})
	})
	t.Run("HttpHandler", func(t *testing.T) {
		w := (http.ResponseWriter)(new(TestResponseWriter))
		r := new(http.Request)
		defs := NewDefs()

		var closer1 HttpCloser
		closer := &closer1

		err := defs.Add(func() A {
			return closer
		}, PerHttpRequest)

		if err != nil {
			t.Fatal(err)
		}

		resolver, err := NewResolver(defs)
		if err != nil {
			t.Fatal(err)
		}

		errHandler := func(err error, w http.ResponseWriter, r *http.Request) {
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

		handlerFn, err := resolver.HttpHandler(handler, errHandler)
		if err != nil {
			t.Fatal(err)
		}

		handlerFn(w, r)
		if closer1.isClosed == false {
			t.Fatal("dependency not closed")
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
		defs := NewDefs()
		err := defs.Add(func() (A, error) {
			return nil, fmt.Errorf("error")
		}, PerHttpRequest)

		if err != nil {
			t.Fatal(err)
		}

		resolver, err := NewResolver(defs)
		if err != nil {
			t.Fatal(err)
		}

		isErrHandlerCalled := false
		errHandler := func(err error, w http.ResponseWriter, r *http.Request) {
			isErrHandlerCalled = true
		}

		handler := func(a A, innerW http.ResponseWriter, innerR *http.Request) {}

		handlerFn, err := resolver.HttpHandler(handler, errHandler)
		if err != nil {
			t.Fatal(err)
		}

		handlerFn(w, r)
		if isErrHandlerCalled == false {
			t.Fatal("err handler never called")
		}
	})
	t.Run("InvalidHandler", func(t *testing.T) {
		resolver, err := NewResolver(NewDefs())
		if err != nil {
			t.Fatal(err)
		}

		_, err = resolver.HttpHandler("hello", nil)
		if err == nil {
			t.Fatal("expecting http handler to fail with invalid args")
		}
	})
}
