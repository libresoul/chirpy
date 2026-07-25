package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaEvent(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Failed to decode parameters, invalid json", nil)
		return
	}

	if params.Data.UserID == uuid.Nil {
		respondWithError(w, 400, "user_id is required", nil)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	_, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, 404, "User not found", nil)
			return
		}
		respondWithError(w, 500, "Failed to upgrade member", err)
		return
	}

	w.WriteHeader(204)
}
