package api

import (
	"net/http"
	"time"
)

type Route struct {
	Path        string
	HandlerFunc http.HandlerFunc
}

type optionConf struct {
	routes []Route
}

type Option = func(*optionConf)

func WithRoutes(m []Route) Option {
	return func(o *optionConf) {
		o.routes = append(o.routes, m...)
	}
}

func Server(options ...Option) (server *http.Server) {
	mux := http.NewServeMux()

	muxOption := &optionConf{
		routes: make([]Route, 0),
	}

	for _, optFn := range options {
		optFn(muxOption)
	}

	for _, route := range muxOption.routes {
		mux.Handle(route.Path, route.HandlerFunc)
	}

	server = &http.Server{
		Addr:        "localhost:1337",
		Handler:     mux,
		IdleTimeout: 30 * time.Second,
	}
	return
}
