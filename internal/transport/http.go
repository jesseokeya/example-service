package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nicholaslam/example-service/internal/endpoint"
)

var (
	errBadRouting = errors.New("inconsistent mapping between route and handler")
	errBadRequest = errors.New("bad request")
)

// MakeCreateHTTPHandler mounts the create endpoint.
func MakeCreateHTTPHandler(endpoint kitendpoint.Endpoint) http.Handler {
	return kithttp.NewServer(
		endpoint,
		decodeCreateRequest,
		encodeResponse,
		kithttp.ServerErrorEncoder(encodeError),
	)
}

// MakeReadHTTPHandler mounts the read endpoint.
func MakeReadHTTPHandler(endpoint kitendpoint.Endpoint) http.Handler {
	return kithttp.NewServer(
		endpoint,
		decodeReadRequest,
		encodeResponse,
		kithttp.ServerErrorEncoder(encodeError),
	)
}

// MakeListHTTPHandler mounts the list endpoint.
func MakeListHTTPHandler(endpoint kitendpoint.Endpoint) http.Handler {
	return kithttp.NewServer(
		endpoint,
		decodeListRequest,
		encodeResponse,
		kithttp.ServerErrorEncoder(encodeError),
	)
}

// MakeDeleteHTTPHandler mounts the delete endpoint.
func MakeDeleteHTTPHandler(endpoint kitendpoint.Endpoint) http.Handler {
	return kithttp.NewServer(
		endpoint,
		decodeDeleteRequest,
		encodeResponse,
		kithttp.ServerErrorEncoder(encodeError),
	)
}

func decodeCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest
	}
	return req, nil
}

func decodeReadRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRouting
	}
	return endpoint.ReadRequest{ID: id}, nil
}

func decodeListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	palindromeRaw := strings.ToLower(r.URL.Query().Get("palindrome"))
	if palindromeRaw == "" {
		return endpoint.ListRequest{Palindrome: nil}, nil
	} else if palindromeRaw == "true" {
		palindrome := true
		return endpoint.ListRequest{Palindrome: &palindrome}, nil
	} else if palindromeRaw == "false" {
		palindrome := false
		return endpoint.ListRequest{Palindrome: &palindrome}, nil
	}
	return nil, errBadRequest
}

func decodeDeleteRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRouting
	}
	return endpoint.DeleteRequest{ID: id}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if response == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("cannot encode nil error")
	}
	w.WriteHeader(statusCode(err))
}

func statusCode(err error) int {
	switch err {
	case endpoint.ErrNotFound:
		return http.StatusNotFound
	case endpoint.ErrBadRequest, errBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
