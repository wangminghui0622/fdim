package database

import (
	"context"
	"time"

	"fdim/pkg/model"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FavoriteMsgMongo struct {
	coll *mongo.Collection
}

func NewFavoriteMsgMongo(db *mongo.Database) *FavoriteMsgMongo {
	coll := db.Collection("favorite_msg")
	// 创建索引
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "create_time", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "client_msg_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	coll.Indexes().CreateMany(ctx, indexes)

	return &FavoriteMsgMongo{coll: coll}
}

func (f *FavoriteMsgMongo) Create(msg *model.FavoriteMsg) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}
	if msg.CreateTime.IsZero() {
		msg.CreateTime = time.Now()
	}

	_, err := f.coll.InsertOne(ctx, msg)
	return err
}

func (f *FavoriteMsgMongo) Delete(userID, favoriteID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := f.coll.DeleteOne(ctx, bson.M{"_id": favoriteID, "user_id": userID})
	return err
}

func (f *FavoriteMsgMongo) DeleteByClientMsgID(userID, clientMsgID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := f.coll.DeleteOne(ctx, bson.M{"user_id": userID, "client_msg_id": clientMsgID})
	return err
}

func (f *FavoriteMsgMongo) GetByUserID(userID string, pageNumber, showNumber int32) ([]*model.FavoriteMsg, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}

	// 获取总数
	total, err := f.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	opts := options.Find().
		SetSort(bson.D{{Key: "create_time", Value: -1}}).
		SetSkip(int64((pageNumber - 1) * showNumber)).
		SetLimit(int64(showNumber))

	cursor, err := f.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []*model.FavoriteMsg
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (f *FavoriteMsgMongo) GetByID(userID, favoriteID string) (*model.FavoriteMsg, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result model.FavoriteMsg
	err := f.coll.FindOne(ctx, bson.M{"_id": favoriteID, "user_id": userID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
