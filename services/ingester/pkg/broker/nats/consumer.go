package nats

import (
	"context"
	"errors"
	"fmt"
	"log"
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
func NewConsumer(natsURL, streamName, consumerGroup string) (*Consumer, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Create a durable, pull-based consumer with an explicit acknowledgement policy.
	sub, err := js.PullSubscribe("", consumerGroup, nats.BindStream(streamName))
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create pull subscription: %w", err)
	}

	log.Printf("NATS Consumer '%s' is bound to stream '%s'", consumerGroup, streamName)
	return &Consumer{sub: sub, nc: nc}, nil
}

func (c *Consumer) FetchBatch(ctx context.Context, maxMessages int, timeout time.Duration) ([]interfaces.Message, error) {
	msgs, err := c.sub.Fetch(maxMessages, nats.MaxWait(timeout))
	if err != nil {
		if errors.Is(err, nats.ErrTimeout) {
			return nil, nil // Not an error, just no messages
		}
		return nil, fmt.Errorf("failed to fetch from NATS: %w", err)
	}

	messages := make([]interfaces.Message, len(msgs))
	for i, msg := range msgs {
		messages[i] = &natsMessage{original: msg}
	}
	return messages, nil
}

func (c *Consumer) CommitMessages(ctx context.Context, messages []interfaces.Message) error {
	for _, msg := range messages {
		if err := msg.Original().(*nats.Msg).Ack(); err != nil {
			return fmt.Errorf("failed to ack message: %w", err)
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
