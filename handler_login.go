package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/samirhembrom/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameter struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(req.Body)
	params := &parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUser(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find user", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access", err)
		return
	}

	expiresInSeconds := time.Hour
	if params.ExpiresInSeconds != nil {
		seconds := time.Duration(expiresInSeconds) * time.Second
		hours := seconds / time.Hour
		if hours < time.Hour {
			expiresInSeconds = hours
		}
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.secret, expiresInSeconds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     jwt,
		},
	})
}
