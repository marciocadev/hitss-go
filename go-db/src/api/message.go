package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/streadway/amqp"
)

func GetRabbitMQChannel() (*amqp.Connection, *amqp.Channel) {
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

func ConsumeInsertClients(ch *amqp.Channel) {
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

	msgs, err := ch.Consume(
		q.Name,
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

	stopChan := make(chan bool)
	go func() {
		// connection string
		psqlConn := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			"postgres", 5432, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
		// open database
		db, err := sql.Open("postgres", psqlConn)
		if err != nil {
			log.Fatal(err)
		}
		// close database
		defer db.Close()

		var insertStmt string = "INSERT INTO hitss.cliente (id, cpf, nome)   values ($1, $2, $3)"
		stmt, err := db.Prepare(insertStmt)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("our Consumer ready, PID: %d", os.Getgid())
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			c := Client{}
			err := json.Unmarshal([]byte(d.Body), &c)
			if err != nil {
				fmt.Println("Error in JSON unmarshalling from json marshalled object:", err)
				return
			}

			res, err := stmt.Exec(c.ID, c.CPF, c.Nome)
			fmt.Println("inserting")
			if err != nil || res == nil {
				log.Fatal(err)
			}

			d.Ack(true)
		}
		stmt.Close()
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-stopChan
}
