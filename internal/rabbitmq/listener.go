package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/Laefye/go-search/internal/service/consumer"
	"github.com/Laefye/go-search/internal/service/dto"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Listener struct {
	consumer *consumer.ConsumerService
	ch       *amqp.Channel
}

func NewListener(
	ch *amqp.Channel,
	consumer *consumer.ConsumerService,
) *Listener {
	return &Listener{
		ch:       ch,
		consumer: consumer,
	}
}

func (l *Listener) Listen(ctx context.Context, queueName string) error {
	msgs, err := l.ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-msgs:
			err := l.consumer.Consume(ctx, &dto.SearchQueryEvent{
				Query:     string(msg.Body),
				Timestamp: time.Now(),
			})
			if err != nil {
				return fmt.Errorf("failed to consume message: %w", err)
			}
		}
	}
}
