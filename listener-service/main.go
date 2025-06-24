package main

import (
	"log"
	"math"
	"time"

	"github.com/abhishek0948/goservices/listener-service/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Connect to RabbitMQ 
	rabbitConn,err := connect();
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		return ;
	}
	defer rabbitConn.Close()

	// Start listening to messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create consumer
	consumer , err := event.NewConsumer(rabbitConn)
	if err!= nil {
		log.Fatal("Error creating consumer");
		return 
	}

	// Watch the queue and consume events
	err = consumer.Listen([] string{"log.INFO","log.WARNING","log.ERROR"})
	if err != nil {
		log.Println(err)
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