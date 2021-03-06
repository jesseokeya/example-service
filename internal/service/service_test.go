package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nicholaslam/example-service/internal/store"
	"github.com/stretchr/testify/require"
)

type mockStore struct {
	msg  store.Message
	msgs []store.Message
	err  error
}

func (ms *mockStore) Create(ctx context.Context, p store.MessagePayload) (store.Message, error) {
	return ms.msg, ms.err
}

func (ms *mockStore) Read(ctx context.Context, id string) (store.Message, error) {
	return ms.msg, ms.err
}

func (ms *mockStore) List(ctx context.Context, p store.ListPayload) ([]store.Message, error) {
	return ms.msgs, ms.err
}

func (ms *mockStore) Delete(ctx context.Context, id string) error {
	return ms.err
}

func toBoolPointer(b bool) *bool {
	return &b
}

func TestNewService(t *testing.T) {
	require.NotNil(t, NewService(&mockStore{}, true))
}

func TestCreate(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	testCases := []struct {
		name    string
		store   store.Store
		strict  bool
		payload MessagePayload
		want    Message
		errMsg  string
	}{
		{
			"strict palindrome",
			&mockStore{
				store.Message{
					ID:         "123",
					Text:       "a toyota",
					Palindrome: false,
					CreatedAt:  now,
				},
				nil,
				nil,
			},
			true,
			MessagePayload{
				Text: "a toyota",
			},
			Message{
				ID:         "123",
				Text:       "a toyota",
				Palindrome: false,
				CreatedAt:  now,
			},
			"",
		},
		{
			"non-strict palindrome",
			&mockStore{
				store.Message{
					ID:         "123",
					Text:       "a toyota",
					Palindrome: true,
					CreatedAt:  now,
				},
				nil,
				nil,
			},
			false,
			MessagePayload{
				Text: "a toyota",
			},
			Message{
				ID:         "123",
				Text:       "a toyota",
				Palindrome: true,
				CreatedAt:  now,
			},
			"",
		},
		{
			"unhandled error",
			&mockStore{
				store.Message{},
				nil,
				errors.New("error"),
			},
			false,
			MessagePayload{
				Text: "",
			},
			Message{},
			"error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewService(tc.store, tc.strict)
			msg, err := svc.Create(context.Background(), tc.payload)
			if tc.errMsg == "" {
				require.NoError(t, err)
				require.Equal(t, tc.want, msg)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
				require.Empty(t, msg)
			}
		})
	}
}

func TestRead(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	testCases := []struct {
		name   string
		store  store.Store
		id     string
		want   Message
		errMsg string
	}{
		{
			"success",
			&mockStore{
				store.Message{
					ID:         "123",
					Text:       "racecar",
					Palindrome: true,
					CreatedAt:  now,
				},
				nil,
				nil,
			},
			"123",
			Message{
				ID:         "123",
				Text:       "racecar",
				Palindrome: true,
				CreatedAt:  now,
			},
			"",
		},
		{
			"store.ErrNotFound",
			&mockStore{
				store.Message{},
				nil,
				store.ErrNotFound,
			},
			"",
			Message{},
			ErrNotFound.Error(),
		},
		{
			"unhandled error",
			&mockStore{
				store.Message{},
				nil,
				errors.New("error"),
			},
			"",
			Message{},
			"error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewService(tc.store, true)
			msg, err := svc.Read(context.Background(), tc.id)
			if tc.errMsg == "" {
				require.NoError(t, err)
				require.Equal(t, tc.want, msg)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
				require.Empty(t, msg)
			}
		})
	}
}

