package store

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/insertopt"
)

type mockCollection struct {
	cur    mongo.Cursor
	docRes *mongo.DocumentResult
	insRes *mongo.InsertOneResult
	err    error
}

func (mc *mockCollection) InsertOne(ctx context.Context, document interface{}, opts ...insertopt.One) (*mongo.InsertOneResult, error) {
	return mc.insRes, mc.err
}

func (mc *mockCollection) FindOne(ctx context.Context, filter interface{}, opts ...findopt.One) *mongo.DocumentResult {
	return mc.docRes
}

func (mc *mockCollection) Find(ctx context.Context, filter interface{}, opts ...findopt.Find) (mongo.Cursor, error) {
	return mc.cur, mc.err
}

func (mc *mockCollection) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...findopt.DeleteOne) *mongo.DocumentResult {
	return mc.docRes
}

func TestNewMongoStore(t *testing.T) {
	c := &mockCollection{}
	require.NotNil(t, NewMongoStore(c))
}

func TestMongoStoreCreate(t *testing.T) {
	testCases := []struct {
		name    string
		payload MessagePayload
		mc      MongoCollection
		want    Message
		errMsg  string
	}{
		{
			"success",
			MessagePayload{
				Text:       "racecar",
				Palindrome: true,
			},
			&mockCollection{
				err: nil,
			},
			Message{
				Text:       "racecar",
				Palindrome: true,
			},
			"",
		},
		{
			"error",
			MessagePayload{
				Text:       "racecar",
				Palindrome: true,
			},
			&mockCollection{
				err: errors.New("error"),
			},
			Message{},
			"error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ms := NewMongoStore(tc.mc)
			msg, err := ms.Create(context.Background(), tc.payload)
			if tc.errMsg == "" {
				require.NoError(t, err)
				require.NotEmpty(t, msg.ID)
				require.Equal(t, tc.want.Text, msg.Text)
				require.Equal(t, tc.want.Palindrome, msg.Palindrome)
				require.NotEmpty(t, msg.CreatedAt)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
				require.Empty(t, msg)
			}
		})
	}
}

func TestMongoStoreRead(t *testing.T) {
}

func TestMongoStoreList(t *testing.T) {
}

func TestMongoStoreDelete(t *testing.T) {
}
