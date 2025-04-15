package producer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Produser struct {
	producer *kafka.Producer
}

func NewProducer(address []string) (*Produser, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	}

	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return &Produser{producer: p}, nil
}

func (p *Produser) Produce(msg string, topic string, key string) error {
	kMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},

		Value: []byte(msg),
		Key:   []byte(key),
	}

	kCh := make(chan kafka.Event)
	if err := p.producer.Produce(kMsg, kCh); err != nil {
		return err
	}

	e := <-kCh
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case *kafka.Error:
		return ev
	default:
		return errors.New("unknow type")
	}
}

func (p *Produser) Close() {
	p.producer.Flush(5000)
	p.producer.Close()
}
