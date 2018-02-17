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
}

// HandlerFunc interface
type HandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

// Route stuct
type Route struct {
	pattern *regexp.Regexp
	handler HandlerFunc
}

// NewRouter constructor
func NewRouter(routes *map[string]HandlerFunc) *Router {
	compiledRoutes := make([]Route, len(*routes))

	i := 0
	for route, handler := range *routes {
		// compiling all routes beforehand so if one is broken we get a panic
		// and block the HTTP server from starting in the first place
		re := regexp.MustCompile(route)
		compiledRoutes[i] = Route{re, handler}
		i++
	}

	return &Router{compiledRoutes, notFoundDefaultHandler}
}

// InitRoutes returns an empty map of routes
func InitRoutes() map[string]HandlerFunc {
	return make(map[string]HandlerFunc)
}

func notFoundDefaultHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404 Not Found")
}

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
