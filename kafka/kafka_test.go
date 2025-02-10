package kafka_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

var kafkaContainer *kafka.KafkaContainer

func TestKafka(t *testing.T) {
	ctx := context.Background()
	var err error
	kafkaContainer, err = kafka.Run(ctx, "confluentinc/confluent-local:7.5.0", kafka.WithClusterID("test-cluster"))
	if err != nil {
		assert.NoError(t, err, "Kafka error shouldnt start with an error")
	}

}
