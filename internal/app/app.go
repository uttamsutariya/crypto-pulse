package app

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/uttamsutariya/crypto-pulse/internal/config"
	"github.com/uttamsutariya/crypto-pulse/internal/kafka"
	"github.com/uttamsutariya/crypto-pulse/internal/websocket"
	"github.com/uttamsutariya/crypto-pulse/pkg/utils"
)

type App struct {
	producer  *kafka.Producer
	wg        sync.WaitGroup
	closeCh   chan struct{}
	wsClients []*websocket.WebSocketClient
}

func NewApp(cfg *config.Config) *App {
	var err error
	var producer *kafka.Producer
	config := config.GetConfig()
	consecutiveFails := 0

	for consecutiveFails <= config.KafkaConnectionFailThreshold {
		producer, err = kafka.NewProducer(config.KafkaBrokers, config.KafkaTopic)
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

	if consecutiveFails >= config.KafkaConnectionFailThreshold {
		log.Fatalf("All retries exhusted to run kafka producer: %v", err)
		log.Fatal("Shutting down service...")
		os.Exit(1)
	}

	return &App{
		closeCh:  make(chan struct{}),
		producer: producer,
	}
}

func (app *App) Start() {
	defer app.producer.Close()
	config := config.GetConfig()
	spotUrl := config.WebSocketURLs[0]

	for _, url := range config.WebSocketURLs {
		wsClient := websocket.NewWebSocketClient(url, &app.wg, app.closeCh)
		app.wsClients = append(app.wsClients, wsClient)
		wsClient.Connect()
		go wsClient.Reconnect()

		app.wg.Add(1)
		go func(client *websocket.WebSocketClient, url string) {
			defer app.wg.Done()
			for batch := range client.Messages() {
				select {
				case <-app.closeCh:
					return
				default:
					var jsonBatch []map[string]interface{}
					if err := json.Unmarshal([]byte(batch), &jsonBatch); err != nil {
						log.Printf("Failed to unmarshal batch :: %v", err)
					}

					isSpotSymbol := spotUrl == url
					kafkaBatchMessages := utils.FormatExchangeMessageToKafka(jsonBatch, isSpotSymbol)
					if err := app.producer.SendMessage(kafkaBatchMessages); err != nil {
						log.Printf("Failed to send message to Kafka: %v", err)
					}
				}
			}
		}(wsClient, url)
	}

	app.handleInterrupt()
	app.wg.Wait()
}

func (app *App) handleInterrupt() {
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
}
