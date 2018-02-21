# gouter

![Travis CI build](https://travis-ci.org/fracasula/gouter.svg?branch=master)

A super simple router for [Go](https://golang.org/) applications.

It uses regular expressions for route matching to give maximum flexibility:

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fracasula/gouter/router"
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

	router.AddMiddleware(func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			f(w, req)
		}
	})

	log.Printf("Listening on port %v...", httpServerAddr)
	http.ListenAndServe(httpServerAddr, r)
}
```
## How to go get it

* In your terminal: `go get -u github.com/fracasula/gouter`
* In your go files: `import "github.com/fracasula/gouter/router"`

## Try it out

Just pull the project, change the `main.go` file as you wish and run `make`.
It will build and run a Docker container with an HTTP server on the 8080 for you to test.
