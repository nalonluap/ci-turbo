package interfaces

import (
	"context"
	"time"
)

// Message is a generic representation of a message received from a broker.
// It contains the raw data (payload) and metadata needed for committing.
type Message interface {
	// Payload returns the raw message content.
	Payload() []byte
	// Original returns the underlying message object from the specific broker library.
	// This is useful for the adapter layer when committing.
	Original() interface{}
}

// Consumer defines the contract for reading messages from a queue.
// This is the "port" through which the worker receives data to process.
type Consumer interface {
	// FetchBatch attempts to read a batch of messages from the topic.
	// It should block for a certain amount of time if no messages are available.
	// `maxMessages` is the maximum number of messages to return in one batch.
	// `timeout` is the maximum time to wait for messages.
	FetchBatch(ctx context.Context, maxMessages int, timeout time.Duration) ([]Message, error)

	// CommitMessages tells the broker that a batch of messages has been successfully processed
	// and should not be delivered again.
	CommitMessages(ctx context.Context, messages []Message) error

	// Close gracefully shuts down the connection to the broker.
	Close() error
}
