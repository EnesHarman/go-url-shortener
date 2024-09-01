package kafka

import (
	envConfig "com.github/EnesHarman/url-shortener/config"
	"com.github/EnesHarman/url-shortener/internal/kafka/model"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/labstack/gommon/log"
)

type ClickEventProducer struct {
	topic    string
	producer sarama.AsyncProducer
}

func NewClickEventProducer() *ClickEventProducer {
	env, err := envConfig.LoadConfig()
	if err != nil {
		panic("Error while loading kafka config!")
	}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_8_0_0
	brokers := []string{env.Kafka.Broker}
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Error(err)
		panic("Failed to create producer")
	}

	return &ClickEventProducer{
		producer: producer,
		topic:    TOPIC_EVENT,
	}
}

func (p *ClickEventProducer) Produce(event model.ClickEvent) {
	msg, err := json.Marshal(event)
	if err != nil {
		log.Error("Error while serializing the message! %v", event)
		return
	}
	kafkaMsg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(msg),
	}

	p.producer.Input() <- kafkaMsg

	select {
	case success := <-p.producer.Successes():
		fmt.Printf("Message sent to partition %d at offset %d\n", success.Partition, success.Offset)
	case err := <-p.producer.Errors():
		log.Printf("Failed to send message: %v", err)
	}
}

func (p *ClickEventProducer) ShutDown() error {
	if err := p.producer.Close(); err != nil {
		log.Printf("Error closing producer: %v", err)
		return err
	}
	log.Printf("Producer successfully shut down.")
	return nil
}
