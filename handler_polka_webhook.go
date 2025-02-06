package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/PajdekPL/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing ApiKey", err)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "wrong ApiKey", err)
		return
	}
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, "")
		return
	}

	userID, err := uuid.Parse(params.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "", err)
		return
	}

	_, err = cfg.db.UpgradeUserToRed(context.Background(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, parameters{})
}
