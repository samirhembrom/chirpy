package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/samirhembrom/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, req *http.Request) {
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Header malformed", err)
		return
	}
	fmt.Printf("APIKEY::%v\nENV::%v\n", apiKey, cfg.apiKey)
	if apiKey != cfg.apiKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access", err)
		return
	}
	type parameter struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	params := &parameter{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpdateUserRed(context.Background(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find uuid", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update uuid", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
