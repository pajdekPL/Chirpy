package main

import (
	"context"
	"net/http"
	"time"

	"github.com/PajdekPL/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	type ReturnUserData struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	userID, err := cfg.db.GetUserIdByRefreshToken(context.Background(), refreshToken)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}
	jwt, err := auth.MakeJWT(userID, cfg.jwtSecret, time.Hour)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token not valid", err)
		return
	}
	respondWithJSON(w, http.StatusOK, ReturnUserData{
		Token: jwt,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	err = cfg.db.RevokeToken(context.Background(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke session", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}
