package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/njeruthuo/comms-service/messaging"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	r, err := messaging.NewRabbitMQChannel()
	if err != nil {
		log.Fatal("There was a problem connecting to the channel: " + err.Error())
	} else {
		log.Println(r.String())
	}
	defer r.Close()

	r.ReadRabbitMQChannel(os.Getenv("REDIS_QUEUE_NAME"))
}
