package datastore

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoDataStore struct {
	client *mongo.Client
	cursor *mongo.Cursor
}

func NewMongoDataStore(uri, db, collection string) (DataStore, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	//ctx.Done()
	cur, err := client.Database(db).Collection(collection).Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}

	return &mongoDataStore{client: client, cursor: cur}, nil
}

func (d *mongoDataStore) HasNext() bool {
	return d.cursor.Next(context.Background())
}

func (d *mongoDataStore) Next() ([]byte, error) {
	doc := make(map[string]interface{})
	err := d.cursor.Decode(&doc)
	if err != nil {
		return nil, err
	}

	marshalledDoc, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	return marshalledDoc, nil
}

func (d *mongoDataStore) Close() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := d.cursor.Close(ctx)
	if err != nil {
		return err
	}

	return d.client.Disconnect(ctx)
}
