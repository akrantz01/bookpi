package main

import (
	"github.com/akrantz01/bookpi/server/assets"
	"net/http"
)

type notFoundServer struct{}

func (nfs *notFoundServer) Open(string) (http.File, error) {
	return assets.HTTP.Open("build/404.html")
}

func notFoundHandler() http.Handler {
	return http.FileServer(new(notFoundServer))
}
