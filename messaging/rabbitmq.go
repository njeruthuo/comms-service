package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/njeruthuo/comms-service/emails"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func NewRabbitMQChannel() (*RabbitMQ, error) {
	url := os.Getenv("REDIS_DATABASE_URL")

	if url == "" {
		return nil, fmt.Errorf("REDIS_DATABASE_URL is not set.")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		ch.Close()
		return nil, err
	}

	return &RabbitMQ{
		Conn: conn,
		Ch:   ch,
	}, nil
}

func (r *RabbitMQ) Close() {
	r.Ch.Close()
	r.Conn.Close()
}

func (r *RabbitMQ) String() string {
	return "RabbitMQ connected!"
}

func (r *RabbitMQ) ReadRabbitMQChannel(queueName string) {
	if queueName == "" {
		log.Fatal("REDIS_QUEUE_NAME is not set.")
	}
	queue, err := r.Ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	msgs, err := r.Ch.Consume(
		queue.Name,
		"",    // consumer tag
		false, // auto-ack — false so we ack manually
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	for msg := range msgs {
		log.Printf("received message: %s", msg.Body)

		var emailPayload emails.EmailPayload
		if err := json.Unmarshal(msg.Body, &emailPayload); err != nil {
			log.Printf("failed to decode message: %v", err)
			msg.Nack(false, false)
			continue
		}

		if err := emails.SendEmailService(emailPayload.Email, "Password reset request", emailPayload.Token); err != nil {
			log.Printf("failed to send email: %v", err)
			msg.Nack(false, false)
			continue
		}

		msg.Ack(false)
	}
}
