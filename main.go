package main

import (
	"log"
	"net/http"
	"regexp"
)

// Router gives you a way to map routes to handlers
type Router struct {
	routes          map[string]func(w http.ResponseWriter, req *http.Request)
	notFoundHandler func(w http.ResponseWriter, req *http.Request)
}

func (c Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for route, handler := range c.routes {
		match, _ := regexp.MatchString(route, req.URL.Path)

		if match {
			go handler(w, req)
			return
		}
	}

	c.notFoundHandler(w, req)
}

const httpServerAddr = ":8080"

func main() {
	routes := make(map[string]func(w http.ResponseWriter, req *http.Request))
	notFound := func(w http.ResponseWriter, req *http.Request) {
		log.Printf("404 Not Found for path %v", req.URL.Path)
	}

	routes["^/help$"] = func(w http.ResponseWriter, req *http.Request) {
		log.Println("Gotcha! You are in the help page =)")
	}

	r := Router{routes, notFound}

	log.Printf("Listening on port %v...", httpServerAddr)
	http.ListenAndServe(httpServerAddr, r)
}
