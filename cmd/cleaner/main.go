package main

import (
	"context"
	"log"
	"time"

	"github.com/Laefye/go-search/internal/config"
	"github.com/Laefye/go-search/internal/repository"
	"github.com/Laefye/go-search/internal/service/cleaner"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func runTicker(ctx context.Context, cleaner *cleaner.CleanerService) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	cleaner.Clean(ctx, time.Now())

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := cleaner.Clean(ctx, time.Now())
			if err != nil {
				log.Println("Error occurred while cleaning:", err)
			} else {
				log.Println("Cleaner executed at", time.Now())
			}
		}
	}
}

func main() {
	godotenv.Load()

	config := config.MustLoad()

	counterRepo := repository.NewQueryStatsRepository(redis.NewClient(&redis.Options{
		Addr:     config.Redis,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	}))

	cleanerService := cleaner.NewCleanerService(counterRepo, config.MinuteDelay)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go runTicker(ctx, cleanerService)

	<-ctx.Done()
}
