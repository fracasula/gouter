package router

import (
	"fmt"
	"net/http"
	"regexp"
)

// Router gives you a simple way to map routes to handlers
type Router struct {
	Routes          []Route
	NotFoundHandler http.HandlerFunc
	Middlewares     []*Middleware
}

// Middleware struct
type Middleware struct {
	Handler *MiddlewareHandlerFunc
	Next    *Middleware
}

// MiddlewareHandlerFunc type
type MiddlewareHandlerFunc func(http.ResponseWriter, *http.Request, *Middleware)

// Route stuct
type Route struct {
	pattern *regexp.Regexp
	handler RouteHandlerFunc
}

// RouteHandlerFunc type
type RouteHandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

// NewRouter constructor
func New(routes *map[string]RouteHandlerFunc) *Router {
	compiledRoutes := make([]Route, len(*routes))

	i := 0
	for route, handler := range *routes {
		// compiling all routes beforehand so if one is broken we get a panic
		// and block the HTTP server from starting in the first place
		re := regexp.MustCompile(route)
		compiledRoutes[i] = Route{re, handler}
		i++
	}

	return &Router{compiledRoutes, notFoundDefaultHandler, nil}
}

// InitRoutes returns an empty map of routes
func InitRoutes() map[string]RouteHandlerFunc {
	return make(map[string]RouteHandlerFunc)
}

func notFoundDefaultHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404 Not Found")
}

func (r *Router) AddMiddleware(mw *MiddlewareHandlerFunc) {
	r.Middlewares = append(r.Middlewares, mw)

	if l := len(r.Middlewares); l > 1 {
		r.Middlewares[l-2].Next = r.Middlewares[l-1]
	}
}

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if len(r.Middlewares) > 1 {
		r.Middlewares[0].Handler()
	}

	for _, route := range r.Routes {
		matches := route.pattern.FindStringSubmatch(req.URL.Path)

		if len(matches) > 0 {
			vars := make(map[string]string)

			for i, name := range route.pattern.SubexpNames() {
				if name != "" { // only named groups are allowed
					vars[name] = matches[i]
				}
			}

			route.handler(w, req, vars)
			return
		}
	}

	r.NotFoundHandler(w, req)
}
