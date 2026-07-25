package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/libresoul/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body    string `json:"body"`
		User_ID string `json:"user_id"`
	}
	type returnVals struct {
		ID         uuid.UUID `json:"id"`
		Created_at string    `json:"Created_at"`
		Updated_at string    `json:"Updated_at"`
		Body       string    `json:"body"`
		User_ID    uuid.UUID `json:"user_id"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Failed to decode parameters, invalid json", nil)
		return
	}

	uid, uid_err := uuid.Parse(params.User_ID)
	if params.Body == "" || uid_err != nil {
		respondWithError(w, 400, "valid body and user_id are required", nil)
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

	respondWithJSON(w, 201, returnVals{
		ID:         chirp.ID,
		Created_at: chirp.CreatedAt.Format(time.DateTime),
		Updated_at: chirp.UpdatedAt.Format(time.DateTime),
		Body:       chirp.Body,
		User_ID:    chirp.UserID,
	})
}
