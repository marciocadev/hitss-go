package main

import "hitss/api"

func main() {
	// api.StartInsertConsumer()

	conn, ch := api.GetRabbitMQChannel()
	api.ConsumeInsertClients(ch)

	defer conn.Close()
	defer ch.Close()
}
