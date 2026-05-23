package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Port          int
	Redis         string
	RedisPassword string
	RedisDB       int
	RabbitMQ      string
	QueryQueue    string
	WindowMinutes int
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func MustLoad() *Config {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		panic(fmt.Errorf("Invalid PORT value: %v", err))
	}
	redisDb, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		panic(fmt.Errorf("Invalid REDIS_DB value: %v", err))
	}
	windowMinutes, err := strconv.Atoi(getEnv("WINDOW_MINUTES", "5"))
	if err != nil {
		panic(fmt.Errorf("Invalid WINDOW_MINUTES value: %v", err))
	}
	return &Config{
		Port:          port,
		Redis:         getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDb,
		RabbitMQ:      getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		QueryQueue:    getEnv("RABBITMQ_QUERY_QUEUE", "search.query"),
		WindowMinutes: windowMinutes,
	}
}

func (c *Config) CreateRedisOptions() *redis.Options {
	return &redis.Options{
		Addr:     c.Redis,
		Password: c.RedisPassword,
		DB:       c.RedisDB,
	}
}
