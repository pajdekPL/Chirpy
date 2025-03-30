package main

import (
	"context"
	"net/http"
	"time"

	"github.com/PajdekPL/Chirpy/internal/auth"
)

// @Summary      Refresh access token
// @Description  Refreshes the access token using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  ReturnUserData
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /refresh [post]
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

// @Summary      Revoke refresh token
// @Description  Revokes the current refresh token, logging out the user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      204  "No Content"
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /revoke [post]
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
