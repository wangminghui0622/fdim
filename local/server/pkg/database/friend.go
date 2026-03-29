package database

import (
	"context"
	"time"

	"fdim/pkg/model"
	"fdim/protocol/constant"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FriendDatabase struct {
	friendCollection  *mongo.Collection
	requestCollection *mongo.Collection
}

func NewFriendDatabase(db *MongoDB) *FriendDatabase {
	return &FriendDatabase{
		friendCollection:  db.GetCollection("friends"),
		requestCollection: db.GetCollection("friend_requests"),
	}
}

// FindFriends Һ
func (d *FriendDatabase) FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) ([]*model.Friend, error) {
	filter := bson.M{
		"owner_user_id":  ownerUserID,
		"friend_user_id": bson.M{"$in": friendUserIDs},
	}
	cursor, err := d.friendCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var friends []*model.Friend
	if err := cursor.All(ctx, &friends); err != nil {
		return nil, err
	}
	return friends, nil
}

// FindFriendsWithError Һѣ飩
func (d *FriendDatabase) FindFriendsWithError(ctx context.Context, ownerUserID string, friendUserIDs []string) ([]*model.Friend, error) {
	return d.FindFriends(ctx, ownerUserID, friendUserIDs)
}

// CheckIn ǷΪ
func (d *FriendDatabase) CheckIn(ctx context.Context, userID1, userID2 string) (bool, bool, error) {
	count1, err := d.friendCollection.CountDocuments(ctx, bson.M{
		"owner_user_id":  userID1,
		"friend_user_id": userID2,
	})
	if err != nil {
		return false, false, err
	}

	count2, err := d.friendCollection.CountDocuments(ctx, bson.M{
		"owner_user_id":  userID2,
		"friend_user_id": userID1,
	})
	if err != nil {
		return false, false, err
	}

	return count1 > 0, count2 > 0, nil
}

