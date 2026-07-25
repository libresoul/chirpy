package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	type resultVals struct {
		ID         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Cannot decode params, invalid json", err)
		return
	}

	if params.Email == "" {
		respondWithError(w, 400, "email is required", nil)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			respondWithError(w, 400, "Email already exists", nil)
			return
		}
		respondWithError(w, 500, "Failed to create user account", err)
		return
	}

	data := resultVals{
		ID:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	}
	respondWithJSON(w, 201, data)
}
