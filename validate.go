package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type requestBody struct {
	Body string `json:"body"`
}

type errorBody struct {
	Error string `json:"error"`
}

type responseBody struct {
	Valid bool `json:"valid"`
}

func (cfg *apiConfig) handlerValidate(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decorder := json.NewDecoder(req.Body)
	data := requestBody{}
	err := decorder.Decode(&data)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		respBody := errorBody{
			Error: "Something went wrong",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshaing %v", err)
		}
		w.Write(dat)
		return
	}
	body := data.Body
	if len(body) > 140 {
		w.WriteHeader(400)
		respBody := errorBody{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshaing %v", err)
		}
		w.Write(dat)
		return
	}
	w.WriteHeader(200)
	respBody := responseBody{
		Valid: true,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshaing %v", err)
	}
	w.Write(dat)
}
