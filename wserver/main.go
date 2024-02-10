package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	id             int
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits += 1
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = 0
}

func (cfg *apiConfig) hitsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf(`<html>
	<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
	</body>
	</html>`, cfg.fileserverHits))
}

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

func getChirps(w http.ResponseWriter, r *http.Request) {
	type resBody struct {
		ID   int    `json:"id"`
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	request := reqBody{}
	err := decoder.Decode(&request)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
	}
	if len(request.Body) > 140 {
		RespondWithError(w, 400, "Chirp is too long")
		return
	}
	RespondWithJSON(w, 200, resBody{Body: request.Body})
}

func postChirps(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Body string `json:"body"`
	}
	type resBody struct {
		ID   int    `json:"id"`
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	request := reqBody{}
	err := decoder.Decode(&request)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
	}
	if len(request.Body) > 140 {
		RespondWithError(w, 400, "Chirp is too long")
		return
	}
	RespondWithJSON(w, 200, resBody{Body: request.Body})
}

func main() {
	config := apiConfig{fileserverHits: 0}

	r := chi.NewRouter()
	api := chi.NewRouter()
	admin := chi.NewRouter()

	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	r.Handle("/app/*", config.middlewareMetricsInc(fsHandler))
	r.Handle("/app", config.middlewareMetricsInc(fsHandler))

	api.Get("/healthz", healthz)
	api.HandleFunc("/reset", config.resetHandler)

	// Chirps CRUD
	api.Get("/chirps", getChirps)
	api.Post("/chirps", postChirps)

	admin.Get("/metrics", config.hitsHandler)

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
