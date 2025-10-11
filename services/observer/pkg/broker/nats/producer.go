package nats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nalonluap/ci-turbo/services/observer/internal/interfaces"

	"github.com/nalonluap/ci-turbo/internal/domain"
	"github.com/nats-io/nats.go"
)

// ensure Producer implements the MessageBroker interface
var _ interfaces.Broker = (*Producer)(nil)

// Producer is a NATS implementation for sending messages.
type Producer struct {
	js nats.JetStreamContext
	nc *nats.Conn
}

// NewProducer creates and configures a new NATS producer.
func NewProducer(natsURL, streamName, subjects string) (*Producer, error) {
	// Connect to the NATS server
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Create a JetStream context for guaranteed delivery
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Ensure the stream exists. This is an idempotent operation.
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subjects},
	})
	if err != nil && !errors.Is(err, nats.ErrStreamNameAlreadyInUse) {
		nc.Close()
		return nil, fmt.Errorf("failed to add NATS stream: %w", err)
	}

	return &Producer{js: js, nc: nc}, nil
}

// Publish serializes a metric event and sends it to the NATS stream.
func (p *Producer) Publish(ctx context.Context, topic string, event *domain.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event to JSON: %w", err)
	}

	// Publish the message. JetStream handles the reliability guarantees.
	_, err = p.js.Publish(topic, payload)
	if err != nil {
		return fmt.Errorf("failed to publish message to NATS: %w", err)
	}

	return nil
}

// Close gracefully shuts down the connection.
func (p *Producer) Close() error {
	p.nc.Close()
	return nil
}
