package helpers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gocql/gocql"
	"github.com/segmentio/kafka-go"
)

// ProcessMessages fetches unpublished messages from the Cassandra table,
// publishes them to Kafka, and marks them as published in the Cassandra table.
func ProcessMessages(ctx context.Context, session *gocql.Session, writer *kafka.Writer) error {
	// Define the query to fetch unpublished messages
	const messagesOutboxQuery = `
		SELECT room_id, published, event_id, chat_id, event_type, payload, created_at
		FROM products_keyspace.message_by_room_outbox
		WHERE published = false;
	`

	iter := session.Query(messagesOutboxQuery).WithContext(ctx).Iter()

	// Variables to hold query results
	var (
		roomID    int64
		published bool
		eventID   gocql.UUID
		chatID    int64
		eventType string
		payload   string
		createdAt time.Time
	)

	// Iterate over the results
	for iter.Scan(&roomID, &published, &eventID, &chatID, &eventType, &payload, &createdAt) {
		// Prepare the Kafka message
		kafkaMessage := kafka.Message{
			Key:   []byte(eventID.String()),
			Value: []byte(payload),
		}

		// Publish the message to Kafka
		if err := writer.WriteMessages(ctx, kafkaMessage); err != nil {
			return fmt.Errorf("failed to publish message to Kafka: %w", err)
		}

		slog.Info("Processed message with eventId", "eventID", eventID)

		// Update the Cassandra table to mark the message as published
		const updateQuery = `
		UPDATE products_keyspace.message_by_room_outbox
		SET published = true
		WHERE room_id = ? AND event_id = ?;
	`
		if err := session.Query(updateQuery, roomID, eventID).Exec(); err != nil {
			return fmt.Errorf("failed to update message status in Cassandra: %w", err)
		}

		slog.Info("Updated message status in Cassandra", "eventID", eventID)
	}

	// Close the iterator and handle any errors
	if err := iter.Close(); err != nil {
		return fmt.Errorf("error iterating over Cassandra rows: %w", err)
	}

	return nil
}

// users-227263068692869121, 227263213010481153, 227263253712007169
// rooms- 227263386235236353, 227263437556740097, 227263480875511809, 227263506896973825
// messages-
