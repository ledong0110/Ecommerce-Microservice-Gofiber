package goose

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/jinzhu/copier"
)


var mongodb *mongo.Database

type Combination[T any] struct {
	Collection string
}

// Some method may be missed, due to my laziness :))
type Methods[T any] interface {
	Find(interface{}, ...*options.FindOptions) ([]T, error)
	FindOne(interface{}, ...*options.FindOneOptions) (T, error)
	FindOneAndReplace() // Not yet
	FindOneAndDelete()	// Not yet
	FindOneAndUpdate()  // Not yet
	UpdateOne(interface{}, interface{}) (*mongo.UpdateResult, error)
	UpdateMany(interface{}, interface{}) (*mongo.UpdateResult, error)
	UpdateById()		// Not yet
	InsertOne(T) (*mongo.InsertOneResult, error)
	InsertMany([]T) (*mongo.InsertManyResult, error)
	DeleteOne(interface{}) (*mongo.DeleteResult, error)
	DeleteMany(interface{}) (*mongo.DeleteResult, error)
	ReplaceOne(interface{}, interface{}) (*mongo.UpdateResult, error)
	Drop() error
	Count(interface{}) (int,error)
}


func (comb Combination[T]) Find(filter interface{}, opts ...*options.FindOptions) ([]T, error){
	var elementDetails []T
	accessCollection := mongodb.Collection(comb.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, queryError := accessCollection.Find(ctx, filter, opts...)
	defer cancel()
	if queryError != nil {
		log.Println(queryError)
		return elementDetails, queryError
	}
	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var referSchema T
		err := cursor.Decode(&referSchema)

		if err == nil {
			var eachElement T
			copier.Copy(&eachElement, &referSchema)
			elementDetails = append(elementDetails, eachElement)
		} else {
			log.Println(err)
			return elementDetails, err
		}
	}
	return elementDetails, nil
}

func (comb Combination[T]) FindOne(filter interface{}, opts ...*options.FindOneOptions) (T, error) {
	accessCollection := mongodb.Collection(comb.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var element T
	queryError := accessCollection.FindOne(ctx, filter, opts...).Decode(&element)
	defer cancel()
	if queryError != nil {
		log.Println(queryError)
		return element, queryError
	}
	return element, nil
}

func (comb Combination[T]) UpdateOne(filter interface{}, update interface{}) (*mongo.UpdateResult, error){
	accessCollection := mongodb.Collection(comb.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	updateResult, err := accessCollection.UpdateOne(ctx, filter, update)
	defer cancel()
	if err != nil {
		log.Println(err)
	}
	return updateResult, err
}

func (comb Combination[T]) UpdateMany(filter interface{}, update interface{}) (*mongo.UpdateResult, error){
	accessCollection := mongodb.Collection(comb.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	updateResults, err := accessCollection.UpdateMany(ctx, filter, update)
	defer cancel()
	if err != nil {
		log.Println(err)
	}
	return updateResults, err
}

func (comb Combination[T]) InsertOne(data T) (*mongo.InsertOneResult, error) {
	accessCollection := mongodb.Collection(comb.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	insertResult, err := accessCollection.InsertOne(ctx, data)
	defer cancel()
	if err != nil {
		log.Println(err)
	}
	return insertResult, err
}

func (comb Combination[T]) InsertMany(data []interface{}) (*mongo.InsertManyResult, error) {
	accessCollection := mongodb.Collection(comb.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	insertResults, err := accessCollection.InsertMany(ctx, data)
	defer cancel()
	if err != nil {
		log.Println(err)
	}
	return insertResults, err
}

func (comb Combination[T]) DeleteOne(filter interface{}) (*mongo.DeleteResult, error) {
	accessCollection := mongodb.Collection(comb.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	deleteResults, err := accessCollection.DeleteOne(ctx, filter)
	defer cancel()
	if err != nil {
		log.Println(err)
	}
	return deleteResults, err
}

func (comb Combination[T]) DeleteMany(filter interface{}) (*mongo.DeleteResult, error) {
	accessCollection := mongodb.Collection(comb.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	deleteResults, err := accessCollection.DeleteMany(ctx, filter)
	defer cancel()
	if err != nil {
		log.Println(err)
	}
	return deleteResults, err
}

func (comb Combination[T]) ReplaceOne(filter interface{}, update interface{}) (*mongo.UpdateResult, error){
	accessCollection := mongodb.Collection(comb.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	replaceResult, err := accessCollection.ReplaceOne(ctx, filter, update)
	defer cancel()
	if err != nil {
		log.Println(err)
	}
	return replaceResult, err
}

func (comb Combination[T]) Drop() error {
	accessCollection := mongodb.Collection(comb.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err := accessCollection.Drop(ctx)
	defer cancel()
	if err != nil {
		log.Println(err)
	}
	return err
}

func (comb Combination[T]) Count(filter interface{}) (int64, error) {
	accessCollection := mongodb.Collection(comb.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	count, err := accessCollection.CountDocuments(ctx, filter)
	defer cancel()
	if err != nil {
		log.Println(err)
	}
	return count, err
}
///////////////////////////////////////////////////////////////////////////
// Connection URI

func Connect(uri, dbName string) {
	log.Println("Database connecting...")
	// Create a new client and connect to the server
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	mongodb = client.Database(dbName)
	log.Println("Database connected")
}

// New

func New[T any](collection string) Combination[T] {
	model := Combination[T]{Collection: collection}
	return model
}
