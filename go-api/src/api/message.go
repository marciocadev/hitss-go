package api

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func PublishingDeleteClient(msg string) {
	conn, ch := getRabbitMWConneectionAndChannel()
	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		os.Getenv("QUEUE_DELETE"),
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

func PublishingInsertNewClient(msg Client) {
	conn, ch := getRabbitMWConneectionAndChannel()
	defer conn.Close()
	defer ch.Close()

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

func getRabbitMWConneectionAndChannel() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(os.Getenv("RMQ_URL"))
	if err != nil {
		log.Fatalf("%s: %v", "Failed to connect to RabbitMQ", err)
		panic(err)
	}

	log.Println("Successfully connected to our RabbitMQ instance")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %v", "Failed to open a channel", err)
		panic(err)
	}

	log.Println("Successfully open a channel")

	return conn, ch
}
