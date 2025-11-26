package pubsub

import (
	"context"
	"encoding/json"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota
	Transient
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return errors.New("can't marshal json")
	}
	publishingError := ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        bytes,
	})
	if publishingError != nil {
		return publishingError
	}
	return nil
}
