package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func emptyHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	// nothing to do here
}

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

	routes["^/contact$"] = func(w http.ResponseWriter, req *http.Request, vars map[string]string) {
		fmt.Fprint(w, "Hello from the contact page!")
	}

	router := NewRouter(&routes)

	w := httptest.NewRecorder()
	body := strings.NewReader("")
	req, _ := http.NewRequest("GET", "/contact", body)

	router.ServeHTTP(w, req)

	t.Errorf("%v, %v", w.Code, w.Body)
}
