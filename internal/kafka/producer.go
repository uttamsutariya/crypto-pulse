package kafka

import (
	"log"

	"github.com/IBM/sarama"
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

func (p *Producer) Close() {
	if p.producer == nil {
		log.Println("Producer already closed")
		return
	}

	if err := p.producer.Close(); err != nil {
		log.Printf("Failed to close Kafka producer: %v", err)
	}
}
