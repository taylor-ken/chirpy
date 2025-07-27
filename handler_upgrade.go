package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/taylor-ken/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpgrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	token, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "cannot get token", err)
		return
	}
	if token != cfg.polka_key {
		respondWithError(w, 401, "API key doesn't match", nil)
		return
	}

	type data struct {
		UserID uuid.UUID `json:"user_id"`
	}

	type requests struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	request := requests{}
	err = decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request", err)
		return
	}

	if request.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeUser(r.Context(), request.Data.UserID)

	if err == sql.ErrNoRows {
		w.WriteHeader(404)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error %s", err)

		return
	}
	w.WriteHeader(204)
}
