package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct {
	
}

func main() {
	app := Config{}

	log.Printf("Starting Broker Service on port %s\n", webPort)

	srv:=&http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err:=srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}