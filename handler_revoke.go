package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/samirhembrom/chirpy/internal/auth"
	"github.com/samirhembrom/chirpy/internal/database"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
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

	err = cfg.db.UpdateRefreshTokenRevokedAt(
		context.Background(),
		database.UpdateRefreshTokenRevokedAtParams{
			UpdatedAt: time.Now(),
			RevokedAt: sql.NullTime{
				Valid: true,
				Time:  time.Now().UTC(),
			},
			Token: refreshToken.Token,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
