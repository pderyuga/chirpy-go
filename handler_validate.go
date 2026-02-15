package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type validResponse struct {
		Valid       bool   `json:"valid"`
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, err.Error(), err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleanedBody := getCleanedBody(params.Body, badWords)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	validResp := validResponse{Valid: true, CleanedBody: cleanedBody}

	respondWithJSON(w, http.StatusOK, validResp)
}

func getCleanedBody(chirp string, badWords map[string]struct{}) string {
	placeholder := "****"
	words := strings.Split(chirp, " ")

	for i, word := range words {
		lowercaseWord := strings.ToLower(word)
		if _, ok := badWords[lowercaseWord]; ok {
			words[i] = placeholder
		}
	}
	cleanedBody := strings.Join(words, " ")
	return cleanedBody
}
