package store

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrNotFound is returned if a Message is not found.
	ErrNotFound = errors.New("not found")
)

// Storer describes a store that allows create, read, list, and delete operations on Messages.
type Storer interface {
	Create(ctx context.Context, p MessagePayload) (Message, error)
	Read(ctx context.Context, id string) (Message, error)
	List(ctx context.Context) ([]Message, error)
	Delete(ctx context.Context, id string) error
}

// MessagePayload represents a payload used to create a Message.
type MessagePayload struct {
	Text       string
	Palindrome bool
}

// Message represents a string that may be a palindrome.
type Message struct {
	ID         string    `bson:"_id"`
	Text       string    `bson:"text"`
	Palindrome bool      `bson:"palindrome"`
	CreatedAt  time.Time `bson:"createdAt"`
}