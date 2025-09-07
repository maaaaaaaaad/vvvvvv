package server

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func NewMux(graphqlHandler http.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(
		"/graphql",
		graphqlHandler,
	)
	mux.HandleFunc(
		"/health",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.Header().Set(
				"Content-Type",
				"application/json",
			)
			bs, _ := json.Marshal(HealthResponse{Status: "ok"})
			w.WriteHeader(http.StatusOK)
			w.Write(bs)
		},
	)
	return mux
}
