package database

import (
	"context"
	"time"

	"fdim/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserDatabase struct {
	collection *mongo.Collection
}

func NewUserDatabase(db *MongoDB) *UserDatabase {
	return &UserDatabase{
		collection: db.GetCollection("users"),
	}
}

// Find û
func (d *UserDatabase) Find(ctx context.Context, userIDs []string) ([]*model.User, error) {
	filter := bson.M{"user_id": bson.M{"$in": userIDs}}
	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID IDȡû
func (d *UserDatabase) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	err := d.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateByMap Mapû
func (d *UserDatabase) UpdateByMap(ctx context.Context, userID string, data map[string]interface{}) error {
	_, err := d.collection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": data})
	return err
}

// Create û
func (d *UserDatabase) Create(ctx context.Context, users []*model.User) error {
	docs := make([]interface{}, len(users))
	for i, u := range users {
		docs[i] = u
	}
	_, err := d.collection.InsertMany(ctx, docs)
	return err
}

// IsExist ûǷ??
func (d *UserDatabase) IsExist(ctx context.Context, userIDs []string) (bool, error) {
	count, err := d.collection.CountDocuments(ctx, bson.M{"user_id": bson.M{"$in": userIDs}})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindWithError û飩
func (d *UserDatabase) FindWithError(ctx context.Context, userIDs []string) ([]*model.User, error) {
	return d.Find(ctx, userIDs)
}

// PageFindUser ҳû??app_manger_level Χ??
func (d *UserDatabase) PageFindUser(ctx context.Context, level1, level2 int32, offset, limit int32) (int64, []*model.User, error) {
	filter := bson.M{
		"app_manger_level": bson.M{"$gte": level1, "$lte": level2},
	}

	// ȡ
	total, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "create_time", Value: -1}})
	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		return 0, nil, err
	}
	return total, users, nil
}

// PageFindUserWithKeyword ݹؼʷҳ??
func (d *UserDatabase) PageFindUserWithKeyword(ctx context.Context, level1, level2 int32, userID, nickname string, offset, limit int32) (int64, []*model.User, error) {
	filter := bson.M{
		"app_manger_level": bson.M{"$gte": level1, "$lte": level2},
	}
	if userID != "" {
		filter["user_id"] = bson.M{"$regex": userID, "$options": "i"}
	}
	if nickname != "" {
		filter["nickname"] = bson.M{"$regex": nickname, "$options": "i"}
	}

	// ȡ
	total, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "create_time", Value: -1}})
	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		return 0, nil, err
	}
	return total, users, nil
}

// SearchUsers ݹؼûģƥuserIDnickname
func (d *UserDatabase) SearchUsers(ctx context.Context, keyword string, offset, limit int32) (int64, []*model.User, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"user_id": bson.M{"$regex": keyword, "$options": "i"}},
			{"nickname": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}

	// ȡ
	total, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "create_time", Value: -1}})
	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		return 0, nil, err
	}
	return total, users, nil
}

// CountUsers ͳû
func (d *UserDatabase) CountUsers(ctx context.Context) (int64, error) {
	return d.collection.CountDocuments(ctx, bson.M{})
}

// CountUsersBeforeTime ͳָʱ֮ǰû
func (d *UserDatabase) CountUsersBeforeTime(ctx context.Context, timestamp int64) (int64, error) {
	filter := bson.M{
		"create_time": bson.M{"$lt": time.UnixMilli(timestamp)},
	}
	return d.collection.CountDocuments(ctx, filter)
}

// GetAllUserID ҳȡûID
func (d *UserDatabase) GetAllUserID(ctx context.Context, offset, limit int32) (int64, []string, error) {
	// ȡ
	total, err := d.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, nil, err
	}

	// ҳѯֻȡ user_id ֶ
	opts := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetProjection(bson.M{"user_id": 1}).
		SetSort(bson.D{{Key: "create_time", Value: -1}})
	cursor, err := d.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var userIDs []string
	for cursor.Next(ctx) {
		var result struct {
			UserID string `bson:"user_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return 0, nil, err
		}
		userIDs = append(userIDs, result.UserID)
	}

	return total, userIDs, nil
}

// Take ȡû飩
func (d *UserDatabase) Take(ctx context.Context, userID string) (*model.User, error) {
	return d.GetUserByID(ctx, userID)
}
