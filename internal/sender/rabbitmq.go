package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"

	"github.com/Sugar-pack/test-task/internal/logging"
	"github.com/Sugar-pack/test-task/internal/model"
)

const JSONType = "application/json"

type Producer interface {
	PublishMessage(ctx context.Context, contentType string, message *model.Message) error
}

type RMQProducer struct {
	Queue            string
	ConnectionString string
}

func NewRMQProducer(queue string, conn string) *RMQProducer {
	return &RMQProducer{
		Queue:            queue,
		ConnectionString: conn,
	}
}

func (x RMQProducer) PublishMessage(ctx context.Context, contentType string, model *model.Message) error {
	logger := logging.FromContext(ctx)
	conn, err := amqp.Dial(x.ConnectionString)
	if err != nil {
		logger.WithError(err).Error("Failed to connect to RabbitMQ")

		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	defer func(conn *amqp.Connection) {
		errClose := conn.Close()
		if errClose != nil {
			logger.WithError(errClose).Error("Failed to close connection to RabbitMQ")
		}
	}(conn)

	channel, err := conn.Channel()
	if err != nil {
		logger.WithError(err).Error("Failed to open a channel")

		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer func(ch *amqp.Channel) {
		errClose := ch.Close()
		if errClose != nil {
			logger.WithError(errClose).Error("Failed to close channel")
		}
	}(channel)

	queueDeclare, err := channel.QueueDeclare(x.Queue, false, false, false, false, nil)
	if err != nil {
		logger.WithError(err).Error("Failed to declare a queue")

		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	body, err := json.Marshal(model)
	if err != nil {
		logger.WithError(err).Error("Failed to marshal message")

		return fmt.Errorf("failed to marshal message: %w", err)
	}
	err = channel.Publish(
		"",
		queueDeclare.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		})

	if err != nil {
		logger.WithError(err).Error("Failed to publish a message")

		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}
