package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PajdekPL/Chirpy/internal/auth"
	"github.com/PajdekPL/Chirpy/internal/database"
)

type UserDataChange struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RespondDataUserChange struct {
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

// @Summary      Update user data
// @Description  Updates the authenticated user's email, name, and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        userData  body      UserDataChange  true  "User data to update"
// @Success      200       {object}  RespondDataUserChange
// @Failure      400  	   {object}  ErrorResponse
// @Failure      404  	   {object}  ErrorResponse
// @Failure      500       {object}  ErrorResponse
// @Router       /users [put]
func (cfg *apiConfig) handlerChangeUserData(w http.ResponseWriter, r *http.Request) {
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
	data := UserDataChange{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Problem decoding data", err)
		return
	}

	if data.Email == "" || data.Password == "" || data.Name == "" {
		respondWithError(w, http.StatusUnauthorized, "Required not empty email, password and name", err)
		return
	}
	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "New password hashing problem", err)
		return
	}
	userDb, err := cfg.db.ChangeUserData(context.Background(), database.ChangeUserDataParams{
		ID:             userID,
		Email:          data.Email,
		UserName:       data.Name,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "", err)
		return
	}
	respondWithJSON(w, http.StatusOK, RespondDataUserChange{
		Email:     userDb.Email,
		UpdatedAt: userDb.UpdatedAt,
		CreatedAt: userDb.CreatedAt,
	})
}

type UserDataChangeName struct {
	UserName string `json:"user_name"`
}
type RespondDataChangeName struct {
	UserName  string    `json:"user_name"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

// @Summary      Change user name data
// @Description  Updates the authenticated user's email, password, and name
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        userData  body      UserDataChangeName  true  "User data to update"
// @Success      200       {object}  RespondDataChangeName
// @Failure      400       {object}  map[string]string
// @Failure      401       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /users [patch]
func (cfg *apiConfig) handlerChangeUserName(w http.ResponseWriter, r *http.Request) {

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
	data := UserDataChangeName{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Problem decoding data", err)
		return
	}

	if data.UserName == "" {
		respondWithError(w, http.StatusUnauthorized, "Required a new user name", err)
		return
	}

	userDb, err := cfg.db.ChangeUserName(context.Background(), database.ChangeUserNameParams{
		ID:       userID,
		UserName: data.UserName,
	})

	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("User with the given name: %s already exists", data.UserName), err)
			return
		}
		respondWithError(w, http.StatusBadRequest, "Problem creating user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, RespondDataChangeName{
		UserName:  userDb.UserName,
		UpdatedAt: userDb.UpdatedAt,
	})
}
