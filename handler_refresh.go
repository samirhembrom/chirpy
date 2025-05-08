package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/samirhembrom/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token ", err)
		return
	}

	refreshToken, err := cfg.db.GetUserFromRefreshToken(context.Background(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token doesn't exist", err)
		return
	}

	fmt.Printf(
		"ExpiresAt::%v\nRevokeAt::%v\nTime::%v\n",
		refreshToken.ExpiresAt,
		refreshToken.RevokedAt.Time,
		time.Now(),
	)

	if time.Now().After(refreshToken.ExpiresAt) ||
		(refreshToken.RevokedAt.Valid && time.Now().After(refreshToken.RevokedAt.Time)) {
		respondWithError(w, http.StatusUnauthorized, "Expired token", err)
		return
	}

	jwt, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: jwt,
	})
}