// AddFriendRequest Ӻ루ݵȣͬһ from/to ֻһظ£
func (d *FriendDatabase) AddFriendRequest(ctx context.Context, fromUserID, toUserID, reqMsg, ex string) error {
	filter := bson.M{
		"from_user_id": fromUserID,
		"to_user_id":   toUserID,
	}
	update := bson.M{
		"$set": bson.M{
			"req_msg":       reqMsg,
			"ex":            ex,
			"create_time":   time.Now(),
			"handle_result": 0,
			"handler_user_id": "",
			"handle_msg":    "",
			"handle_time":   time.Time{},
		},
		"$setOnInsert": bson.M{
			"from_user_id": fromUserID,
			"to_user_id":   toUserID,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := d.requestCollection.UpdateOne(ctx, filter, update, opts)
	return err
}

// BecomeFriends Ϊ
func (d *FriendDatabase) BecomeFriends(ctx context.Context, ownerUserID string, friendUserIDs []string, addSource int32) error {
	now := time.Now()
	docs := make([]interface{}, len(friendUserIDs))
	for i, friendUserID := range friendUserIDs {
		docs[i] = model.Friend{
			OwnerUserID:  ownerUserID,
			FriendUserID: friendUserID,
			CreateTime:   now,
			AddSource:    addSource,
		}
	}
	_, err := d.friendCollection.InsertMany(ctx, docs)
	return err
}

// AgreeFriendRequest ͬ
func (d *FriendDatabase) AgreeFriendRequest(ctx context.Context, req *model.FriendRequest) error {
	// º״??
	filter := bson.M{
		"from_user_id": req.FromUserID,
		"to_user_id":   req.ToUserID,
	}
	update := bson.M{
		"$set": bson.M{
			"handle_result":   req.HandleResult,
			"handle_msg":      req.HandleMsg,
			"handle_time":     req.HandleTime,
			"handler_user_id": req.HandlerUserID,
		},
	}
	_, err := d.requestCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	// ˫ѹϵ
	now := time.Now()
	friends := []interface{}{
		model.Friend{
			ID:           primitive.NewObjectID(),
			OwnerUserID:  req.FromUserID,
			FriendUserID: req.ToUserID,
			CreateTime:   now,
			AddSource:    constant.BecomeFriendByApply, // ͨӺ
		},
		model.Friend{
			ID:           primitive.NewObjectID(),
			OwnerUserID:  req.ToUserID,
			FriendUserID: req.FromUserID,
			CreateTime:   now,
			AddSource:    constant.BecomeFriendByApply, // ͨӺ
		},
	}
	_, err = d.friendCollection.InsertMany(ctx, friends)
	return err
}

// RefuseFriendRequest ܾ
func (d *FriendDatabase) RefuseFriendRequest(ctx context.Context, req *model.FriendRequest) error {
	filter := bson.M{
		"from_user_id": req.FromUserID,
		"to_user_id":   req.ToUserID,
	}
	update := bson.M{
		"$set": bson.M{
			"handle_result":   req.HandleResult,
			"handle_msg":      req.HandleMsg,
			"handle_time":     req.HandleTime,
			"handler_user_id": req.HandlerUserID,
		},
	}
	_, err := d.requestCollection.UpdateOne(ctx, filter, update)
	return err
}

// Delete ɾ
func (d *FriendDatabase) Delete(ctx context.Context, ownerUserID string, friendUserIDs []string) error {
	filter := bson.M{
		"owner_user_id":  ownerUserID,
		"friend_user_id": bson.M{"$in": friendUserIDs},
	}
	_, err := d.friendCollection.DeleteMany(ctx, filter)
	return err
}

// UpdateRemark ºѱע
func (d *FriendDatabase) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) error {
	filter := bson.M{
		"owner_user_id":  ownerUserID,
		"friend_user_id": friendUserID,
	}
	update := bson.M{
		"$set": bson.M{
			"remark": remark,
		},
	}
	_, err := d.friendCollection.UpdateOne(ctx, filter, update)
	return err
}

// UpdateFriends ºϢ
func (d *FriendDatabase) UpdateFriends(ctx context.Context, ownerUserID string, friendUserIDs []string, data map[string]interface{}) error {
	filter := bson.M{
		"owner_user_id":  ownerUserID,
		"friend_user_id": bson.M{"$in": friendUserIDs},
	}
	update := bson.M{
		"$set": data,
	}
	_, err := d.friendCollection.UpdateMany(ctx, filter, update)
	return err
}

// FindFriendUserIDs ҺûIDб
func (d *FriendDatabase) FindFriendUserIDs(ctx context.Context, ownerUserID string) ([]string, error) {
	filter := bson.M{
		"owner_user_id": ownerUserID,
	}
	cursor, err := d.friendCollection.Find(ctx, filter, options.Find().SetProjection(bson.M{"friend_user_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var friendIDs []string
	for cursor.Next(ctx) {
		var result struct {
			FriendUserID string `bson:"friend_user_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		friendIDs = append(friendIDs, result.FriendUserID)
	}
	return friendIDs, nil
}

// PageFriendRequestToMe ҳȡ͸ҵĺ
func (d *FriendDatabase) PageFriendRequestToMe(ctx context.Context, toUserID string, handleResults []int32, offset, limit int32) (int64, []*model.FriendRequest, error) {
	filter := bson.M{"to_user_id": toUserID}
	if len(handleResults) > 0 {
		filter["handle_result"] = bson.M{"$in": handleResults}
	}

	// ȡ
	total, err := d.requestCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "create_time", Value: -1}})
	cursor, err := d.requestCollection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var requests []*model.FriendRequest
	if err := cursor.All(ctx, &requests); err != nil {
		return 0, nil, err
	}
	return total, requests, nil
}

// PageFriendRequestFromMe ҳȡҷ͵ĺ
func (d *FriendDatabase) PageFriendRequestFromMe(ctx context.Context, fromUserID string, handleResults []int32, offset, limit int32) (int64, []*model.FriendRequest, error) {
	filter := bson.M{"from_user_id": fromUserID}
	if len(handleResults) > 0 {
		filter["handle_result"] = bson.M{"$in": handleResults}
	}

	// ȡ
	total, err := d.requestCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "create_time", Value: -1}})
	cursor, err := d.requestCollection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var requests []*model.FriendRequest
	if err := cursor.All(ctx, &requests); err != nil {
		return 0, nil, err
	}
	return total, requests, nil
}

// GetUnhandledFriendRequestCount ȡδĺ
func (d *FriendDatabase) GetUnhandledFriendRequestCount(ctx context.Context, userID string, ts int64) (int64, error) {
	filter := bson.M{
		"to_user_id":    userID,
		"handle_result": 0, // δ??
	}
	if ts != 0 {
		filter["create_time"] = bson.M{"$gt": time.UnixMilli(ts)}
	}
	return d.requestCollection.CountDocuments(ctx, filter)
}

// GetSelfUnhandledFriendRequestCount ȡԼ͵δ??
func (d *FriendDatabase) GetSelfUnhandledFriendRequestCount(ctx context.Context, userID string, ts int64) (int64, error) {
	filter := bson.M{
		"from_user_id":  userID,
		"handle_result": 0, // δ??
	}
	if ts != 0 {
		filter["create_time"] = bson.M{"$gt": time.UnixMilli(ts)}
	}
	return d.requestCollection.CountDocuments(ctx, filter)
}

// PageOwnerFriends ҳȡб
func (d *FriendDatabase) PageOwnerFriends(ctx context.Context, ownerUserID string, offset, limit int32) (int64, []*model.Friend, error) {
	filter := bson.M{"owner_user_id": ownerUserID}

	// ȡ
	total, err := d.friendCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ is_pinned Ȼ _id ??
	opts := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{
			{Key: "is_pinned", Value: -1},
			{Key: "_id", Value: 1},
		})
	cursor, err := d.friendCollection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var friends []*model.Friend
	if err := cursor.All(ctx, &friends); err != nil {
		return 0, nil, err
	}
	return total, friends, nil
}

// FindBothFriendRequests ˫루??GetDesignatedFriendsApply??
func (d *FriendDatabase) FindBothFriendRequests(ctx context.Context, fromUserID, toUserID string) ([]*model.FriendRequest, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"from_user_id": fromUserID, "to_user_id": toUserID},
			{"from_user_id": toUserID, "to_user_id": fromUserID},
		},
	}
	cursor, err := d.requestCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var requests []*model.FriendRequest
	if err := cursor.All(ctx, &requests); err != nil {
		return nil, err
	}
	return requests, nil
}
