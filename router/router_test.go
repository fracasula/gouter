package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestInitRoutes(t *testing.T) {
	routes := InitRoutes()

	if len(routes) != 0 {
		t.Errorf("InitRoutes should return an empty map got %d elements instead", len(routes))
	}

	if fmt.Sprintf("%v", routes) != "map[]" {
		t.Errorf("InitRoutes should return a map got %v instead", routes)
	}
}

func TestNewRouter(t *testing.T) {
	routes := InitRoutes()
	emptyHandler := func(w http.ResponseWriter, req *http.Request, vars map[string]string) {}

	routes["^/help$"] = emptyHandler
	routes["^/contact$"] = emptyHandler

	router := New(&routes)

	if l := len(router.Routes); l != 2 {
		t.Errorf("Expected two routes in the router, got %d instead", l)
	}

	for _, r := range router.Routes {
		if rp := r.pattern.String(); rp != "^/help$" && rp != "^/contact$" {
			t.Errorf("Unexpected route pattern, got '%s'", rp)
		}

		if htype := reflect.TypeOf(r.handler).String(); htype != "router.RouteHandlerFunc" {
			t.Errorf("Expected first route handler to be 'router.RouteHandlerFunc', got '%s' instead", htype)
		}
	}
}

func TestNewRouterPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected New to panic with broken regular expressions")
		}
	}()

	routes := InitRoutes()
	routes["a broken regex???"] = func(w http.ResponseWriter, req *http.Request, vars map[string]string) {}

	New(&routes)
}

func TestRouting(t *testing.T) {
	routes := InitRoutes()

	routes["^/help$"] = func(w http.ResponseWriter, req *http.Request, _ map[string]string) {
		fmt.Fprint(w, "Hello from the help page!")
	}

	routes["^/another/page$"] = func(w http.ResponseWriter, req *http.Request, _ map[string]string) {
		fmt.Fprintf(w, "Hello from another page!")
	}

	router := New(&routes)

	// Testing simple page
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/help", strings.NewReader(""))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || w.Body.String() != "Hello from the help page!" {
		t.Errorf("Routing to help page failed, got %d: %v", w.Code, w.Body)
	}

	// Testing another page
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/another/page", strings.NewReader(""))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || w.Body.String() != "Hello from another page!" {
		t.Errorf("Routing to help page failed, got %d: %v", w.Code, w.Body)
	}
}

func TestRoutingWithParams(t *testing.T) {
	routes := InitRoutes()

	routes["^/product$"] = func(w http.ResponseWriter, req *http.Request, _ map[string]string) {
		fmt.Fprint(w, "Hello from the product page!")
	}

	routes["^/product/(?P<pid>[0-9]+)$"] = func(w http.ResponseWriter, req *http.Request, vars map[string]string) {
		fmt.Fprintf(w, "Product ID is %v", vars["pid"])
	}

	router := New(&routes)

	// Testing params
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/product/135", strings.NewReader(""))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || w.Body.String() != "Product ID is 135" {
		t.Errorf("Routing to help page failed, got %d: %v", w.Code, w.Body)
	}
}

func TestRoutingNotFound(t *testing.T) {
	routes := InitRoutes()
	routes["^/page$"] = func(w http.ResponseWriter, req *http.Request, _ map[string]string) {
		fmt.Fprint(w, "Hello from a page!")
	}

	router := New(&routes)

	// Testing not found handler
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/a/not/found/page", strings.NewReader(""))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound || w.Body.String() != "404 Not Found" {
		t.Errorf("Routing to not found handler failed, got %d: %v", w.Code, w.Body)
	}
}

func TestMiddlewares(t *testing.T) {
	routes := InitRoutes()
	routes["^/page$"] = func(w http.ResponseWriter, req *http.Request, _ map[string]string) {
		fmt.Fprint(w, "Hello from a page!")
	}

	router := New(&routes)

	router.AddMiddleware(func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			f(w, req)
		}
	})

	router.AddMiddleware(func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("X-Custom-ID", "112233")
			f(w, req)
		}
	})

	// Testing not found handler
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/page", strings.NewReader(""))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || w.Body.String() != "Hello from a page!" {
		t.Errorf("Routing to page failed, got %d: %v", w.Code, w.Body)
	}

	if h := w.Header().Get("X-Custom-ID"); h != "112233" {
		t.Errorf("Custom middleware failed for 'X-Custom-ID' header, got %v", h)
	}

	if h := w.Header().Get("Content-Type"); h != "application/json" {
		t.Errorf("Custom middleware failed for 'Content-Type' header, got %v", h)
	}
}
