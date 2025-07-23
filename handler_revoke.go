package main

import (
	"net/http"

	"github.com/taylor-ken/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

	// 1. Extract the bearer token (with error handling)
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "cannot get token", err)
		return
	}

	// 2. Set revoked_at and updated_at to now for the token row in DB (with error handling)
	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "could not revoke token", err)
		return
	}

	// 3. Respond with 204 status, no content
	w.WriteHeader(http.StatusNoContent)
}
