package database

import (
	"context"

	"fdim/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlackDatabase struct {
	collection *mongo.Collection
}

func NewBlackDatabase(db *MongoDB) *BlackDatabase {
	return &BlackDatabase{
		collection: db.GetCollection("blacks"),
	}
}

// AddBlack 添加黑名�?
func (d *BlackDatabase) AddBlack(ctx context.Context, black *model.Black) error {
	_, err := d.collection.InsertOne(ctx, black)
	return err
}

// RemoveBlack 移除黑名�?
func (d *BlackDatabase) RemoveBlack(ctx context.Context, ownerUserID, blockUserID string) error {
	filter := bson.M{
		"owner_user_id": ownerUserID,
		"block_user_id": blockUserID,
	}
	_, err := d.collection.DeleteOne(ctx, filter)
	return err
}

// FindBlackInfos 查找黑名单信�?
func (d *BlackDatabase) FindBlackInfos(ctx context.Context, ownerUserID string, blockUserIDs []string) ([]*model.Black, error) {
	filter := bson.M{
		"owner_user_id": ownerUserID,
		"block_user_id": bson.M{"$in": blockUserIDs},
	}
	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var blacks []*model.Black
	if err := cursor.All(ctx, &blacks); err != nil {
		return nil, err
	}
	return blacks, nil
}

// IsBlack 检查是否在黑名�?
func (d *BlackDatabase) IsBlack(ctx context.Context, ownerUserID, blockUserID string) (bool, error) {
	count, err := d.collection.CountDocuments(ctx, bson.M{
		"owner_user_id": ownerUserID,
		"block_user_id": blockUserID,
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// PageOwnerBlacks 分页获取黑名�?
func (d *BlackDatabase) PageOwnerBlacks(ctx context.Context, ownerUserID string, offset, limit int32) (int64, []*model.Black, error) {
	filter := bson.M{"owner_user_id": ownerUserID}

	// 获取总数
	total, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// 分页查询
	opts := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "create_time", Value: -1}})
	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var blacks []*model.Black
	if err := cursor.All(ctx, &blacks); err != nil {
		return 0, nil, err
	}
	return total, blacks, nil
}
