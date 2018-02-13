package main

import (
	"log"
	"net/http"
	"regexp"
)

// Router gives you a way to map routes to handlers
type Router struct {
	routes          map[string]func(w http.ResponseWriter, req *http.Request, vars map[string]string)
	notFoundHandler func(w http.ResponseWriter, req *http.Request)
}

func (c Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for route, handler := range c.routes {
		re := regexp.MustCompile(route)
		matches := re.FindStringSubmatch(req.URL.Path)

		if len(matches) > 0 {
			vars := make(map[string]string)

			for i, name := range re.SubexpNames() {
				if name != "" {
					vars[name] = matches[i]
				}
			}

			go handler(w, req, vars)
			return
		}
	}

	c.notFoundHandler(w, req)
}

const httpServerAddr = ":8080"

func main() {
	notFound := func(w http.ResponseWriter, req *http.Request) {
		log.Printf("404 Not Found for path %v", req.URL.Path)
	}

	routes := make(map[string]func(w http.ResponseWriter, req *http.Request, vars map[string]string))

	routes["^/help$"] = func(w http.ResponseWriter, req *http.Request, vars map[string]string) {
		log.Println("Gotcha! You are in the help page =)")
	}

	routes["^/product/(?P<pid>[0-9]+)$"] = func(w http.ResponseWriter, req *http.Request, vars map[string]string) {
		log.Printf("Product ID is %v", vars["pid"])
	}

	routes["^/category/(?P<cid>[0-9]+)$"] = func(w http.ResponseWriter, req *http.Request, vars map[string]string) {
		log.Printf("Category ID is %v", vars["cid"])
	}

	r := Router{routes, notFound}

	log.Printf("Listening on port %v...", httpServerAddr)
	http.ListenAndServe(httpServerAddr, r)
}
