package main

import (
	"log"
	"os"

	"github.com/uttamsutariya/crypto-pulse/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		os.Exit(1)
	}

	log.Println("Kafka Broker ::", cfg.KafkaBrokers)
	log.Println("Kafka Topic ::", cfg.KafkaTopic)
	log.Println("WebSocketURLs ::", cfg.WebSocketURLs)
}
