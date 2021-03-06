package store

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type mongoStore struct {
	collection *mongo.Collection
}

// NewMongoStore returns a new store that persists Messages in MongoDB.
func NewMongoStore(c *mongo.Collection) Store {
	return &mongoStore{
		collection: c,
	}
}

func (ms *mongoStore) Create(ctx context.Context, p MessagePayload) (Message, error) {
	msg := Message{
		ID:         objectid.New().Hex(),
		Text:       p.Text,
		Palindrome: p.Palindrome,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339Nano),
	}
	_, err := ms.collection.InsertOne(ctx, msg)
	if err != nil {
		return Message{}, err
	}
	return msg, nil
}

func (ms *mongoStore) Read(ctx context.Context, id string) (Message, error) {
	filter := bson.NewDocument(bson.EC.String("_id", id))
	var msg Message
	err := ms.collection.FindOne(ctx, filter).Decode(&msg)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Message{}, ErrNotFound
		}
		return Message{}, err
	}
	return msg, nil
}

func (ms *mongoStore) List(ctx context.Context, p ListPayload) ([]Message, error) {
	var filter *bson.Document
	if p.Palindrome != nil {
		filter = bson.NewDocument(bson.EC.Boolean("palindrome", *p.Palindrome))
	}
	cur, err := ms.collection.Find(ctx, filter)
	defer cur.Close(ctx)
	if err != nil {
		return []Message{}, err
	}
	msgs := []Message{}
	for cur.Next(ctx) {
		var msg Message
		err := cur.Decode(&msg)
		if err != nil {
			return []Message{}, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func (ms *mongoStore) Delete(ctx context.Context, id string) error {
	filter := bson.NewDocument(bson.EC.String("_id", id))
	var msg Message
	err := ms.collection.FindOneAndDelete(ctx, filter).Decode(&msg)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	return nil
}
