package main

import (
	"encoding/json"
	"fmt"
	"internal/database"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ApiConfig struct {
	db             *database.DB
	fileserverHits int
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits += 1
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = 0
}

func (cfg *ApiConfig) HitsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf(`<html>
	<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
	</body>
	</html>`, cfg.fileserverHits))
}

func (cfg *ApiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}
	sort.Slice(chirps, func(i, j int) bool { return chirps[i].ID < chirps[j].ID })

	RespondWithJSON(w, 200, chirps)
}

func (cfg *ApiConfig) GetChirp(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "chirpId"))
	if err != nil {
		RespondWithError(w, 400, "something wrong with the id")
		return
	}
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	for _, c := range chirps {
		if c.ID == id {
			RespondWithJSON(w, 200, c)
			return
		}
	}

	RespondWithError(w, 404, "chirp not found")
}

func (cfg *ApiConfig) PostChirps(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	request := reqBody{}
	err := decoder.Decode(&request)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}
	if len(request.Body) > 140 {
		RespondWithError(w, 400, "Chirp is too long")
		return
	}

	chirp, err := cfg.db.CreateChirp(request.Body)
	if err != nil {
		return
	}

	RespondWithJSON(w, 201, chirp)
}

func (cfg *ApiConfig) PostUsers(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	request := reqBody{}
	err := decoder.Decode(&request)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	user, err := cfg.db.CreateUser(request.Email)
	if err != nil {
		return
	}

	RespondWithJSON(w, 201, user)
}
