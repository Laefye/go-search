package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/Laefye/go-search/internal/rabbitmq/events"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch        *amqp.Channel
	queueName string
}

func NewPublisher(ch *amqp.Channel, queueName string) *Publisher {
	return &Publisher{ch: ch, queueName: queueName}
}

func (p *Publisher) PublishQuery(ctx context.Context, query events.SearchQueryEvent) error {
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return err
	}

	return p.ch.PublishWithContext(
		ctx,
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        queryBytes,
		},
	)
}
