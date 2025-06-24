package main

import (
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Connect to RabbitMQ 
	rabbiqmqConn,err := connect();
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		return ;
	}
	defer rabbiqmqConn.Close()

}

func connect() (*amqp.Connection,error) {
	var counts int64
	var backoff = 1*time.Second
	var connection *amqp.Connection

	for {
		c , err := amqp.Dial("amqp://guest:guest@localhost:5672")
		if err != nil {
			log.Println("RabbitMQ is not yet ready");
			counts++
		} else {
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