func TestList(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	testCases := []struct {
		name        string
		listPayload ListPayload
		store       store.Store
		want        []Message
		errMsg      string
	}{
		{
			"no palindrome query",
			ListPayload{Palindrome: nil},
			&mockStore{
				store.Message{},
				[]store.Message{
					{
						ID:         "123",
						Text:       "racecar",
						Palindrome: true,
						CreatedAt:  now,
					},
					{
						ID:         "456",
						Text:       "a toyota",
						Palindrome: false,
						CreatedAt:  now,
					},
					{
						ID:         "789",
						Text:       "abc",
						Palindrome: false,
						CreatedAt:  now,
					},
				},
				nil,
			},
			[]Message{
				{
					ID:         "123",
					Text:       "racecar",
					Palindrome: true,
					CreatedAt:  now,
				},
				{
					ID:         "456",
					Text:       "a toyota",
					Palindrome: false,
					CreatedAt:  now,
				},
				{
					ID:         "789",
					Text:       "abc",
					Palindrome: false,
					CreatedAt:  now,
				},
			},
			"",
		},
		{
			"palindrome=true",
			ListPayload{Palindrome: toBoolPointer(true)},
			&mockStore{
				store.Message{},
				[]store.Message{
					{
						ID:         "123",
						Text:       "racecar",
						Palindrome: true,
						CreatedAt:  now,
					},
				},
				nil,
			},
			[]Message{
				{
					ID:         "123",
					Text:       "racecar",
					Palindrome: true,
					CreatedAt:  now,
				},
			},
			"",
		},
		{
			"palindrome=false",
			ListPayload{Palindrome: toBoolPointer(false)},
			&mockStore{
				store.Message{},
				[]store.Message{
					{
						ID:         "456",
						Text:       "a toyota",
						Palindrome: false,
						CreatedAt:  now,
					},
					{
						ID:         "789",
						Text:       "abc",
						Palindrome: false,
						CreatedAt:  now,
					},
				},
				nil,
			},
			[]Message{
				{
					ID:         "456",
					Text:       "a toyota",
					Palindrome: false,
					CreatedAt:  now,
				},
				{
					ID:         "789",
					Text:       "abc",
					Palindrome: false,
					CreatedAt:  now,
				},
			},
			"",
		},
		{
			"unhandled error",
			ListPayload{Palindrome: nil},
			&mockStore{
				store.Message{},
				nil,
				errors.New("error"),
			},
			nil,
			"error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewService(tc.store, true)
			msgs, err := svc.List(context.Background(), tc.listPayload)
			if tc.errMsg == "" {
				require.NoError(t, err)
				require.Equal(t, tc.want, msgs)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
				require.Empty(t, msgs)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	testCases := []struct {
		name   string
		store  store.Store
		id     string
		errMsg string
	}{
		{
			"success",
			&mockStore{
				store.Message{
					ID:         "123",
					Text:       "racecar",
					Palindrome: true,
					CreatedAt:  now,
				},
				nil,
				nil,
			},
			"123",
			"",
		},
		{
			"store.ErrNotFound",
			&mockStore{
				store.Message{
					ID:         "123",
					Text:       "racecar",
					Palindrome: true,
					CreatedAt:  now,
				},
				nil,
				store.ErrNotFound,
			},
			"456",
			"",
		},
		{
			"unhandled error",
			&mockStore{
				store.Message{},
				nil,
				errors.New("error"),
			},
			"",
			"error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewService(tc.store, true)
			err := svc.Delete(context.Background(), tc.id)
			if tc.errMsg == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestToMessage(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	testCases := []struct {
		name string
		msg  store.Message
		want Message
	}{
		{
			"success",
			store.Message{
				ID:         "123",
				Text:       "racecar",
				Palindrome: true,
				CreatedAt:  now,
			},
			Message{
				ID:         "123",
				Text:       "racecar",
				Palindrome: true,
				CreatedAt:  now,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := toMessage(tc.msg)
			require.Equal(t, tc.want, msg)
		})
	}
}

func TestToSlice(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	testCases := []struct {
		name string
		msgs []store.Message
		want []Message
	}{
		{
			"success",
			[]store.Message{
				{
					ID:         "123",
					Text:       "racecar",
					Palindrome: true,
					CreatedAt:  now,
				},
				{
					ID:         "456",
					Text:       "abc",
					Palindrome: false,
					CreatedAt:  now,
				},
			},
			[]Message{
				{
					ID:         "123",
					Text:       "racecar",
					Palindrome: true,
					CreatedAt:  now,
				},
				{
					ID:         "456",
					Text:       "abc",
					Palindrome: false,
					CreatedAt:  now,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msgs := toSlice(tc.msgs)
			require.Equal(t, tc.want, msgs)
		})
	}
}
