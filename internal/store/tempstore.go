package store

import (
	"context"
	"time"

	"github.com/satori/go.uuid"
)

type tempStore struct {
	messages map[string]Message
}

// NewTempStore returns a new store that persists Messages in memory.
func NewTempStore() Store {
	return &tempStore{
		messages: map[string]Message{},
	}
}

func (ts *tempStore) Create(ctx context.Context, p MessagePayload) (Message, error) {
	id := uuid.NewV4().String()
	msg := Message{
		ID:         id,
		Text:       p.Text,
		Palindrome: p.Palindrome,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339Nano),
	}
	ts.messages[id] = msg
	return msg, nil
}

func (ts *tempStore) Read(ctx context.Context, id string) (Message, error) {
	msg, ok := ts.messages[id]
	if !ok {
		return Message{}, ErrNotFound
	}
	return msg, nil
}

func (ts *tempStore) List(ctx context.Context, p ListPayload) ([]Message, error) {
	msgs := toSlice(ts.messages)
	if p.Palindrome == nil {
		return msgs, nil
	}
	var retMsgs []Message
	for _, m := range msgs {
		if m.Palindrome == *p.Palindrome {
			retMsgs = append(retMsgs, m)
		}
	}
	return retMsgs, nil
}

func (ts *tempStore) Delete(ctx context.Context, id string) error {
	_, ok := ts.messages[id]
	if !ok {
		return ErrNotFound
	}
	delete(ts.messages, id)
	return nil
}

func toSlice(m map[string]Message) []Message {
	s := make([]Message, len(m))
	i := 0
	for _, v := range m {
		s[i] = v
		i++
	}
	return s
}
