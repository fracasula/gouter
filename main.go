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
		log.Printf("%v", req.URL.Path)
		http.Error(w, "Not Found", 404)
	}
}

func (c controller) help(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "This is the help page!")
}

func main() {
	c := controller{}

	http.HandleFunc("/", c.home)
	http.HandleFunc("/help", c.help)

	log.Print("Listening...")
	http.ListenAndServe("localhost:8080", nil)
}
