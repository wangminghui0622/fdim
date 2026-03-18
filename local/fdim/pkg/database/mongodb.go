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

// NewMongoDBFromClient 使用现有的MongoDB客户端创建MongoDB实例
func NewMongoDBFromClient(client *mongo.Client, databaseName string) *MongoDB {
	return &MongoDB{
		client:   client,
		database: client.Database(databaseName),
	}
}

func NewMongoDB(host, database, username, password string) *MongoDB {
	var uri string

	// 判断 host 是否已经是完整的 MongoDB URI（含凭据或查询参数）
	isFullURI := strings.Contains(host, "@") || strings.Contains(host, "?")

	if isFullURI {
		// host 已经是完整 URI（如 mongodb://root:test625@127.0.0.1:27017/openim_v3?authSource=admin），直接使用
		uri = host
	} else if username != "" && password != "" {
		// 去掉可能存在的 mongodb:// 前缀，避免重复拼接
		cleanHost := strings.TrimPrefix(host, "mongodb://")
		// 认证 URI + authSource=admin（root 用户在 admin 库）
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
