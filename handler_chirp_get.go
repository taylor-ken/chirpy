package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	// other imports
)

func (cfg *apiConfig) handlerChirpGet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 400, "Couldn't parse id", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(ctx, id)
	if err == sql.ErrNoRows {
		respondWithError(w, 404, "Couldn't retrieve chirp", err)
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirp", err)
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID.UUID, // This should work
		Body:      dbChirp.Body,
	}

	respondWithJSON(w, 200, chirp)
}
