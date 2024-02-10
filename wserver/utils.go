package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func StripProfanity(s string) string {
	profannityList := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	clean := []string{}
	words := strings.Split(s, " ")
	for _, w := range words {
		pr := false
		for _, p := range profannityList {
			if strings.ToLower(w) == p {
				clean = append(clean, "****")
				pr = true
				continue
			}
		}
		if !pr {
			clean = append(clean, w)
		}
	}
	return strings.Join(clean, " ")
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func RespondWithError(w http.ResponseWriter, code int, msg string) error {
	return RespondWithJSON(w, code, map[string]string{"error": msg})
}
