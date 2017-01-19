package mware

import (
	"net/http"
	"reflect"
	"testing"
)

func testHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func TestNewMware(t *testing.T) {
	mw := New(testHandler)
	if len(mw.handlers) != 1 {
		t.Errorf("Expected one middleware, got %d", len(mw.handlers))
	}

	mw = New(testHandler, testHandler)
	if len(mw.handlers) != 2 {
		t.Errorf("Expected 2 middleware, got %d", len(mw.handlers))
	}
}

func TestAppend(t *testing.T) {
	mw := New(testHandler)
	mw.Append(testHandler)
	if len(mw.handlers) != 2 {
		t.Errorf("Expected 2 middleware, got %d", len(mw.handlers))
	}

	mw.Append(testHandler)
	if len(mw.handlers) != 3 {
		t.Errorf("Expected 3 middleware, got %d", len(mw.handlers))
	}
}

func TestRun(t *testing.T) {
	mw := New()
	if mw.Run(nil) != http.DefaultServeMux {
		t.Error("Expected http.DefaultServeMux")
	}

	testHttpHandler := testHandler(http.DefaultServeMux)
	if !assertFuncs(mw.Run(testHttpHandler), testHttpHandler) {
		t.Error("Run does not work")
	}
}

func TestRunFunc(t *testing.T) {
	mw := New()
	f1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	if !assertFuncs(mw.RunFunc(f1), f1) {
		t.Error("RunFunc does not work")
	}
}

func assertFuncs(f1, f2 interface{}) bool {
	v1 := reflect.ValueOf(f1)
	v2 := reflect.ValueOf(f2)

	return v1.Pointer() == v2.Pointer()
}
