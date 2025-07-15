package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	type chirpRequest struct {
		Body string `json:"body"`
	}

	type successResponse struct {
		Body bool `json:"valid"`
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
		sucresp := successResponse{Body: true}    // Fix: use true, not "success"
		successResp, err := json.Marshal(sucresp) // Fix: use a different variable name
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500) // Fix: should be 500 for server error, not 400
			return
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		w.Write(successResp) // Fix: use the correct variable
	}

}
