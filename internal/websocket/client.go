package websocket

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/uttamsutariya/crypto-pulse/internal/config"
)

type WebSocketClient struct {
	Url                        string
	Conn                       *websocket.Conn
	messageCh                  chan []byte
	closeCh                    chan struct{}
	reconnectCh                chan struct{}
	wg                         *sync.WaitGroup
	mutex                      sync.Mutex
	consecutiveConnectionFails int // can have multiple connection, for each connection maintain `consecutiveConnectionFails`
}

func NewWebSocketClient(url string, wg *sync.WaitGroup, closeCh chan struct{}) *WebSocketClient {
	return &WebSocketClient{
		Url:                        url,
		messageCh:                  make(chan []byte),
		reconnectCh:                make(chan struct{}),
		closeCh:                    closeCh,
		wg:                         wg,
		consecutiveConnectionFails: 0,
	}
}

func (client *WebSocketClient) Connect() {
	var err error
	config := config.GetConfig()

	for client.consecutiveConnectionFails < config.WebSocketConnectionFailThreshold {
		client.Conn, _, err = websocket.DefaultDialer.Dial(client.Url, nil)
		if err == nil {
			// Connection successful
			client.consecutiveConnectionFails = 0
			log.Printf("Connected to WebSocket :: %s", client.Url)
			client.wg.Add(1)
			go client.readMessages()
			return
		}

		// Connection failed
		client.consecutiveConnectionFails++
		log.Printf("Failed to connect to WebSocket %s (attempt %d/%d): %v", client.Url, client.consecutiveConnectionFails, config.WebSocketConnectionFailThreshold, err)
	}

	// exhausted the failThreshold
	log.Println("All retries exhausted to connect WebSocket, shutting down!")
	os.Exit(1)
}

func (client *WebSocketClient) readMessages() {
	defer client.wg.Done()

	// Timer for tracking idle time
	idleTimer := time.NewTimer(1 * time.Minute)
	defer idleTimer.Stop()

	for {
		select {
		case <-client.closeCh:
			log.Printf("Closing WebSocket connection to %s", client.Url)
			return
		case <-idleTimer.C:
			log.Printf("No data points received from last 1 minute. Reconnecting...")
			client.reconnectCh <- struct{}{}
			return
		default:
			_, message, err := client.Conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message from WebSocket %s: %v", client.Url, err)
				client.reconnectCh <- struct{}{}
				return
			}
			client.messageCh <- message

			// Resetting the idle timer when a new data point is received
			if !idleTimer.Stop() {
				<-idleTimer.C
			}
			idleTimer.Reset(1 * time.Minute)
		}
	}
}

func (client *WebSocketClient) Reconnect() {
	for range client.reconnectCh {
		client.mutex.Lock()
		client.Conn.Close()
		client.Connect()
		client.mutex.Unlock()
	}
}

func (client *WebSocketClient) Close() {
	err := client.Conn.Close()
	if err != nil {
		log.Printf("Error closing WebSocket connection for :: %s", client.Url)
	}
}

func (client *WebSocketClient) Messages() <-chan []byte {
	return client.messageCh
}
