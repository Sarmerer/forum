package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatastore struct {
	Collection *mongo.Collection
	Client     *mongo.Client
}

func Connect(collection string) (md *MongoDatastore, err error) {
	md = &MongoDatastore{}
	uri := "mongodb+srv://sarmerer:Cvtifhbrb125879%21@forum.tmgwm.mongodb.net/forum?retryWrites=true&w=majority"
	if md.Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri)); err != nil {
		return nil, err
	}
	if md.Collection = md.Client.Database("forum").Collection(collection); err != nil {
		return nil, err
	}
	return md, err
}
