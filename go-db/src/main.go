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

	updateChannel := api.GetRabbitMQChannel(conn)
	defer updateChannel.Close()

	queueInsertName := os.Getenv("QUEUE_INSERT")
	api.GetRabbitMQQueue(insertChannel, queueInsertName)

	queueDeleteName := os.Getenv("QUEUE_DELETE")
	api.GetRabbitMQQueue(deleteChannel, queueDeleteName)

	queueUpdateName := os.Getenv("QUEUE_UPDATE")
	api.GetRabbitMQQueue(updateChannel, queueUpdateName)

	go api.ConsumeInsertClient(insertChannel, queueInsertName)
	go api.ConsumeDeleteClient(insertChannel, queueDeleteName)
	go api.ConsumeUpdateClient(updateChannel, queueUpdateName)

	// Wait for termination signal (Ctrl+C)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
