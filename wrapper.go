package main

import "net/http"

type muxWrapper struct {
	handler http.Handler
}

func (m muxWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.RemoteAddr = r.Header.Get("X-Forwarded-For")
	m.handler.ServeHTTP(w, r)
}
