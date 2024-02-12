package main

import (
	"encoding/json"
	"fmt"
	"internal/database"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type ApiConfig struct {
	secret         string
	pkey           string
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
	a := r.URL.Query().Get("author_id")
	s := r.URL.Query().Get("sort")
	aid, _ := strconv.Atoi(a)

	chirps, err := cfg.db.GetChirps(aid)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}
	if len(s) > 0 {
		if s == "asc" {
			sort.Slice(chirps, func(i, j int) bool { return chirps[i].ID < chirps[j].ID })
			RespondWithJSON(w, 200, chirps)
			return
		} else {
			sort.Slice(chirps, func(i, j int) bool { return chirps[i].ID > chirps[j].ID })
			RespondWithJSON(w, 200, chirps)
			return
		}
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
	chirps, err := cfg.db.GetChirps(0)
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

	// checking for token
	bearer := r.Header.Get("Authorization")
	if len(bearer) == 0 {
		RespondWithError(w, 401, "unauthorized")
		return
	}
	value := strings.Split(bearer, " ")[1]
	token, err := jwt.ParseWithClaims(value, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.secret), nil
	})
	if err != nil {
		RespondWithError(w, 401, "unauthorized")
		return
	}
	userId, _ := token.Claims.GetSubject()
	issued, err := token.Claims.GetIssuer()
	if err != nil || issued == "chirpy-refresh" {
		RespondWithError(w, 401, "unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	request := reqBody{}
	err = decoder.Decode(&request)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}
	if len(request.Body) > 140 {
		RespondWithError(w, 400, "Chirp is too long")
		return
	}

	id, _ := strconv.Atoi(userId)
	chirp, err := cfg.db.CreateChirp(request.Body, id)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	RespondWithJSON(w, 201, chirp)
}

func (cfg *ApiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpId, err := strconv.Atoi(chi.URLParam(r, "chirpId"))
	if err != nil {
		RespondWithError(w, 400, "something wrong with the id")
		return
	}
	chirps, err := cfg.db.GetChirps(0)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	// checking for token
	bearer := r.Header.Get("Authorization")
	if len(bearer) == 0 {
		RespondWithError(w, 401, "unauthorized")
		return
	}
	value := strings.Split(bearer, " ")[1]
	token, err := jwt.ParseWithClaims(value, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.secret), nil
	})
	if err != nil {
		RespondWithError(w, 401, "unauthorized")
		return
	}
	userId, _ := token.Claims.GetSubject()
	issued, err := token.Claims.GetIssuer()
	if err != nil || issued == "chirpy-refresh" {
		RespondWithError(w, 401, "unauthorized")
		return
	}

	id, _ := strconv.Atoi(userId)
	for _, c := range chirps {
		if c.ID == chirpId {
			if c.AuthorID == id {
				cfg.db.DeleteChirp(chirpId)
				RespondWithJSON(w, 200, "")
				return
			} else {
				RespondWithError(w, 403, "unauthorized")
			}
		}
	}

	RespondWithError(w, 404, "chirp not found")
}

func (cfg *ApiConfig) PostUsers(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	request := reqBody{}
	err := decoder.Decode(&request)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	user, err := cfg.db.CreateUser(request.Email, hash)
	if err != nil {
		RespondWithError(w, 400, "user already exists")
		return
	}
	userToSend := struct {
		ID     int    `json:"id"`
		Email  string `json:"email"`
		Chirpy bool   `json:"is_chirpy_red"`
	}{
		ID:     user.ID,
		Email:  user.Email,
		Chirpy: user.IsRed,
	}

	RespondWithJSON(w, 201, userToSend)
}

