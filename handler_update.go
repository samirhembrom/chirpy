package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/samirhembrom/chirpy/internal/auth"
	"github.com/samirhembrom/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdate(w http.ResponseWriter, req *http.Request) {
	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(req.Body)
	params := &parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token", err)
		return
	}

	user_id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't valid token", err)
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create hashed password", err)
		return
	}

	user, err := cfg.db.UpdateUser(context.Background(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_password,
		ID:             user_id,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find user", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
