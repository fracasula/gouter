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
	Middlewares     []Middleware
}

// Middleware type
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Route stuct
type Route struct {
	pattern *regexp.Regexp
	handler RouteHandlerFunc
}

// RouteHandlerFunc type
type RouteHandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

// New constructor
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

// AddMiddleware function
func (r *Router) AddMiddleware(mw Middleware) {
	r.Middlewares = append(r.Middlewares, mw)
}

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler := func(w http.ResponseWriter, req *http.Request) {
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

	buildMiddlewaresChain(handler, r.Middlewares...)(w, req)
}

func notFoundDefaultHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404 Not Found")
}

func buildMiddlewaresChain(f http.HandlerFunc, mws ...Middleware) http.HandlerFunc {
	if len(mws) == 0 {
		return f
	}

	return mws[0](buildMiddlewaresChain(f, mws[1:cap(mws)]...))
}
