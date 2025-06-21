package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse {
		Message: "Broker Service is running",
		Data: nil,
		Error: false,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}