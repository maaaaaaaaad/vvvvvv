package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	mux := NewMux(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				w.WriteHeader(http.StatusMethodNotAllowed)
			},
		),
	)
	r := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)
	w := httptest.NewRecorder()
	mux.ServeHTTP(
		w,
		r,
	)
	if w.Code != 200 {
		t.Fatalf(
			"status=%d",
			w.Code,
		)
	}
	if got := w.Body.String(); got == "" {
		t.Fatalf("empty body")
	}
}
