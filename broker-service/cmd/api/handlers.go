package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"

	"github.com/abhishek0948/goservices/broker-service/event"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log LogPayLoad `json:"log,omitempty"`
	Mail MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayLoad struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From string `json:"from"`
	To string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// HandleSubmission is the main point of entry into the broker. It accepts a JSON
// payload and performs an action based on the value of "action" in that JSON.
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		// Without RabbitMQ 
		// app.logItem(w,requestPayload.Log)
		// For Rabbit
		// app.logEventViaRabbit(w,requestPayload.Log)
		// For RPC
		app.logItemViaRPC(w,requestPayload.Log)
	case "mail":
		app.sendMail(w,requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

// WithRabbit or RPC
func (app *Config) logItem(w http.ResponseWriter,entry LogPayLoad) {
	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		app.errorJSON(w,err);
		return ;
	}

	logServiceURL := "http://logger-service/log"

	request,err := http.NewRequest("POST",logServiceURL,bytes.NewBuffer(jsonData));
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp,err := client.Do(request);
	if err!= nil {
		app.errorJSON(w,err);
		return 
	}
	defer resp.Body.Close();

	if resp.StatusCode != http.StatusAccepted {
		app.errorJSON(w,err);
		return 
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Log Successful"

	app.writeJSON(w,http.StatusAccepted,payload);
}

// RabbitMQ function
func (app *Config) logEventViaRabbit(w http.ResponseWriter,l LogPayLoad) {
	err := app.pushToQueue(l.Name,l.Data)
	if err != nil  {
		app.errorJSON(w,err);
		return 
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	app.writeJSON(w,http.StatusAccepted,payload);
}

func (app *Config) pushToQueue(name,msg string) error {
	emitter,err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayLoad {
		Name : name,
		Data : msg,
	}

	j,_ := json.MarshalIndent(&payload,"","\t");
	err = emitter.Push(string(j),"log.INFO");
	if err!= nil {
		return err
	}

	return nil
}

type RPCPayload struct {
	Name string
	Data string
}

// With RPC
func (app *Config) logItemViaRPC(w http.ResponseWriter,l LogPayLoad) {
	client,err := rpc.Dial("tcp","logger-service:5001");
	if err != nil {
		app.errorJSON(w,err);
		return ;
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo",rpcPayload,&result);
	if err != nil {
		app.errorJSON(w,err)
		return
	}

	payload := jsonResponse {
		Error: false,
		Message: result,
	}

	app.writeJSON(w,http.StatusAccepted,payload);
}

func (app *Config) sendMail(w http.ResponseWriter,mail MailPayload) {
	jsonData, err := json.MarshalIndent(mail,"","\t");
	if err != nil {
		app.errorJSON(w,err);
		return ;
	}

	mailServiceURL := "http://mail-service/send"

	request,err := http.NewRequest("POST",mailServiceURL,bytes.NewBuffer(jsonData));
	if err != nil {
		app.errorJSON(w,err);
		return ;
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response,err := client.Do(request);
	if err!= nil {
		app.errorJSON(w,err);
		return 
	}
	defer response.Body.Close();

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w,errors.New("error form mail-service"))
		return
	}

	var payload jsonResponse
	payload.Error = false;
	payload.Message = "Message sent to:" + mail.To

	app.writeJSON(w,http.StatusAccepted,payload)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} 

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

