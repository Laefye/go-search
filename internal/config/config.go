package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port          int
	Redis         string
	RedisPassword string
	RedisDB       int
	RabbitMQ      string
	QueryQueue    string
	MinuteDelay   int
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
	minuteDelay, err := strconv.Atoi(getEnv("MINUTE_DELAY", "1"))
	if err != nil {
		panic(fmt.Errorf("Invalid MINUTE_DELAY value: %v", err))
	}
	return &Config{
		Port:          port,
		Redis:         getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDb,
		RabbitMQ:      getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		QueryQueue:    getEnv("RABBITMQ_QUERY_QUEUE", "search.query"),
		MinuteDelay:   minuteDelay,
	}
}
