package database

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoDBFromClient ʹеMongoDBͻ˴MongoDBʵ
func NewMongoDBFromClient(client *mongo.Client, databaseName string) *MongoDB {
	return &MongoDB{
		client:   client,
		database: client.Database(databaseName),
	}
}

func NewMongoDB(host, database, username, password string) *MongoDB {
	var uri string

	// ж host ǷѾ MongoDB URIƾݻѯ
	isFullURI := strings.Contains(host, "@") || strings.Contains(host, "?")

	if isFullURI {
		// host Ѿ URI mongodb://root:test625@127.0.0.1:27017/FDIM_v3?authSource=adminֱʹ
		uri = host
	} else if username != "" && password != "" {
		// ȥܴڵ mongodb:// ǰ׺ظƴ
		cleanHost := strings.TrimPrefix(host, "mongodb://")
		// ֤ URI + authSource=adminroot û admin ⣩
		uri = fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=admin", username, password, cleanHost, database)
	} else {
		cleanHost := strings.TrimPrefix(host, "mongodb://")
		uri = fmt.Sprintf("mongodb://%s/%s", cleanHost, database)
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
	}

	db := client.Database(database)
	return &MongoDB{
		client:   client,
		database: db,
	}
}

func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.database.Collection(name)
}

func (m *MongoDB) GetClient() *mongo.Client {
	return m.client
}

func (m *MongoDB) GetDatabase() *mongo.Database {
	return m.database
}

func (m *MongoDB) Close(ctx context.Context) error {
	if m.client != nil {
		return m.client.Disconnect(ctx)
	}
	return nil
}
