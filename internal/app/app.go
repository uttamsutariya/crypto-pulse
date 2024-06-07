package app

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/uttamsutariya/crypto-pulse/internal/config"
	"github.com/uttamsutariya/crypto-pulse/internal/kafka"
)

type App struct {
	config   *config.Config
	producer *kafka.Producer
	wg       sync.WaitGroup
	closeCh  chan struct{}
}

func NewApp(cfg *config.Config) *App {
	return &App{
		config:  cfg,
		closeCh: make(chan struct{}),
	}
}

func (app *App) Start() {
	consecutiveFails := 0
	var err error

	for consecutiveFails <= app.config.KafkaConnectionFailThreshold {
		app.producer, err = kafka.NewProducer(app.config.KafkaBrokers, app.config.KafkaTopic)
		if err != nil {
			log.Fatalf("Failed to create Kafka producer: %v", err)
			consecutiveFails++
			time.Sleep(2 * time.Second)
		} else {
			log.Println("Kafka Producer Initialized!")
			consecutiveFails = 0
			break
		}
	}

	if consecutiveFails >= app.config.KafkaConnectionFailThreshold {
		log.Fatalf("All retries exhusted to run kafka producer: %v", err)
		log.Fatal("Shutting down service...")
		os.Exit(1)
	}

	defer app.producer.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-interrupt
		log.Println("Interrupt signal received, shutting down...")

		close(app.closeCh)

		// Close Kafka producer if it was initialized
		if app.producer != nil {
			app.producer.Close()
			log.Println("Kafka producer closed!")
		}

		log.Println("Shutdown complete!")
		os.Exit(0)
	}()

	app.wg.Wait()
}
