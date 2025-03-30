package main

import (
	"net/http"
	"time"

	"github.com/PajdekPL/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type ReturnUserData struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	UserName    string    `json:"user_name"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

// @Summary      Get user data
// @Description  Retrieves authenticated user's profile information
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  ReturnUserData
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /users [get]
func (cfg *apiConfig) handlerGetUserData(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bad token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bad token", err)
		return
	}

	user, err := cfg.db.GetUserByID(r.Context(), userID)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get user data", err)
		return
	}

	responseData := ReturnUserData{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		UserName:    user.UserName,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, responseData)
}
