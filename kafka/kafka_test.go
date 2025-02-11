package kafka_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	kafkaMessage "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	testKafka "github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/yaninyzwitty/gqlgen-subscriptions-go/kafka"
)

var brokers []string

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Start Kafka container
	container, err := testKafka.Run(ctx,
		"confluentinc/confluent-local:7.5.0",
		testKafka.WithClusterID("test-cluster"),
	)
	if err != nil {
		log.Fatal("Failed to start Kafka container: ", err)
	}

	// Ensure container cleanup after tests
	defer func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			log.Printf("Failed to terminate Kafka container: %s", err)
		}
	}()

	// Retrieve broker addresses
	brokers, err = container.Brokers(ctx)
	if err != nil {
		log.Fatalf("Failed to get Kafka brokers: %s", err)
	}

	// Run the test suite.
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestKafkaConnectivity(t *testing.T) {
	// Set a timeout for the test context.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup Kafka configuration and instance.
	kafkaConfig := &kafka.KafkaInputConfig{
		BootstrapServers: brokers[0],
		// Ensure your CreateKafkaReader is configured to read earlier messages,
		// e.g., by setting StartOffset to kafka.FirstOffset if needed.
	}
	kafkaInstance := kafka.NewKafka(kafkaConfig)

	// Create a Kafka writer to produce a message.
	writer, err := kafkaInstance.CreateKafkaWriter(ctx, "test-topic")
	require.NoError(t, err, "failed to create Kafka writer")
	require.NotNil(t, writer, "writer should not be nil")
	defer writer.Close()

	// Define the test message.
	message := kafkaMessage.Message{
		Key:   []byte("test-key"),
		Value: []byte("test-value"),
	}

	// Write the test message.
	err = writer.WriteMessages(ctx, message)
	require.NoError(t, err, "failed to write message to Kafka")

	// Create a Kafka reader to consume the message.
	reader, err := kafkaInstance.CreateKafkaReader("test-topic", "test-group")
	require.NoError(t, err, "failed to create Kafka reader")
	require.NotNil(t, reader, "reader should not be nil")
	defer reader.Close()

	// Read the message from Kafka.
	msg, err := reader.ReadMessage(ctx)
	require.NoError(t, err, "failed to read message from Kafka")
	require.Equal(t, message.Key, msg.Key, "message key does not match")
	require.Equal(t, message.Value, msg.Value, "message value does not match")
}
