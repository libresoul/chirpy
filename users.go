package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/libresoul/chirpy/internal/auth"
	"github.com/libresoul/chirpy/internal/database"
)

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID         uuid.UUID `json:"id"`
	Created_at string    `json:"created_at"`
	Updated_at string    `json:"updated_at"`
	Email      string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	params := credentials{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Cannot decode params, invalid json", err)
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(w, 400, "email and password required", nil)
		return
	}

	hashedPw, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Failed to create account", err)
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPw,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			respondWithError(w, 400, "Email already exists", nil)
			return
		}
		respondWithError(w, 500, "Failed to create user account", err)
		return
	}

	data := User{
		ID:         user.ID,
		Created_at: user.CreatedAt.Format(time.DateTime),
		Updated_at: user.UpdatedAt.Format(time.DateTime),
		Email:      user.Email,
	}
	respondWithJSON(w, 201, data)
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	params := credentials{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Failed to decode params invlid json", nil)
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(w, 400, "email and password required", nil)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password", nil)
		return
	}

	isMatched, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 500, "Failed to look up user info", err)
		return
	}

	if !isMatched {
		respondWithError(w, 401, "Incorrect email or password", nil)
		return
	}

	resBody := User{
		ID:         user.ID,
		Created_at: user.CreatedAt.Format(time.DateTime),
		Updated_at: user.UpdatedAt.Format(time.DateTime),
		Email:      user.Email,
	}
	respondWithJSON(w, 200, resBody)
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("PLATFORM") != "dev" {
		respondWithError(w, 403, "Not available on non-dev environments", nil)
		return
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "Failed to delete user accounts", err)
		return
	}

	w.WriteHeader(200)
}
