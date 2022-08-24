package main

import (
	v1 "retry/v1"
	"strings"
)

const version = "v1"

var (
	ctrl v1.Controller
)

func init() {
	ctrl.Cobra()
	ctrl.ServiceVersion = version

	// Workaround for known bug in lib Viper: https://github.com/spf13/viper/issues/380
	kafkaBrokers := strings.Split(ctrl.KafkaBrokers, ",")

	ctrl.NewConsumer(kafkaBrokers, ctrl.KafkaConsumerGroupID, ctrl.KafkaUser, ctrl.KafkaPassword,
		ctrl.KafkaAuthEnabled, ctrl.KafkaRetryTopic)
	ctrl.NewProducer(kafkaBrokers, ctrl.KafkaUser, ctrl.KafkaPassword, ctrl.KafkaAuthEnabled)
	ctrl.Log()
}

func main() {
	ctrl.Run()
}
