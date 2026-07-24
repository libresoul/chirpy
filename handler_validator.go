package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coudn't decode parameters", err)
		return
	}

	if params.Body == "" {
		respondWithError(w, http.StatusBadRequest, "body is required", nil)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
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

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: strings.Join(words, " "),
	})
}
