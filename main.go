package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

	go r.ReadRabbitMQChannel(os.Getenv("REDIS_PASSWORD_RESET_QUEUE_NAME"), "REDIS_PASSWORD_RESET_QUEUE_NAME")
	go r.ReadRabbitMQChannel(os.Getenv("REDIS_EMAIL_VERIFICATION_QUEUE_NAME"), "REDIS_EMAIL_VERIFICATION_QUEUE_NAME")
	go r.ReadRabbitMQChannel(os.Getenv("REDIS_PHONE_VERIFICATION_QUEUE_NAME"), "REDIS_PHONE_VERIFICATION_QUEUE_NAME")

	// Block until interrupted so the consumer goroutines keep running.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down comms-service")
}
