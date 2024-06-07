package main

import (
	"log"
	"os"

	"github.com/uttamsutariya/crypto-pulse/internal/app"
	"github.com/uttamsutariya/crypto-pulse/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		os.Exit(1)
	}

	log.Println("Starting crypto pulse!")

	application := app.NewApp(cfg)
	application.Start()
}
