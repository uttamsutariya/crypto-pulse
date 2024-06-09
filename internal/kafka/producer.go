package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/uttamsutariya/crypto-pulse/internal/types"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

// makes sure no new producer created every single time
var (
	producer sarama.SyncProducer
)

// returns new producer
func NewProducer(brokers []string, topic string) (*Producer, error) {
	if producer != nil {
		return &Producer{producer: producer}, nil
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	var err error
	producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (p *Producer) SendMessage(batchMessage []types.KafkaMessage) error {
	kafkaBatchMessages := make([]*sarama.ProducerMessage, 0)

	for _, msg := range batchMessage {
		messageBytes, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling kafka message :: %v", err)
		}

		kafkaMsg := &sarama.ProducerMessage{
			Topic: p.topic,
			Value: sarama.ByteEncoder(messageBytes),
		}

		kafkaBatchMessages = append(kafkaBatchMessages, kafkaMsg)
	}

	err := producer.SendMessages(kafkaBatchMessages)
	if err != nil {
		log.Printf("Error sending kafka messages :: %v", err)
	}

	return err
}

func (p *Producer) Close() {
	if p.producer == nil {
		log.Println("Producer already closed")
		return
	}

	if err := p.producer.Close(); err != nil {
		log.Printf("Failed to close Kafka producer: %v", err)
	}
}
