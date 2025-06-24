package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// Connect to RabbitMQ 
	rabbitConn,err := connect();
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		return ;
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection,error) {
	var counts int64
	var backoff = 1*time.Second
	var connection *amqp.Connection

	for {
		c , err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Println("RabbitMQ is not yet ready");
			counts++
		} else {
			log.Println("Connection to rabbitMQ established")
			connection = c;
			break;
		}

		if counts > 5 {
			log.Println(err)
			return nil,err
		}

		backoff = time.Duration(math.Pow(float64(counts),2)) * time.Second
		log.Println("Backing of for", backoff)
		time.Sleep(backoff);
		continue;
	}

	return connection,nil
}