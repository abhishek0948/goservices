package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer,error) {
	consumer := Consumer {
		conn: conn,
	}

	err := consumer.setup();
	if err != nil {
		return Consumer{},err
	}

	return consumer,nil
}

func (consumer *Consumer) setup() error {
	ch, err := consumer.conn.Channel();
	if err != nil {
		return err
	}

	return declareExchange(ch);
}

type Payload struct{
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch,err := consumer.conn.Channel();
	if err != nil {
		return err
	}
	defer ch.Close()

	q,err := declareRandomQueue(ch);
	if err != nil {
		return err
	}

	for _ , s:= range topics {
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages,err := ch.Consume(q.Name,"",true,false,false,false,nil);
	if err != nil {
		return err
	}

	forever := make(chan bool);
	go func(){
		for d:= range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body,&payload);

			go handlePayload(payload);
		}
	}()

	fmt.Printf("Waiting for msg [Exchange,Queue] [log_topics,%s]",q.Name);
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log","event":
		err:= logEvent(payload);
		if err!=nil {
			log.Println(err)
		}
	case "authenticate":
	
	default:
		err:= logEvent(payload);
		if err!=nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		return err ;
	}

	logServiceURL := "http://logger-service/log"

	request,err := http.NewRequest("POST",logServiceURL,bytes.NewBuffer(jsonData));
	if err != nil {
		return err;
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp,err := client.Do(request);
	if err!= nil {
		return err;
	}
	defer resp.Body.Close();

	if resp.StatusCode != http.StatusAccepted {
		return err;
	}

	return nil
}