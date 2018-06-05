package endpoint

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nicholaslam/palindrome-service/internal/service"
	"github.com/stretchr/testify/require"
)

type mockService struct {
	msg  service.Message
	msgs []service.Message
	err  error
}

func (ms *mockService) Create(ctx context.Context, p service.MessagePayload) (service.Message, error) {
	return ms.msg, ms.err
}

func (ms *mockService) Read(ctx context.Context, id string) (service.Message, error) {
	return ms.msg, ms.err
}

func (ms *mockService) List(ctx context.Context) ([]service.Message, error) {
	return ms.msgs, ms.err
}

func (ms *mockService) Delete(ctx context.Context, id string) error {
	return ms.err
}

func toStringPointer(s string) *string {
	return &s
}

func TestMakeCreateEndpoint(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)

	testCases := []struct {
		name   string
		svc    service.Servicer
		req    CreateRequest
		want   MessageResponse
		errMsg string
	}{
		{
			"success",
			&mockService{
				service.Message{
					ID:         "123",
					Text:       "racecar",
					Palindrome: true,
					CreatedAt:  now,
				},
				nil,
				nil,
			},
			CreateRequest{
				Text: toStringPointer("racecar"),
			},
			MessageResponse{
				ID:         "123",
				Text:       "racecar",
				Palindrome: true,
				CreatedAt:  now,
			},
			"",
		},
		{
			"ErrBadRequest",
			&mockService{},
			CreateRequest{},
			MessageResponse{},
			ErrBadRequest.Error(),
		},
		{
			"unhandled error",
			&mockService{
				service.Message{},
				nil,
				errors.New("error"),
			},
			CreateRequest{
				Text: toStringPointer("racecar"),
			},
			MessageResponse{},
			"error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fn := MakeCreateEndpoint(tc.svc)
			res, err := fn(context.Background(), tc.req)
			msgRes, ok := res.(MessageResponse)
			require.True(t, ok)
			if tc.errMsg == "" {
				require.NoError(t, err)
				require.Equal(t, tc.want, msgRes)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
				require.Empty(t, msgRes)
			}
		})
	}
}

func TestMakeReadEndpoint(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)

	testCases := []struct {
		name   string
		svc    service.Servicer
		req    ReadRequest
		want   MessageResponse
		errMsg string
	}{
		{
			"success",
			&mockService{
				service.Message{
					ID:         "123",
					Text:       "racecar",
					Palindrome: true,
					CreatedAt:  now,
				},
				nil,
				nil,
			},
			ReadRequest{
				ID: "123",
			},
			MessageResponse{
				ID:         "123",
				Text:       "racecar",
				Palindrome: true,
				CreatedAt:  now,
			},
			"",
		},
		{
			"service.ErrNotFound",
			&mockService{
				service.Message{},
				nil,
				service.ErrNotFound,
			},
			ReadRequest{
				ID: "123",
			},
			MessageResponse{},
			ErrNotFound.Error(),
		},
		{
			"unhandled error",
			&mockService{
				service.Message{},
				nil,
				errors.New("error"),
			},
			ReadRequest{
				ID: "123",
			},
			MessageResponse{},
			"error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fn := MakeReadEndpoint(tc.svc)
			res, err := fn(context.Background(), tc.req)
			msgRes, ok := res.(MessageResponse)
			require.True(t, ok)
			if tc.errMsg == "" {
				require.NoError(t, err)
				require.Equal(t, tc.want, msgRes)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
				require.Empty(t, msgRes)
			}
		})
	}
}

func TestMakeListEndpoint(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)

	testCases := []struct {
		name   string
		svc    service.Servicer
		want   []MessageResponse
		errMsg string
	}{
		{
			"success",
			&mockService{
				service.Message{},
				[]service.Message{
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
				nil,
			},
			[]MessageResponse{
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
			"",
		},
		{
			"unhandled error",
			&mockService{
				service.Message{},
				nil,
				errors.New("error"),
			},
			nil,
			"error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fn := MakeListEndpoint(tc.svc)
			res, err := fn(context.Background(), nil)
			msgRes, ok := res.([]MessageResponse)
			require.True(t, ok)
			if tc.errMsg == "" {
				require.NoError(t, err)
				require.Equal(t, tc.want, msgRes)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
				require.Empty(t, msgRes)
			}
		})
	}
}

func TestMakeDeleteEndpoint(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)

	testCases := []struct {
		name   string
		svc    service.Servicer
		req    DeleteRequest
		errMsg string
	}{
		{
			"success",
			&mockService{
				service.Message{
					ID:         "123",
					Text:       "racecar",
					Palindrome: true,
					CreatedAt:  now,
				},
				nil,
				nil,
			},
			DeleteRequest{
				ID: "123",
			},
			"",
		},
		{
			"unhandled error",
			&mockService{
				service.Message{},
				nil,
				errors.New("error"),
			},
			DeleteRequest{},
			"error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fn := MakeDeleteEndpoint(tc.svc)
			res, err := fn(context.Background(), tc.req)
			require.Nil(t, res)
			if tc.errMsg == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
			}
		})
	}
}
