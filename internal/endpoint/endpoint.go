// Package endpoint implements utilities to marshal and unmarshal JSON.
package endpoint

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/nicholaslam/example-service/internal/service"
)

var (
	// ErrNotFound is returned if a Message is not found.
	ErrNotFound = errors.New("not found")

	// ErrBadRequest is returned if a request is invalid.
	ErrBadRequest = errors.New("bad request")
)

// CreateRequest represents a payload used to create a Message.
type CreateRequest struct {
	Text *string `json:"text,omitempty"`
}

// ReadRequest represents a payload used to read a Message.
type ReadRequest struct {
	ID string `json:"id"`
}

// ListRequest represents a payload used to list Messages.
type ListRequest struct {
	Palindrome *bool
}

// DeleteRequest represents a payload used to delete a Message.
type DeleteRequest struct {
	ID string `json:"id"`
}

// MessageResponse represents a single Message response.
type MessageResponse struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	Palindrome bool   `json:"palindrome"`
	CreatedAt  string `json:"createdAt"`
}

// MakeCreateEndpoint returns a new endpoint for creating Messages.
func MakeCreateEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		if req.Text == nil {
			return MessageResponse{}, ErrBadRequest
		}
		p := service.MessagePayload{
			Text: *req.Text,
		}
		msg, err := svc.Create(ctx, p)
		if err != nil {
			return MessageResponse{}, err
		}
		return toMessageResponse(msg), nil
	}
}

// MakeReadEndpoint returns a new endpoint for reading Messages.
func MakeReadEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ReadRequest)
		msg, err := svc.Read(ctx, req.ID)
		if err != nil {
			if err == service.ErrNotFound {
				return MessageResponse{}, ErrNotFound
			}
			return MessageResponse{}, err
		}
		return toMessageResponse(msg), nil
	}
}

// MakeListEndpoint returns a new endpoint for listing Messages.
func MakeListEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListRequest)
		p := service.ListPayload{
			Palindrome: req.Palindrome,
		}
		msgs, err := svc.List(ctx, p)
		if err != nil {
			return []MessageResponse{}, err
		}
		res := []MessageResponse{}
		for _, msg := range msgs {
			res = append(res, toMessageResponse(msg))
		}
		return res, nil
	}
}

// MakeDeleteEndpoint returns a new endpoint for deleting Messages.
func MakeDeleteEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)
		err := svc.Delete(ctx, req.ID)
		return nil, err
	}
}

func toMessageResponse(msg service.Message) MessageResponse {
	return MessageResponse{
		ID:         msg.ID,
		Text:       msg.Text,
		Palindrome: msg.Palindrome,
		CreatedAt:  msg.CreatedAt,
	}
}
