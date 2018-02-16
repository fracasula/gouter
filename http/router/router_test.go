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

	routes["^/help$"] = emptyHandler
	routes["^/contact$"] = emptyHandler

	router := NewRouter(&routes)

	if l := len(router.Routes); l != 2 {
		t.Errorf("Expected two routes in the router, got %d instead", l)
	}

	for _, r := range router.Routes {
		if rp := r.pattern.String(); rp != "^/help$" && rp != "^/contact$" {
			t.Error("Unexpected route pattern, got '%s'", rp)
		}

		if htype := reflect.TypeOf(r.handler).String(); htype != "router.HandlerFunc" {
			t.Error("Expected first route handler to be 'router.HandlerFunc', got '%s' instead", htype)
		}
	}
}

func TestNewRouterPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected NewRouter to panic with broken regular expressions")
		}
	}()

	routes := InitRoutes()
	routes["a broken regex???"] = emptyHandler

	NewRouter(&routes)
}

func TestRouting(t *testing.T) {
	routes := InitRoutes()

	routes["^/help$"] = func(w http.ResponseWriter, req *http.Request, vars map[string]string) {
		fmt.Fprint(w, "Hello from the help page!")
	}

	routes["^/product/(?P<pid>[0-9]+)$"] = func(w http.ResponseWriter, req *http.Request, vars map[string]string) {
		fmt.Fprintf(w, "Product ID is %v", vars["pid"])
	}

	router := NewRouter(&routes)

	// Testing simple page
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/help", strings.NewReader(""))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || w.Body.String() != "Hello from the help page!" {
		t.Errorf("Routing to help page failed, got %d: %v", w.Code, w.Body)
	}

	// Testing params
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/product/135", strings.NewReader(""))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || w.Body.String() != "Product ID is 135" {
		t.Errorf("Routing to help page failed, got %d: %v", w.Code, w.Body)
	}

	// Testing not found handler
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/a/not/found/page", strings.NewReader(""))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound || w.Body.String() != "404 Not Found" {
		t.Errorf("Routing to help page failed, got %d: %v", w.Code, w.Body)
	}
}

// functions

func emptyHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	// nothing to do here
}
