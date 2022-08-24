package v1

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Controller struct {
	Env
	Consumer
	Producer
}

func (ctrl *Controller) Run() {

	defer ctrl.CloseConsumer()
	defer ctrl.CloseProducer()

	var batch []kafka.Message
	ctx := context.Background()

	// Receive messages
	for {
		m, err := ctrl.reader.FetchMessage(ctx)
		if err != nil {
			logrus.Warnf("failed to receive messages: %v", err)
			continue
		}

		logrus.Debugf("received message: %v : %v", string(m.Key), string(m.Value))

		// Handle message header
		retrys := ctrl.getHeaderStringValueAsIntByKey(m, ctrl.KafkaRetrysHeaderName)
		if retrys >= ctrl.KafkaRetryLimit {
			logrus.Debugf("maximum retry limit reached for message %s. Sending message to dead letter.", string(m.Value))
			m.Topic = ctrl.KafkaDeadLetterTopic
		} else {
			m = ctrl.updateMessage(retrys, m)
		}

		// Create message batch
		batch = append(batch, m)
		if len(batch) == ctrl.KafkaBatchSize {
			break
		}
	}

	// Publish messages
	for _, m := range batch {
		if err := ctrl.writer.WriteMessages(ctx, kafka.Message{
			Topic:   m.Topic,
			Key:     m.Key,
			Value:   m.Value,
			Headers: m.Headers,
		}); err != nil {
			logrus.Panicf("failed to write message: %v", err)
		}
	}

	// Commit messages
	if err := ctrl.reader.CommitMessages(ctx, batch...); err != nil {
		logrus.Panicf("failed to commit message batch: %v", err)
	}

	logrus.Infof("successfully published all retry messages")
}

func (ctrl *Controller) updateMessage(retrys int, m kafka.Message) kafka.Message {

	retrys++
	retryValue := strconv.Itoa(retrys)

	failedConsumerGroupId := ctrl.getHeaderStringValueByKey(m, ctrl.KafkaFailedConsumerGroupIdHeaderName)

	m.Topic = ctrl.KafkaEingangsTopic

	m.Headers = []kafka.Header{
		{
			Key:   ctrl.KafkaRetrysHeaderName,
			Value: []byte(retryValue),
		},
		{
			Key:   ctrl.KafkaFailedConsumerGroupIdHeaderName,
			Value: []byte(failedConsumerGroupId),
		},
	}

	return m
}

func (ctrl *Controller) getHeaderStringValueAsIntByKey(m kafka.Message, key string) int {

	for _, h := range m.Headers {
		if h.Key == key {

			val := string(h.Value[0:])
			i, err := strconv.Atoi(val)
			if err != nil {
				logrus.Warnf("failed to convert string to int from header: %s\n", key)
				// if parsing failed routing message to deadletter
				return ctrl.KafkaRetryLimit + 1
			}
			return i
		}
	}

	return 0
}

func (ctrl *Controller) getHeaderStringValueByKey(m kafka.Message, key string) string {

	for _, h := range m.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}

	return ""
}
