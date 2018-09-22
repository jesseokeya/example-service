package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/nicholaslam/example-service/internal/endpoint"
	"github.com/nicholaslam/example-service/internal/service"
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

func (ms *mockService) List(ctx context.Context, p service.ListPayload) ([]service.Message, error) {
	return ms.msgs, ms.err
}

func (ms *mockService) Delete(ctx context.Context, id string) error {
	return ms.err
}

func toStringPointer(s string) *string {
	return &s
}

func toBoolPointer(b bool) *bool {
	return &b
}

func TestMakeCreateHTTPHandler(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	testCases := []struct {
		name    string
		payload endpoint.CreateRequest
		svc     service.Service
		status  int
		want    endpoint.MessageResponse
	}{
		{
			"success",
			endpoint.CreateRequest{
				Text: toStringPointer("racecar"),
			},
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
			http.StatusOK,
			endpoint.MessageResponse{
				ID:         "123",
				Text:       "racecar",
				Palindrome: true,
				CreatedAt:  now,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			b, _ := json.Marshal(tc.payload)
			r, _ := http.NewRequest("POST", "/api/v1/messages", bytes.NewReader(b))
			MakeCreateHTTPHandler(endpoint.MakeCreateEndpoint(tc.svc)).ServeHTTP(w, r)
			require.Equal(t, tc.status, w.Code)
			var res endpoint.MessageResponse
			json.Unmarshal(w.Body.Bytes(), &res)
			require.Equal(t, tc.want, res)
		})
	}
}

func TestMakeReadHTTPHandler(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	testCases := []struct {
		name    string
		payload endpoint.CreateRequest
		svc     service.Service
		status  int
		want    endpoint.MessageResponse
	}{
		{
			"success",
			endpoint.CreateRequest{
				Text: toStringPointer("racecar"),
			},
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
			http.StatusOK,
			endpoint.MessageResponse{
				ID:         "123",
				Text:       "racecar",
				Palindrome: true,
				CreatedAt:  now,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			b, _ := json.Marshal(tc.payload)
			r, _ := http.NewRequest("GET", "/api/v1/messages/123", bytes.NewReader(b))
			r = mux.SetURLVars(r, map[string]string{"id": "123"})
			MakeReadHTTPHandler(endpoint.MakeReadEndpoint(tc.svc)).ServeHTTP(w, r)
			require.Equal(t, tc.status, w.Code)
			var res endpoint.MessageResponse
			json.Unmarshal(w.Body.Bytes(), &res)
			require.Equal(t, tc.want, res)
		})
	}
}

func TestMakeListHTTPHandler(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	testCases := []struct {
		name   string
		svc    service.Service
		status int
		want   []endpoint.MessageResponse
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
			http.StatusOK,
			[]endpoint.MessageResponse{
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
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/api/v1/messages", nil)
			MakeListHTTPHandler(endpoint.MakeListEndpoint(tc.svc)).ServeHTTP(w, r)
			require.Equal(t, tc.status, w.Code)
			var res []endpoint.MessageResponse
			json.Unmarshal(w.Body.Bytes(), &res)
			require.Equal(t, tc.want, res)
		})
	}
}

func TestMakeDeleteHTTPHandler(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	testCases := []struct {
		name    string
		payload endpoint.CreateRequest
		svc     service.Service
		status  int
	}{
		{
			"success",
			endpoint.CreateRequest{
				Text: toStringPointer("racecar"),
			},
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
			http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			b, _ := json.Marshal(tc.payload)
			r, _ := http.NewRequest("DELETE", "/api/v1/messages/123", bytes.NewReader(b))
			r = mux.SetURLVars(r, map[string]string{"id": "123"})
			MakeDeleteHTTPHandler(endpoint.MakeDeleteEndpoint(tc.svc)).ServeHTTP(w, r)
			require.Equal(t, tc.status, w.Code)
		})
	}
}

