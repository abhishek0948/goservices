package main

import (
	"net/http"

	"github.com/abhishek0948/goservices/logger-service/data"
)

type JSONPayLoad struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayLoad JSONPayLoad

	if err := app.readJSON(w, r, &requestPayLoad); err != nil {
		app.errorJSON(w, err)
		return
	}

	logEntry := data.LogEntry {
		Name: requestPayLoad.Name,
		Data: requestPayLoad.Data,
	}

	err := app.Models.LogEntry.Insert(logEntry);
	if err!= nil {
		app.errorJSON(w,err);
		return;
	}

	resp := jsonResponse{
		Error: false,
		Message: "Log entry success",
	}

	app.writeJSON(w, http.StatusAccepted, resp, nil)
}