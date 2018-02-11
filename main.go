package main

import (
	"fmt"
	"log"
	"net/http"
)

type controller struct{}

func (c controller) home(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "This is the homepage!")
	default:
		log.Printf("404 Not Found: %v", req.URL.Path)
		http.Error(w, "Not Found", 404)
	}
}

func (c controller) help(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "This is the help page!")
}

const httpServerAddr = ":8080"

func main() {
	c := controller{}

	http.HandleFunc("/", c.home)
	http.HandleFunc("/help", c.help)

	log.Printf("Listening on port %v...", httpServerAddr)
	http.ListenAndServe(httpServerAddr, nil)
}
