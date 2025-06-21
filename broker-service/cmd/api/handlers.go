package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Message string `json:"message"`
	Error bool `json:"error"`
	Data interface{} `json:"data,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse {
		Message: "Broker Service is running",
		Data: nil,
		Error: false,
	}

	out,_ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out);
}