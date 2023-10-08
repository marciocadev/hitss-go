package api

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func PublishingDeleteClient(id string) {
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

	idBytes := new(bytes.Buffer)
	json.NewEncoder(idBytes).Encode(id)

	m := amqp.Publishing{
		Body: idBytes.Bytes(),
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

func PublishingInsertNewClient(c Client) {
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
	json.NewEncoder(reqBodyBytes).Encode(c)

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

func PublishingUpdateClient(id string, c Client) {
	conn, ch := getRabbitMWConneectionAndChannel()
	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		os.Getenv("QUEUE_UPDATE"),
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

	c.ID = id
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(c)

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
