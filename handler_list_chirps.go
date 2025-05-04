package main

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) handlerListChirps(w http.ResponseWriter, req *http.Request) {
	chirpsDb, err := cfg.db.ListChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, chirp := range chirpsDb {
		chirps = append(chirps, Chirp{
			ID:         chirp.ID,
			Created_At: chirp.CreatedAt,
			Updated_At: chirp.UpdatedAt,
			Body:       chirp.Body,
			User_Id:    chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
