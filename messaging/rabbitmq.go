package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/njeruthuo/comms-service/emails"
	"github.com/njeruthuo/comms-service/sms"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	passwordResetQueue     = "REDIS_PASSWORD_RESET_QUEUE_NAME"
	emailVerificationQueue = "REDIS_EMAIL_VERIFICATION_QUEUE_NAME"
	phoneVerificationQueue = "REDIS_PHONE_VERIFICATION_QUEUE_NAME"
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

func (r *RabbitMQ) ReadRabbitMQChannel(queueName, name string) {
	if queueName == "" {
		log.Fatalf("%s is not set", name)
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

		var err error
		switch name {
		case passwordResetQueue:
			err = handleEmail(msg.Body, "Password reset request")
		case emailVerificationQueue:
			err = handleEmail(msg.Body, "Verify your email address")
		case phoneVerificationQueue:
			err = handleSMS(msg.Body)
		default:
			log.Printf("no handler for queue %q, discarding message", name)
			msg.Nack(false, false)
			continue
		}

		if err != nil {
			log.Printf("failed to process message from %s: %v", name, err)
			msg.Nack(false, false)
			continue
		}

		msg.Ack(false)
	}
}

// handleEmail decodes an email payload and sends it through the existing
// email service. It is reused for both password reset and email verification.
func handleEmail(body []byte, subject string) error {
	var payload emails.EmailPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("failed to decode email payload: %w", err)
	}
	return emails.SendEmailService(payload.Email, subject, payload.Token)
}

// handleSMS decodes an SMS payload and sends it through the SMS service.
func handleSMS(body []byte) error {
	var payload sms.SMSPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("failed to decode sms payload: %w", err)
	}
	return sms.SendSMSService(payload.Phone, payload.Token)
}
