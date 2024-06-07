package websocket

import (
	"log"
	"os"
	"sync"

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