func TestDecodeListRequest(t *testing.T) {
	testCases := []struct {
		name   string
		query  string
		want   endpoint.ListRequest
		errMsg string
	}{
		{
			"no palindrome query",
			"",
			endpoint.ListRequest{Palindrome: nil},
			"",
		},
		{
			"palindrome=true",
			"true",
			endpoint.ListRequest{Palindrome: toBoolPointer(true)},
			"",
		},
		{
			"palindrome=false",
			"false",
			endpoint.ListRequest{Palindrome: toBoolPointer(false)},
			"",
		},
		{
			"invalid palindrome query",
			"invalid",
			endpoint.ListRequest{},
			errBadRequest.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("GET", "/api/v1/messages", nil)
			q := r.URL.Query()
			q.Add("palindrome", tc.query)
			r.URL.RawQuery = q.Encode()
			req, err := decodeListRequest(context.Background(), r)
			if tc.errMsg == "" {
				require.NoError(t, err)
				require.Equal(t, tc.want, req)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
				require.Nil(t, req)
			}
		})
	}
}

func TestDecodeCreateRequestError(t *testing.T) {
	reader := bytes.NewReader([]byte("hello, world!"))
	r, _ := http.NewRequest("POST", "/api/v1/messages", reader)
	req, err := decodeCreateRequest(context.Background(), r)
	require.Error(t, err)
	require.Equal(t, endpoint.ErrBadRequest.Error(), err.Error())
	require.Empty(t, req)
}

func TestDecodeReadRequestError(t *testing.T) {
	r, _ := http.NewRequest("GET", "/api/v1/messages/123", nil)
	req, err := decodeReadRequest(context.Background(), r)
	require.Error(t, err)
	require.Equal(t, errBadRouting.Error(), err.Error())
	require.Empty(t, req)
}

func TestDecodeDeleteRequestError(t *testing.T) {
	r, _ := http.NewRequest("DELETE", "/api/v1/messages/123", nil)
	req, err := decodeDeleteRequest(context.Background(), r)
	require.Error(t, err)
	require.Equal(t, errBadRouting.Error(), err.Error())
	require.Empty(t, req)
}

func TestEncodeResponse(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	testCases := []struct {
		name string
		res  interface{}
		err  error
	}{
		{
			"success",
			endpoint.MessageResponse{
				ID:         "123",
				Text:       "racecar",
				Palindrome: true,
				CreatedAt:  now,
			},
			nil,
		},
		{
			"nil response",
			nil,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := encodeResponse(context.Background(), w, tc.res)
			if tc.res != nil {
				contentType := w.Header()["Content-Type"]
				require.Equal(t, []string{"application/json; charset=utf-8"}, contentType)
			}
			require.NoError(t, err)
		})
	}
}

func TestEncodeError(t *testing.T) {
	testCases := []struct {
		name string
		err  error
		want int
	}{
		{
			"endpoint.ErrNotFound",
			endpoint.ErrNotFound,
			http.StatusNotFound,
		},
		{
			"endpoint.ErrBadRequest",
			endpoint.ErrBadRequest,
			http.StatusBadRequest,
		},
		{
			"unhandled error",
			errors.New("error"),
			http.StatusInternalServerError,
		},
		{
			"nil error",
			nil,
			0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			fn := func() {
				encodeError(context.Background(), tc.err, w)
			}
			if tc.err == nil {
				require.Panics(t, fn)
			} else {
				fn()
				require.Equal(t, tc.want, w.Code)
			}
		})
	}
}

func TestStatusCode(t *testing.T) {
	testCases := []struct {
		name string
		err  error
		want int
	}{
		{
			"endpoint.ErrNotFound",
			endpoint.ErrNotFound,
			http.StatusNotFound,
		},
		{
			"endpoint.ErrBadRequest",
			endpoint.ErrBadRequest,
			http.StatusBadRequest,
		},
		{
			"errBadRequest",
			errBadRequest,
			http.StatusBadRequest,
		},
		{
			"unhandled error",
			errors.New("error"),
			http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code := statusCode(tc.err)
			require.Equal(t, tc.want, code)
		})
	}
}
