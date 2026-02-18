package main

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/pderyuga/chirpy-go/internal/auth"
	"github.com/pderyuga/chirpy-go/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type loginResponse struct {
		database.User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	expiresIn := math.Max(float64(params.ExpiresInSeconds), 3600) // 1 hour

	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(expiresIn)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error creating access token", err)
		return
	}

	response := loginResponse{
		User: database.User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: jwtToken,
	}

	respondWithJSON(w, http.StatusOK, response)
}
