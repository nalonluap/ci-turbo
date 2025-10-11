package nats

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nalonluap/ci-turbo/services/ingester/internal/interfaces"
	"github.com/nats-io/nats.go"
)

// ensure natsMessage implements the interface
var _ interfaces.Message = (*natsMessage)(nil)

// natsMessage is a wrapper around the nats.Msg type.
type natsMessage struct {
	original *nats.Msg
}

func (m *natsMessage) Payload() []byte {
	return m.original.Data
}

func (m *natsMessage) Original() interface{} {
	return m.original
}

// Consumer is a NATS implementation of the MessageConsumer interface.
type Consumer struct {
	sub *nats.Subscription
	nc  *nats.Conn
}

// NewConsumer creates and configures a new NATS consumer.
func NewConsumer(natsURL, consumerGroup string) (*Consumer, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Create a durable, pull-based consumer. This allows our worker to control the flow of messages.
	sub, err := js.PullSubscribe("metrics.*", consumerGroup)
	if err != nil {
		return nil, fmt.Errorf("failed to create pull subscription: %w", err)
	}

	return &Consumer{sub: sub, nc: nc}, nil
}

// FetchBatch reads a batch of messages from the NATS subscription.
func (c *Consumer) FetchBatch(ctx context.Context, maxMessages int, timeout time.Duration) ([]interfaces.Message, error) {
	// Fetch will block until the timeout is reached or the batch is full.
	msgs, err := c.sub.Fetch(maxMessages, nats.MaxWait(timeout))
	if err != nil {
		// nats.ErrTimeout is a normal occurrence when no messages are available.
		if errors.Is(err, nats.ErrTimeout) {
			return nil, nil // Return an empty slice, not an error
		}
		return nil, fmt.Errorf("failed to fetch messages from NATS: %w", err)
	}

	messages := make([]interfaces.Message, len(msgs))
	for i, msg := range msgs {
		messages[i] = &natsMessage{original: msg}
	}

	return messages, nil
}

// CommitMessages acknowledges to NATS that the messages have been processed.
func (c *Consumer) CommitMessages(ctx context.Context, messages []interfaces.Message) error {
	for _, msg := range messages {
		// Acknowledge each message individually.
		if err := msg.Original().(*nats.Msg).Ack(); err != nil {
			// In a real scenario, you might want to handle this more gracefully,
			// perhaps by trying to ack the others.
			return fmt.Errorf("failed to acknowledge message: %w", err)
		}
	}
	return nil
}

// Close gracefully shuts down the NATS connection.
func (c *Consumer) Close() error {
	if err := c.sub.Unsubscribe(); err != nil {
		return err
	}
	c.nc.Close()
	return nil
}
