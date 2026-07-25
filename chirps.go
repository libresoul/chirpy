package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/libresoul/chirpy/internal/auth"
	"github.com/libresoul/chirpy/internal/database"
)

type Chirp struct {
	ID         uuid.UUID `json:"id"`
	Created_at string    `json:"Created_at"`
	Updated_at string    `json:"Updated_at"`
	Body       string    `json:"body"`
	User_ID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", nil)
		return
	}

	uid, err := auth.ValidateJWT(token, cfg.jwt_secret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	params := parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Failed to decode parameters, invalid json", nil)
		return
	}

	if params.Body == "" {
		respondWithError(w, 400, "valid body is required", nil)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	forbiddenWords := []string{"profane", "fornax", "sharbert", "kerfuffle"}
	words := strings.Fields(params.Body)

	for i, s := range words {
		for _, f := range forbiddenWords {
			if strings.ToLower(s) == f {
				words[i] = "****"
			}
		}
	}
	cleanedBody := strings.Join(words, " ")

	chirpParams := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: uid,
	}
	chirp, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, 500, "Failed to create chirp", err)
		return
	}

	respondWithJSON(w, 201, Chirp{
		ID:         chirp.ID,
		Created_at: chirp.CreatedAt.Format(time.DateTime),
		Updated_at: chirp.UpdatedAt.Format(time.DateTime),
		Body:       chirp.Body,
		User_ID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Failed to retrieve chirps", err)
		return
	}

	resBody := []Chirp{}
	for _, c := range chirps {
		resBody = append(resBody, Chirp{
			ID:         c.ID,
			Created_at: c.CreatedAt.Format(time.DateTime),
			Updated_at: c.UpdatedAt.Format(time.DateTime),
			Body:       c.Body,
			User_ID:    c.UserID,
		})
	}

	respondWithJSON(w, 200, resBody)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")
	if chirpId == "" {
		respondWithError(w, 400, "chirp id is required", nil)
		return
	}

	id, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, 400, "Invalid chirp id", nil)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "Chirp not found", nil)
		return
	}

	resBody := Chirp{
		ID:         chirp.ID,
		Created_at: chirp.CreatedAt.Format(time.DateTime),
		Updated_at: chirp.UpdatedAt.Format(time.DateTime),
		Body:       chirp.Body,
		User_ID:    chirp.UserID,
	}
	respondWithJSON(w, 200, resBody)
}
