package main

import (
	"fmt"
	"log"
	"net/http"

	"./router"
)

const httpServerAddr = ":8080"

func main() {
	routes := router.InitRoutes()

	routes["^/help$"] = func(w http.ResponseWriter, req *http.Request, _ map[string]string) {
		fmt.Fprintf(w, "Gotcha! You are in the help page =)")
	}

	routes["^/product/(?P<pid>[0-9]+)$"] = func(w http.ResponseWriter, req *http.Request, vars map[string]string) {
		fmt.Fprintf(w, "Product ID is %v", vars["pid"])
	}

	routes["^/category/(?P<cid>[0-9]+)$"] = func(w http.ResponseWriter, req *http.Request, vars map[string]string) {
		fmt.Fprintf(w, "Category ID is %v", vars["cid"])
	}

	r := router.New(&routes)

	log.Printf("Listening on port %v...", httpServerAddr)
	http.ListenAndServe(httpServerAddr, r)
}
