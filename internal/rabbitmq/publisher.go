package rabbitmq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch        *amqp.Channel
	queueName string
}

func NewPublisher(ch *amqp.Channel, queueName string) *Publisher {
	return &Publisher{ch: ch, queueName: queueName}
}

func (p *Publisher) PublishQuery(ctx context.Context, query string) error {
	return p.ch.PublishWithContext(
		ctx,
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(query),
			Timestamp:   time.Now(),
		},
	)
}
