package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	type chirpRequest struct {
		Body string `json:"body"`
	}

	type successResponse struct {
		Body string `json:"cleaned_body"`
	}

	type errorResponse struct {
		Body string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := chirpRequest{}
	err := decoder.Decode(&chirp)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(chirp.Body) > 140 {
		errresp := errorResponse{Body: "Chirp is too long"} // What message goes here?
		resp, err := json.Marshal(errresp)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(400)
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		w.Write(resp)
	} else {
		cleaned_body := removeProfanity(chirp.Body)
		sucresp := successResponse{Body: cleaned_body}
		successResp, err := json.Marshal(sucresp)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		w.Write(successResp)
	}

}

func removeProfanity(s string) string {
	profanities := []string{"kerfuffle", "sharbert", "fornax"}
	out := []string{}
	var found bool
	words := strings.Split(s, " ")
	for _, word := range words {
		found = false
		for _, profanity := range profanities {
			if strings.ToLower(word) == profanity {
				found = true
				break
			}
		}
		if found {
			out = append(out, "****")
		} else {
			out = append(out, word)
		}
	}
	out_string := strings.Join(out, " ")
	return strings.TrimSpace(out_string)
}
