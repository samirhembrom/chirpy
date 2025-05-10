package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, req *http.Request) {
	s := req.URL.Query().Get("author_id")
	if s != "" {
		user_id, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't parse auther id", err)
			return
		}
		chirpsDb, err := cfg.db.ListChirpsByUser(context.Background(), user_id)
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
		return
	}
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

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, req *http.Request) {
	urlParams := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(urlParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse user_id", err)
		return
	}

	chirp, err := cfg.db.GetChirp(context.Background(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:         chirp.ID,
		Created_At: chirp.CreatedAt,
		Updated_At: chirp.UpdatedAt,
		Body:       chirp.Body,
		User_Id:    chirp.UserID,
	})
}
