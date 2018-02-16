package main

import (
	"log"
	"net/http"

	"./http/router"
)

const httpServerAddr = ":8080"

func main() {
	routes := router.InitRoutes()

	routes["^/help$"] = func(w http.ResponseWriter, req *http.Request, vars map[string]string) {
		log.Println("Gotcha! You are in the help page =)")
	}

	routes["^/product/(?P<pid>[0-9]+)$"] = func(w http.ResponseWriter, req *http.Request, vars map[string]string) {
		log.Printf("Product ID is %v", vars["pid"])
	}

	routes["^/category/(?P<cid>[0-9]+)$"] = func(w http.ResponseWriter, req *http.Request, vars map[string]string) {
		log.Printf("Category ID is %v", vars["cid"])
	}

	r := router.NewRouter(&routes)

	log.Printf("Listening on port %v...", httpServerAddr)
	http.ListenAndServe(httpServerAddr, r)
}
