# gouter

![Travis CI build](https://travis-ci.org/fracasula/gouter.svg?branch=master)

A super simple router for [Go](https://golang.org/) applications.
It's still a work in progress, you can can check [here](https://github.com/fracasula/gouter/issues/2)
what needs to be done (feel free to post suggestions too).

It uses regular expressions for route matching to give maximum flexibility:

```go
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

	routes["^/help$"] = func(w http.ResponseWriter, req *http.Request, vars map[string]string) {
		fmt.Fprintf(w, "Gotcha! You are in the help page =)")
	}

	routes["^/product/(?P<pid>[0-9]+)$"] = func(w http.ResponseWriter, req *http.Request, vars map[string]string) {
		fmt.Fprintf(w, "Product ID is %v", vars["pid"])
	}

	routes["^/category/(?P<cid>[0-9]+)$"] = func(w http.ResponseWriter, req *http.Request, vars map[string]string) {
		fmt.Fprintf(w, "Category ID is %v", vars["cid"])
	}

	r := router.NewRouter(&routes)

	log.Printf("Listening on port %v...", httpServerAddr)
	http.ListenAndServe(httpServerAddr, r)
}
```

## Try it out

Just pull the project, change the `main.go` file and run `make`.
It will build a Docker container and expose an HTTP server on the 8080 for you to test.
