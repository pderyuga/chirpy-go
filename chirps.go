package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pderyuga/chirpy-go/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := database.CreateChirpParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
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

	newParams := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: params.UserID,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), newParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	respondWithJSON(w, http.StatusCreated, chirp)
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

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	respondWithJSON(w, http.StatusOK, chirps)
}
