package di

import (
	"net/http"
	"reflect"
)

var requestType = reflect.TypeOf((**http.Request)(nil)).Elem()
var responseWriterType = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
