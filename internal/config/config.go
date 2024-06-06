package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	KafkaBrokers  []string
	KafkaTopic    string
	WebSocketURLs []string
}

func LoadConfig() (*Config, error) {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(".env file found but could not be loaded")
		}
	} else {
		log.Println(".env file not found, loading configuration from environment variables")
	}

	return &Config{
		KafkaBrokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		KafkaTopic:   os.Getenv("KAFKA_CRYPTO_TOPIC"),
		WebSocketURLs: []string{
			os.Getenv("BINANCE_SPOT_WS_URL"),
			os.Getenv("BINANCE_USD_M_FUT_WS_URL"),
			os.Getenv("BINANCE_COIN_M_FUT_WS_URL"),
		},
	}, nil
}