package store

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func toBoolPointer(b bool) *bool {
	return &b
}

func TestNewTempStore(t *testing.T) {
	require.NotNil(t, NewTempStore())
}

func TestCreate(t *testing.T) {
	testCases := []struct {
		name    string
		payload MessagePayload
		want    Message
	}{
		{
			"success",
			MessagePayload{
				Text:       "racecar",
				Palindrome: true,
			},
			Message{
				Text:       "racecar",
				Palindrome: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := NewTempStore()
			msg, err := ts.Create(context.Background(), tc.payload)
			require.NoError(t, err)
			require.NotEmpty(t, msg.ID)
			require.Equal(t, tc.want.Text, msg.Text)
			require.Equal(t, tc.want.Palindrome, msg.Palindrome)
			require.NotEmpty(t, msg.CreatedAt)
		})
	}
}

func TestRead(t *testing.T) {
	testCases := []struct {
		name    string
		payload MessagePayload
		errMsg  string
	}{
		{
			"success",
			MessagePayload{},
			"",
		},
		{
			"ErrNotFound",
			MessagePayload{},
			"not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := NewTempStore()
			cMsg, _ := ts.Create(context.Background(), tc.payload)
			if tc.errMsg == "" {
				rMsg, err := ts.Read(context.Background(), cMsg.ID)
				require.NoError(t, err)
				require.Equal(t, cMsg, rMsg)
			} else {
				rMsg, err := ts.Read(context.Background(), "uuid")
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
				require.Empty(t, rMsg)
			}
		})
	}
}

func TestList(t *testing.T) {
	testCases := []struct {
		name            string
		messagePayloads []MessagePayload
		listPayload     ListPayload
		length          int
	}{
		{
			"empty store",
			[]MessagePayload{},
			ListPayload{Palindrome: nil},
			0,
		},
		{
			"no palindrome query",
			[]MessagePayload{
				{
					Text:       "racecar",
					Palindrome: true,
				},
				{
					Text:       "a toyota",
					Palindrome: false,
				},
				{
					Text:       "abc",
					Palindrome: false,
				},
			},
			ListPayload{Palindrome: nil},
			3,
		},
		{
			"palindrome=true",
			[]MessagePayload{
				{
					Text:       "racecar",
					Palindrome: true,
				},
				{
					Text:       "a toyota",
					Palindrome: false,
				},
				{
					Text:       "abc",
					Palindrome: false,
				},
			},
			ListPayload{Palindrome: toBoolPointer(true)},
			1,
		},
		{
			"palindrome=false",
			[]MessagePayload{
				{
					Text:       "racecar",
					Palindrome: true,
				},
				{
					Text:       "a toyota",
					Palindrome: false,
				},
				{
					Text:       "abc",
					Palindrome: false,
				},
			},
			ListPayload{Palindrome: toBoolPointer(false)},
			2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := NewTempStore()
			for _, p := range tc.messagePayloads {
				ts.Create(context.Background(), p)
			}
			msgs, err := ts.List(context.Background(), tc.listPayload)
			require.NoError(t, err)
			require.Equal(t, tc.length, len(msgs))
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		name    string
		payload MessagePayload
		errMsg  string
	}{
		{
			"success",
			MessagePayload{},
			"",
		},
		{
			"ErrNotFound",
			MessagePayload{},
			"not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := NewTempStore()
			msg, _ := ts.Create(context.Background(), tc.payload)
			if tc.errMsg == "" {
				err := ts.Delete(context.Background(), msg.ID)
				require.NoError(t, err)
			} else {
				err := ts.Delete(context.Background(), "uuid")
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
			}
		})
	}
}

func TestToSlice(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	testCases := []struct {
		name string
		m    map[string]Message
	}{
		{
			"empty map",
			map[string]Message{},
		},
		{
			"non-empty map",
			map[string]Message{
				"123": {
					ID:         "123",
					Text:       "racecar",
					Palindrome: true,
					CreatedAt:  now,
				},
				"456": {
					ID:         "456",
					Text:       "a toyota",
					Palindrome: false,
					CreatedAt:  now,
				},
				"789": {
					ID:         "789",
					Text:       "abc",
					Palindrome: false,
					CreatedAt:  now,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msgs := toSlice(tc.m)
			require.Equal(t, len(tc.m), len(msgs))
			for _, v := range tc.m {
				require.Contains(t, msgs, v)
			}
		})
	}
}
