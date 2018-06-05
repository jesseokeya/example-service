package service

import (
	"context"
	"errors"
	"time"

	"github.com/nicholaslam/palindrome-service/internal/store"
	"github.com/nicholaslam/palindrome-service/pkg/palindrome"
)

var (
	// ErrNotFound is returned if a Message is not found.
	ErrNotFound = errors.New("not found")
)

// Servicer describes a service that stores Messages.
type Servicer interface {
	Create(ctx context.Context, p MessagePayload) (Message, error)
	Read(ctx context.Context, id string) (Message, error)
	List(ctx context.Context) ([]Message, error)
	Delete(ctx context.Context, id string) error
}

// MessagePayload represents a payload used to create a Message.
type MessagePayload struct {
	Text string
}

// Message represents a string that may be a palindrome.
type Message struct {
	ID         string
	Text       string
	Palindrome bool
	CreatedAt  string
}

type service struct {
	store            store.Storer
	strictPalindrome bool
}

// NewService returns a new service.
func NewService(s store.Storer, strict bool) Servicer {
	return &service{
		store:            s,
		strictPalindrome: strict,
	}
}

func (s *service) Create(ctx context.Context, p MessagePayload) (Message, error) {
	var pal bool
	if s.strictPalindrome {
		pal = palindrome.IsPalindromeStrict(p.Text)
	} else {
		pal = palindrome.IsPalindrome(p.Text)
	}
	payload := store.MessagePayload{
		Text:       p.Text,
		Palindrome: pal,
	}
	msg, err := s.store.Create(ctx, payload)
	if err != nil {
		return Message{}, err
	}
	return toMessage(msg), nil
}

func (s *service) Read(ctx context.Context, id string) (Message, error) {
	msg, err := s.store.Read(ctx, id)
	if err != nil {
		if err == store.ErrNotFound {
			return Message{}, ErrNotFound
		}
		return Message{}, err
	}
	return toMessage(msg), nil
}

func (s *service) List(ctx context.Context) ([]Message, error) {
	msgs, err := s.store.List(ctx)
	if err != nil {
		return []Message{}, err
	}
	return toSlice(msgs), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	err := s.store.Delete(ctx, id)
	if err == store.ErrNotFound {
		return nil
	}
	return err
}

func toMessage(msg store.Message) Message {
	return Message{
		ID:         msg.ID,
		Text:       msg.Text,
		Palindrome: msg.Palindrome,
		CreatedAt:  msg.CreatedAt.Format(time.RFC3339),
	}
}

func toSlice(msgs []store.Message) []Message {
	var res []Message
	for _, msg := range msgs {
		res = append(res, toMessage(msg))
	}
	return res
}
