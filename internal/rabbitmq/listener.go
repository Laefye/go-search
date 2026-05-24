package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Laefye/go-search/internal/rabbitmq/events"
	"github.com/Laefye/go-search/internal/service/consumer"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Listener struct {
	consumer  *consumer.ConsumerService
	ch        *amqp.Channel
	queueName string
}

func NewListener(
	queueName string,
	ch *amqp.Channel,
	consumer *consumer.ConsumerService,
) *Listener {
	return &Listener{
		queueName: queueName,
		ch:        ch,
		consumer:  consumer,
	}
}

func (l *Listener) Listen(ctx context.Context) error {
	msgs, err := l.ch.Consume(l.queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-msgs:
			var queryEvent events.SearchQueryEvent

			err := json.Unmarshal(msg.Body, &queryEvent)
			if err != nil {
				fmt.Printf("Failed to unmarshal message: %v\n", err)
				continue
			}

			err = l.consumer.Consume(ctx, &queryEvent)
			if err != nil {
				return fmt.Errorf("failed to consume message: %w", err)
			}
		}
	}
}
