package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/samirhembrom/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, req *http.Request) {
	token_string, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access", err)
		return
	}

	user_id, err := auth.ValidateJWT(token_string, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Unauthorized access", err)
		return
	}

	urlParams := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(urlParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse chirp id", err)
		return
	}

	chirp, err := cfg.db.GetChirp(context.Background(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp", err)
		return
	}

	if chirp.UserID != user_id {
		respondWithError(w, http.StatusForbidden, "Forbidden action", err)
		return
	}

	err = cfg.db.DeleteChirp(context.Background(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(204)
}
