package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/taylor-ken/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "cannot get token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtsecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	if dbChirp.UserID.UUID != userID {
		respondWithError(w, 403, "user is not the author, cannot delete chirp", err)
		return
	}
	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	// For 204 No Content, just set the status - no body needed
	w.WriteHeader(http.StatusNoContent)
}
