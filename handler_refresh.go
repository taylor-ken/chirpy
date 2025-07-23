package main

import (
	"net/http"
	"time"

	"github.com/taylor-ken/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "cannot get token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "token not found, expired, or revoked", err)
		return
	}

	userJWT, err := auth.MakeJWT(user.ID, cfg.jwtsecret, 3600*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": userJWT,
	})
}
