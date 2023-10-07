package main

import (
	"hitss/api"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	conn := api.GetRabbitMQConn()
	defer conn.Close()

	insertChannel := api.GetRabbitMQChannel(conn)
	defer insertChannel.Close()

	deleteChannel := api.GetRabbitMQChannel(conn)
	defer deleteChannel.Close()

	queueInsertName := os.Getenv("QUEUE_INSERT")
	api.GetRabbitMQQueue(insertChannel, queueInsertName)

	queueDeleteName := os.Getenv("QUEUE_DELETE")
	api.GetRabbitMQQueue(deleteChannel, queueDeleteName)

	go api.ConsumeInsertClient(insertChannel, queueInsertName)
	go api.ConsumeDeleteClient(insertChannel, queueDeleteName)

	// Wait for termination signal (Ctrl+C)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	// ch := api.GetRabbitMQChannel()
	// api.ConsumeDeleteClient(ch)
	// api.ConsumeInsertClient(ch)

}
