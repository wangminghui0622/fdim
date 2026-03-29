package database

import (
	"context"

	"fdim/protocol/sdkws"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BotDatabase Botݿӿ
type BotDatabase interface {
	// Agent
	CreateAgent(ctx context.Context, agent *Agent) error
	UpdateAgent(ctx context.Context, userID string, update map[string]interface{}) error
	PageFindAgent(ctx context.Context, userIDs []string, pagination *sdkws.RequestPagination) (int64, []*Agent, error)
	DeleteAgent(ctx context.Context, userIDs []string) error
}

// Agent Botģ
type Agent struct {
	UserID     string `bson:"user_id"`
	Nickname   string `bson:"nickname"`
	FaceURL    string `bson:"face_url"`
	URL        string `bson:"url"`
	Key        string `bson:"key"`
	Identity   string `bson:"identity"`
	Model      string `bson:"model"`
	Prompts    string `bson:"prompts"`
	CreateTime int64  `bson:"create_time"`
}

// NewBotDatabase Botݿʵ
func NewBotDatabase(mongoDB *MongoDB) BotDatabase {
	return &botDatabase{
		collection: mongoDB.GetCollection("bot_agent"),
	}
}

type botDatabase struct {
	collection *mongo.Collection
}

func (d *botDatabase) CreateAgent(ctx context.Context, agent *Agent) error {
	_, err := d.collection.InsertOne(ctx, agent)
	return err
}

func (d *botDatabase) UpdateAgent(ctx context.Context, userID string, update map[string]interface{}) error {
	_, err := d.collection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": update})
	return err
}

func (d *botDatabase) PageFindAgent(ctx context.Context, userIDs []string, pagination *sdkws.RequestPagination) (int64, []*Agent, error) {
	// ѯ
	filter := bson.M{}
	if len(userIDs) > 0 {
		filter["user_id"] = bson.M{"$in": userIDs}
	}

	// ҳ
	offset := int64(0)
	limit := int64(10)
	if pagination != nil {
		if pagination.PageNumber > 0 {
			offset = int64((pagination.PageNumber - 1) * pagination.ShowNumber)
		}
		if pagination.ShowNumber > 0 {
			limit = int64(pagination.ShowNumber)
		}
	}

	// ȡ
	total, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "create_time", Value: -1}})

	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var results []*Agent
	if err := cursor.All(ctx, &results); err != nil {
		return 0, nil, err
	}

	return total, results, nil
}

func (d *botDatabase) DeleteAgent(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	_, err := d.collection.DeleteMany(ctx, bson.M{"user_id": bson.M{"$in": userIDs}})
	return err
}
