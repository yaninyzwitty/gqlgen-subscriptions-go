package graph

import (
	"github.com/gocql/gocql"
	"github.com/segmentio/kafka-go"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Session *gocql.Session
	Reader  *kafka.Reader
}
