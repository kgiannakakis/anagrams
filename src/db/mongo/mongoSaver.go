package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Saver provides methods for writing data to a MongoDb Collection. It currently supports
// only connecting to a single connection at a time
type Saver struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// Config wraps options for connecting to a mongo db collection
type Config struct {
	URI        string
	Database   string
	Collection string
}

// Connect connects to a MongoDb Collection.
func (m *Saver) Connect(config interface{}) (err error) {

	cfg := config.(Config)

	clientOptions := options.Client().ApplyURI(cfg.URI)

	m.client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return
	}

	err = m.client.Ping(context.TODO(), nil)

	m.collection = m.client.Database(cfg.Database).Collection(cfg.Collection)

	return
}

// Disconnect disconnects from the MongoDb collection
func (m *Saver) Disconnect() (err error) {
	err = m.client.Disconnect(context.TODO())
	return
}

// Insert inserts a single item to the open collection
func (m *Saver) Insert(item interface{}) (result interface{}, err error) {
	result, err = m.collection.InsertOne(context.TODO(), item)
	return
}

// InsertMany inserts many items to the open collection at the same time
func (m *Saver) InsertMany(items []interface{}) (result interface{}, err error) {
	result, err = m.collection.InsertMany(context.TODO(), items)
	return
}
