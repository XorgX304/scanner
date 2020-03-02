package dao

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Connect to db
func Connect(URI string, Db string, Coll string) (*mongo.Collection, error) {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}
	collection := client.Database(Db).Collection(Coll)
	return collection, nil
}

// Drop collection
func Drop(c *mongo.Collection) error {
	if err := c.Drop(context.TODO()); err != nil {
		return err
	}
	return nil
}

//InsertOne document
func InsertOne(c *mongo.Collection, d interface{}) error {
	_, err := c.InsertOne(context.TODO(), d)
	if err != nil {
		return err
	}
	return nil
}

//InsertMany document
func InsertMany(c *mongo.Collection, d []interface{}) error {
	_, err := c.InsertMany(context.TODO(), d)
	if err != nil {
		return err
	}
	return nil
}
