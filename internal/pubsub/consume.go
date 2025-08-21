package pubsub

import amqp "github.com/rabbitmq/amqp091-go"

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType simpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	var isDurable bool
	if queueType == Durable {
		isDurable = true
	} else {
		isDurable = false
	}

	queue, err := channel.QueueDeclare(queueName, isDurable, !isDurable, !isDurable, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	bindErr := channel.QueueBind(queueName, key, exchange, false, nil)
	if bindErr != nil {
		return nil, amqp.Queue{}, bindErr
	}

	return channel, queue, nil
}
