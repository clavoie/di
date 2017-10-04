// Package di is a dependency injection framework
/*
di supplies several dependency lifetime caching policies, provides dependency aware http handlers compatible with net/http,
and provides a way to clean up dependencies instantiated during an http request.

di only resolves dependencies which are interfaces, the resolver itself, http.ResponseWriter, and *http.Request.
*/
package di
