package main

import (
	"fmt"
	"internal/database"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

func main() {
	db, err := database.NewDB("database.json")
	if err != nil {
		fmt.Println(err)
	}
	config := ApiConfig{fileserverHits: 0, db: db}

	r := chi.NewRouter()
	api := chi.NewRouter()
	admin := chi.NewRouter()

	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	r.Handle("/app/*", config.MiddlewareMetricsInc(fsHandler))
	r.Handle("/app", config.MiddlewareMetricsInc(fsHandler))

	api.Get("/healthz", healthz)
	api.HandleFunc("/reset", config.ResetHandler)

	// Chirps CRUD
	api.Get("/chirps", config.GetChirps)
	api.Get("/chirps/{chirpId}", config.GetChirp)
	api.Post("/chirps", config.PostChirps)

	// Users CRUD
	api.Post("/users", config.PostUsers)

	admin.Get("/metrics", config.HitsHandler)

	r.Mount("/api", api)
	r.Mount("/admin", admin)

	// Adding cors middleware so that requests from boot.dev can work.
	corsMux := middlewareCors(r)

	server := http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}
	server.ListenAndServe()
}
