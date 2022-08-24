package v1

import (
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/sirupsen/logrus"
)

type Producer struct {
	writer *kafka.Writer
}

func (p *Producer) NewProducer(kafkaBrokers []string, kafkaUser string, kafkaPassword string, kafkaAuthEnabled bool) {

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	if kafkaAuthEnabled {
		mechanism, err := scram.Mechanism(scram.SHA512, kafkaUser, kafkaPassword)
		if err != nil {
			logrus.Panicf("Failed to create sasl.Mechanism to securely transmit the provided credentials to Kafka")
		}
		dialer.SASLMechanism = mechanism
	}

	p.writer = &kafka.Writer{
		Addr:     kafka.TCP(kafkaBrokers...),
		Balancer: &kafka.LeastBytes{},
		Transport: &kafka.Transport{
			Dial: dialer.DialFunc,
			SASL: dialer.SASLMechanism,
		},
	}
}

func (p *Producer) CloseProducer() {
	err := p.writer.Close()
	if err != nil {
		logrus.Panicf("failed to close writer: %s", err)
	}
}
