package main

import (
	"context"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.db.ListChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	authorID := uuid.Nil
	authorIDString := req.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
	}

	sortString := "asc"
	sortParams := req.URL.Query().Get("sort")
	if sortParams != "" {
		sortString = sortParams
	}

	chirps := []Chirp{}

	for _, dbChirp := range dbChirps {
		if authorID != uuid.Nil && dbChirp.UserID != authorID {
			continue
		}

		chirps = append(chirps, Chirp{
			ID:         dbChirp.ID,
			Created_At: dbChirp.CreatedAt,
			Updated_At: dbChirp.UpdatedAt,
			User_Id:    dbChirp.UserID,
			Body:       dbChirp.Body,
		})
	}

	if sortString == "asc" {
		sort.Slice(
			chirps,
			func(i, j int) bool { return chirps[i].Created_At.Before(chirps[j].Created_At) },
		)
	} else {
		sort.Slice(
			chirps,
			func(i, j int) bool { return chirps[i].Created_At.After(chirps[j].Created_At) },
		)
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
