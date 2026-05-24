package main

import (
	"net/http"
	"strconv"

	"github.com/Laefye/go-search/internal/config"
	httphandler "github.com/Laefye/go-search/internal/http"
	"github.com/Laefye/go-search/internal/rabbitmq"
	"github.com/Laefye/go-search/internal/repository"
	"github.com/Laefye/go-search/internal/service/search"
	"github.com/Laefye/go-search/internal/service/top"
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

	mux := http.NewServeMux()

	counterRepo := repository.NewRedisRepository(redis.NewClient(&redis.Options{
		Addr:     config.Redis,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	}))

	topService := top.NewTopService(counterRepo)
	searchPublisher := rabbitmq.NewPublisher(ch, q.Name)
	searchService := search.NewSearchService(searchPublisher)

	handler := httphandler.NewHandler(topService, searchService)
	handler.RegisterRoutes(mux)

	http.ListenAndServe(":"+strconv.Itoa(config.Port), mux)
}
