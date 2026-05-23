package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Laefye/go-search/internal/config"
	"github.com/Laefye/go-search/internal/rabbitmq"
	"github.com/Laefye/go-search/internal/repository"
	"github.com/Laefye/go-search/internal/service/consumer"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func main() {
	godotenv.Load()

	config := config.MustLoad()

	conn, err := amqp.Dial(config.RabbitMQ)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	defer ch.Close()

	q, err := rabbitmq.DeclareQueryQueue(ch, config.QueryQueue)
	if err != nil {
		panic(err)
	}

	redis := redis.NewClient(config.CreateRedisOptions())

	counterRepo := repository.NewRedisRepository(redis)

	guardRepo := repository.NewRedisRepository(redis)

	consumerService := consumer.NewConsumerService(counterRepo, guardRepo)
	listener := rabbitmq.NewListener(ch, consumerService)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	if err := listener.Listen(ctx, q.Name); err != nil {
		panic(err)
	}
}
