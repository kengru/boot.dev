package main

import "strings"

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
