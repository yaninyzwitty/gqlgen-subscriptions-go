package kafka

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

// KafkaMethods defines the methods available for Kafka operations.
type KafkaMethods interface {
	CreateKafkaWriter(ctx context.Context, topic string) (*kafka.Writer, error)
	CreateKafkaReader(topic, groupID string) (*kafka.Reader, error)
}

// KafkaInputConfig holds the configuration needed to connect to Kafka.
type KafkaInputConfig struct {
	Username         string
	Password         string
	BootstrapServers string
	SecurityProtocol string
	SASLMechanism    string
	RegistryUrl      string
}

// NewKafka initializes and returns an instance of KafkaInputConfig as KafkaMethods.
func NewKafka(cfg *KafkaInputConfig) KafkaMethods {
	return &KafkaInputConfig{
		Username:         cfg.Username,
		Password:         cfg.Password,
		BootstrapServers: cfg.BootstrapServers,
		SecurityProtocol: cfg.SecurityProtocol,
		SASLMechanism:    cfg.SASLMechanism,
		RegistryUrl:      cfg.RegistryUrl,
	}
}

// createDialer creates a dialer for Kafka connections with SASL and TLS.
func (c *KafkaInputConfig) createDialer() (*kafka.Dialer, error) {
	// Create the SASL mechanism
	mechanism := plain.Mechanism{
		Username: c.Username,
		Password: c.Password,
	}

	// Dialer for connecting securely to Kafka
	dialer := &kafka.Dialer{
		SASLMechanism: mechanism,
		TLS:           &tls.Config{},
	}

	return dialer, nil
}

// CreateKafkaWriter sets up a Kafka writer for producing messages to a specific topic.
func (c *KafkaInputConfig) CreateKafkaWriter(ctx context.Context, topic string) (*kafka.Writer, error) {
	// Create dialer using the shared method
	dialer, err := c.createDialer()
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka writer dialer: %w", err)
	}

	// Create Kafka writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{c.BootstrapServers},
		Topic:   topic,
		Async:   true,

		Dialer:       dialer,
		RequiredAcks: int(kafka.RequireOne),
		Balancer:     &kafka.LeastBytes{},
	})

	if writer == nil {
		return nil, fmt.Errorf("failed to create Kafka writer")
	}

	slog.Info("Kafka writer created succesfully")

	return writer, nil
}

// CreateKafkaReader sets up a Kafka reader for consuming messages from a specific topic and group ID.
func (c *KafkaInputConfig) CreateKafkaReader(topic, groupID string) (*kafka.Reader, error) {
	// Create dialer using the shared method
	dialer, err := c.createDialer()
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka reader dialer: %w", err)
	}

	// Create Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{c.BootstrapServers},
		Topic:          topic,
		GroupID:        groupID,
		Dialer:         dialer,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: 1 * time.Second,
	})

	if reader == nil {
		return nil, fmt.Errorf("failed to create Kafka reader")
	}

	slog.Info("Kafka reader created succesfully")

	return reader, nil
}