func (cfg *ApiConfig) PutUsers(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// checking for token
	bearer := r.Header.Get("Authorization")
	if len(bearer) == 0 {
		RespondWithError(w, 401, "unauthorized")
		return
	}
	value := strings.Split(bearer, " ")[1]
	token, err := jwt.ParseWithClaims(value, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.secret), nil
	})
	if err != nil {
		RespondWithError(w, 401, "unauthorized")
		return
	}
	userId, _ := token.Claims.GetSubject()
	issued, err := token.Claims.GetIssuer()
	if err != nil || issued == "chirpy-refresh" {
		RespondWithError(w, 401, "unauthorized")
		return
	}

	// decoding the request
	decoder := json.NewDecoder(r.Body)
	request := reqBody{}
	err = decoder.Decode(&request)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}
	user, err := cfg.db.UpdateUser(id, request.Email, hash)
	if err != nil {
		return
	}
	userToSend := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		ID:    user.ID,
		Email: user.Email,
	}

	RespondWithJSON(w, 200, userToSend)
}

func (cfg *ApiConfig) Polka(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		} `json:"data"`
	}

	// checking for token
	bearer := r.Header.Get("Authorization")
	if len(bearer) == 0 {
		RespondWithError(w, 401, "unauthorized")
		return
	}
	value := strings.Split(bearer, " ")[1]
	if value != cfg.pkey {
		RespondWithError(w, 401, "unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	request := reqBody{}
	err := decoder.Decode(&request)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	if request.Event != "user.upgraded" {
		RespondWithJSON(w, 200, "")
		return
	}

	err = cfg.db.UpdateUserMembership(request.Data.UserId)
	if err != nil {
		RespondWithError(w, 404, "")
		return
	}

	RespondWithJSON(w, 200, "")
}

func (cfg *ApiConfig) PostLogin(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	request := reqBody{}
	err := decoder.Decode(&request)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}
	user, err := cfg.db.GetUserByEmail(request.Email)
	if err != nil {
		RespondWithError(w, 404, "Something went wrong")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		RespondWithError(w, 401, "Unauthorized")
		return
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour).UTC()),
		Subject:   fmt.Sprintf("%v", user.ID),
	})
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-refresh",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 34 * 60)),
		Subject:   fmt.Sprintf("%v", user.ID),
	})
	signedAcc, _ := access.SignedString([]byte(cfg.secret))
	signedRef, _ := refresh.SignedString([]byte(cfg.secret))

	userToSend := struct {
		ID           int    `json:"id"`
		Email        string `json:"email"`
		AccessToken  string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		Chirpy       bool   `json:"is_chirpy_red"`
	}{
		ID:           user.ID,
		Email:        user.Email,
		AccessToken:  signedAcc,
		RefreshToken: signedRef,
		Chirpy:       user.IsRed,
	}

	RespondWithJSON(w, 200, userToSend)
}

func (cfg *ApiConfig) PostRefresh(w http.ResponseWriter, r *http.Request) {
	// checking for token
	bearer := r.Header.Get("Authorization")
	if len(bearer) == 0 {
		RespondWithError(w, 401, "unauthorized")
		return
	}
	value := strings.Split(bearer, " ")[1]
	token, err := jwt.ParseWithClaims(value, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.secret), nil
	})
	if err != nil {
		RespondWithError(w, 401, "unauthorized")
		return
	}
	userId, _ := token.Claims.GetSubject()
	issued, err := token.Claims.GetIssuer()
	if err != nil || issued != "chirpy-refresh" || cfg.db.IsRevoked(value) {
		RespondWithError(w, 401, "unauthorized")
		return
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour).UTC()),
		Subject:   fmt.Sprintf("%v", userId),
	})
	signedAcc, _ := access.SignedString([]byte(cfg.secret))

	tokenToSend := struct {
		AccessToken string `json:"token"`
	}{
		AccessToken: signedAcc,
	}

	RespondWithJSON(w, 200, tokenToSend)
}

func (cfg *ApiConfig) PostRevoke(w http.ResponseWriter, r *http.Request) {
	// checking for token
	bearer := r.Header.Get("Authorization")
	if len(bearer) == 0 {
		RespondWithError(w, 401, "unauthorized")
		return
	}
	value := strings.Split(bearer, " ")[1]
	_, err := jwt.ParseWithClaims(value, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.secret), nil
	})
	if err != nil {
		RespondWithError(w, 401, "unauthorized")
		return
	}

	err = cfg.db.RevokeToken(value)
	if err != nil {
		RespondWithError(w, 401, "unauthorized")
		return
	}

	RespondWithJSON(w, 200, "")
}
