package api

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func StartPublishing(msg Client) {
	conn, err := amqp.Dial(os.Getenv("RMQ_URL"))
	if err != nil {
		log.Fatalf("%s: %v", "Failed to connect to RabbitMQ", err)
		panic(err)
	}
	defer conn.Close()

	log.Println("Successfully connected to our RabbitMQ instance")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %v", "Failed to open a channel", err)
		panic(err)
	}
	defer ch.Close()

	log.Println("Successfully open a channel")

	q, err := ch.QueueDeclare(
		os.Getenv("QUEUE_INSERT"),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s: %v", "Failed to declare a queue", err)
		panic(err)
	}

	log.Println("Successfully declare a queue")

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(msg)

	m := amqp.Publishing{
		Body: reqBodyBytes.Bytes(),
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		m,
	)
	if err != nil {
		log.Fatalf("%s: %v", "Failed to delivery the message", err)
	}

	log.Println("Message delivery to the queue")
}
