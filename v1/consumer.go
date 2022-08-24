package v1

import (
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	reader *kafka.Reader
}

func (c *Consumer) NewConsumer(kafkaBrokers []string, kafkaGroupID string, kafkaUser string, kafkaPassword string,
	kafkaAuthEnabled bool, kafkaRetryTopic string) {

	r := kafka.ReaderConfig{
		Brokers:        kafkaBrokers,
		GroupID:        kafkaGroupID,
		Topic:          kafkaRetryTopic,
		CommitInterval: time.Second,
		QueueCapacity:  5000,
		Dialer:         c.newDialer(kafkaUser, kafkaPassword, kafkaAuthEnabled),
		MaxBytes:       10e6, // 10MB
	}

	c.reader = kafka.NewReader(r)
}

func (c *Consumer) CloseConsumer() {
	err := c.reader.Close()
	if err != nil {
		logrus.Panicf("failed to close consumer: %s", err)
	}
}

func (c *Consumer) newDialer(kafkaUser string, kafkaPassword string, kafkaAuthEnabled bool) *kafka.Dialer {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	if kafkaAuthEnabled {
		mechanism, err := scram.Mechanism(scram.SHA512, kafkaUser, kafkaPassword)
		if err != nil {
			logrus.Panicf("failed to create SASLMechanism to securely transmit the provided credentials to kafka")
		}
		dialer.SASLMechanism = mechanism
	}

	return dialer
}
