package v1

import (
	"encoding/binary"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestController_getHeaderIntValueByKey(t *testing.T) {

	headerKey := "header"

	m := kafka.Message{
		Headers: []kafka.Header{
			{
				Key:   headerKey,
				Value: []byte("2"),
			},
		},
	}

	ctrl := Controller{
		Env:      Env{},
		Consumer: Consumer{},
		Producer: Producer{},
	}

	headerValue := ctrl.getHeaderStringValueAsIntByKey(m, headerKey)

	assert.Equal(t, 2, headerValue)
}

func TestController_getHeaderStringValueByKey(t *testing.T) {

	headerKey := "header"

	m := kafka.Message{
		Headers: []kafka.Header{
			{
				Key:   headerKey,
				Value: []byte("value"),
			},
		},
	}

	ctrl := Controller{
		Env:      Env{},
		Consumer: Consumer{},
		Producer: Producer{},
	}

	headerValue := ctrl.getHeaderStringValueByKey(m, headerKey)

	assert.Equal(t, "value", headerValue)
}

func TestController_updateMessageHeaderWithFailedConsumerGroupHeader(t *testing.T) {

	m := kafka.Message{
		Headers: []kafka.Header{
			{
				Key:   "retrys",
				Value: []byte("2"),
			},
			{
				Key:   "failed.consumer.group.id",
				Value: []byte("value"),
			},
		},
	}

	ctrl := Controller{
		Env: Env{
			Meta: Meta{
				KafkaRetrysHeaderName:                "retrys",
				KafkaFailedConsumerGroupIdHeaderName: "failed.consumer.group.id",
			},
		},
		Consumer: Consumer{},
		Producer: Producer{},
	}

	m = ctrl.updateMessage(2, m)

	assert.Equal(t, 3, ctrl.getHeaderStringValueAsIntByKey(m, ctrl.KafkaRetrysHeaderName))
	assert.Equal(t, "value", ctrl.getHeaderStringValueByKey(m, ctrl.KafkaFailedConsumerGroupIdHeaderName))
}

func TestController_updateMessageHeaderWithoutRetrysHeader(t *testing.T) {

	intValue := make([]byte, 4)
	binary.LittleEndian.PutUint32(intValue, uint32(2))

	m := kafka.Message{
		Headers: []kafka.Header{
			{
				Key:   "failed.consumer.group.id",
				Value: []byte("value"),
			},
		},
	}

	ctrl := Controller{
		Env: Env{
			Meta: Meta{
				KafkaRetrysHeaderName:                "retrys",
				KafkaFailedConsumerGroupIdHeaderName: "failed.consumer.group.id",
			},
		},
		Consumer: Consumer{},
		Producer: Producer{},
	}

	m = ctrl.updateMessage(0, m)

	assert.Equal(t, 1, ctrl.getHeaderStringValueAsIntByKey(m, ctrl.KafkaRetrysHeaderName))
	assert.Equal(t, "value", ctrl.getHeaderStringValueByKey(m, ctrl.KafkaFailedConsumerGroupIdHeaderName))
}

func TestController_updateMessageHeaderWithoutFailedConsumerGroupHeader(t *testing.T) {

	m := kafka.Message{
		Headers: []kafka.Header{
			{
				Key:   "retrys",
				Value: []byte("2"),
			},
		},
	}

	ctrl := Controller{
		Env: Env{
			Meta: Meta{
				KafkaRetrysHeaderName:                "retrys",
				KafkaFailedConsumerGroupIdHeaderName: "failed.consumer.group.id",
			},
		},
		Consumer: Consumer{},
		Producer: Producer{},
	}

	m = ctrl.updateMessage(2, m)

	assert.Equal(t, 3, ctrl.getHeaderStringValueAsIntByKey(m, ctrl.KafkaRetrysHeaderName))
	assert.Equal(t, "", ctrl.getHeaderStringValueByKey(m, ctrl.KafkaFailedConsumerGroupIdHeaderName))
}
