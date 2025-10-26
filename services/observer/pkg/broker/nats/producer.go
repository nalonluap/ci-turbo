package nats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

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
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Check if the stream exists
	_, err = js.StreamInfo(streamName)
	streamConfig := &nats.StreamConfig{
		Name:      streamName,
		Subjects:  []string{subjects},
		Retention: nats.WorkQueuePolicy, // Important: ensures messages are kept until a worker processes them
	}

	if err != nil {
		if errors.Is(err, nats.ErrStreamNotFound) {
			_, err = js.AddStream(streamConfig)
			if err != nil {
				nc.Close()
				return nil, fmt.Errorf("failed to create NATS stream: %w", err)
			}
			log.Printf("NATS Stream '%s' created", streamName)
		} else {
			nc.Close()
			return nil, fmt.Errorf("failed to get NATS stream info: %w", err)
		}
	}

	return &Producer{js: js, nc: nc}, nil
}

func (p *Producer) Publish(ctx context.Context, topic string, event *domain.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Add a timeout to the publish operation
	ack, err := p.js.Publish(topic, payload, nats.ExpectStream("METRICS"), nats.AckWait(5*time.Second))
	if err != nil {
		log.Printf("failed to publish message to NATS: %v", err)
		return fmt.Errorf("failed to publish message to NATS: %w", err)
	}

	log.Printf("Message for JobID: %s published successfully to stream %s", event.JobID, ack.Stream)
	return nil
}

// Close gracefully shuts down the connection.
func (p *Producer) Close() error {
	p.nc.Close()
	return nil
}
