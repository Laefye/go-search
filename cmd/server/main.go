package main

import (
	"net/http"
	"strconv"

	"github.com/Laefye/go-search/internal/config"
	httphandler "github.com/Laefye/go-search/internal/http"
	"github.com/Laefye/go-search/internal/repository"
	"github.com/Laefye/go-search/internal/service/top"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	godotenv.Load()

	config := config.MustLoad()

	mux := http.NewServeMux()

	counterRepo := repository.NewRedisQueryStatsRepository(redis.NewClient(&redis.Options{
		Addr:     config.Redis,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	}))

	topService := top.NewTopService(counterRepo, config.WindowMinutes)

	handler := httphandler.NewHandler(topService)
	handler.RegisterRoutes(mux)

	http.ListenAndServe(":"+strconv.Itoa(config.Port), mux)
}
