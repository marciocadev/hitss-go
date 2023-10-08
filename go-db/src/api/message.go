package api

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func GetRabbitMQConn() *amqp.Connection {
	conn, err := amqp.Dial(os.Getenv("RMQ_URL"))
	if err != nil {
		log.Fatalf("%s: %v", "Failed to connect to RabbitMQ", err)
		panic(err)
	}

	log.Println("Successfully connected to our RabbitMQ instance")

	return conn
}

func GetRabbitMQChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %v", "Failed to open a channel", err)
		panic(err)
	}

	log.Println("Successfully open a channel")

	return ch
}

func GetRabbitMQQueue(ch *amqp.Channel, queueName string) {
	_, err := ch.QueueDeclare(
		queueName,
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
}

func ConsumeUpdateClient(ch *amqp.Channel, queueName string) {
	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s: %v", "Failed to delivery the message", err)
	}

	db := OpenConn()
	defer db.Close()

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		c := Client{}
		err := json.Unmarshal([]byte(d.Body), &c)
		if err != nil {
			log.Fatalf("Error in JSON unmarshalling from json marshalled object: %v", err)
			return
		}

		stmt, params := GetUpdateStatement(db, c)
		defer stmt.Close()

		UpdateClient(stmt, params)

		// if success remove from queue
		d.Ack(true)
	}
}

func ConsumeDeleteClient(ch *amqp.Channel, queueName string) {
	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s: %v", "Failed to delivery the message", err)
	}

	db := OpenConn()
	defer db.Close()

	stmt := GetDeleteStatement(db)
	defer stmt.Close()

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		var id string = ""
		err := json.Unmarshal([]byte(d.Body), &id)
		if err != nil {
			log.Fatalf("Error in JSON unmarshalling from json marshalled object: %v", err)
			return
		}

		DeleteClient(stmt, id)
		// if success remove from queue
		d.Ack(true)
	}
}

func ConsumeInsertClient(ch *amqp.Channel, queueName string) {
	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s: %v", "Failed to delivery the message", err)
	}

	db := OpenConn()
	defer db.Close()

	stmt := GetInsertStatement(db)
	defer stmt.Close()

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		c := Client{}
		err := json.Unmarshal([]byte(d.Body), &c)
		if err != nil {
			log.Fatalf("Error in JSON unmarshalling from json marshalled object: %v", err)
			return
		}

		dt := convertStringToDate(c.DtNascimento)
		InsertClient(stmt, c.ID, c.Nome, c.Sobrenome, c.Contato, c.Endereco, dt, c.CPF)
		// if success remove from queue
		d.Ack(true)
	}
}

func convertStringToDate(s string) time.Time {
	layout := "02/01/2006"
	dt, _ := time.Parse(layout, s)
	return dt
}
