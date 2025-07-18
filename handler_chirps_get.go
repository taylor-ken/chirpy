package main

import (
	"context"
	"net/http"
	// other imports
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	dbChirps, err := cfg.db.GetChirps(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	var chirps []Chirp
	for _, dbChirp := range dbChirps {
		chirp := Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID.UUID, // This should work
			Body:      dbChirp.Body,
		}
		chirps = append(chirps, chirp)
	}

	respondWithJSON(w, 200, chirps)
}
