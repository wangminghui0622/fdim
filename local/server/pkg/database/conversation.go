package database

import (
	"context"
	"time"

	"fdim/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConversationDatabase struct {
	collection *mongo.Collection
}

func NewConversationDatabase(db *MongoDB) *ConversationDatabase {
	coll := db.GetCollection("conversations")
	// 
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "owner_user_id", Value: 1},
				{Key: "conversation_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "conversation_id", Value: 1},
			},
		},
	}
	coll.Indexes().CreateMany(context.Background(), indexes)

	return &ConversationDatabase{
		collection: coll,
	}
}

// Create Ự
func (d *ConversationDatabase) Create(ctx context.Context, conversations []*model.Conversation) error {
	if len(conversations) == 0 {
		return nil
	}
	docs := make([]interface{}, len(conversations))
	for i, conv := range conversations {
		docs[i] = conv
	}
	_, err := d.collection.InsertMany(ctx, docs)
	return err
}

// UpdateByMap Map»Ự
func (d *ConversationDatabase) UpdateByMap(ctx context.Context, userIDs []string, conversationID string, args map[string]interface{}) (int64, error) {
	if len(args) == 0 || len(userIDs) == 0 {
		return 0, nil
	}
	filter := bson.M{
		"conversation_id": conversationID,
		"owner_user_id":   bson.M{"$in": userIDs},
	}
	result, err := d.collection.UpdateMany(ctx, filter, bson.M{"$set": args})
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

// UpdateUserConversations ûлỰ
func (d *ConversationDatabase) UpdateUserConversations(ctx context.Context, userID string, args map[string]interface{}) ([]*model.Conversation, error) {
	if len(args) == 0 {
		return nil, nil
	}
	filter := bson.M{"user_id": userID}

	// ȲҪµĻỰ
	cursor, err := d.collection.Find(ctx, filter, options.Find().SetProjection(bson.M{"_id": 0, "owner_user_id": 1, "conversation_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var conversations []*model.Conversation
	if err := cursor.All(ctx, &conversations); err != nil {
		return nil, err
	}

	// 
	_, err = d.collection.UpdateMany(ctx, filter, bson.M{"$set": args})
	if err != nil {
		return nil, err
	}

	return conversations, nil
}

// Update »Ự
func (d *ConversationDatabase) Update(ctx context.Context, conversation *model.Conversation) error {
	filter := bson.M{
		"owner_user_id":   conversation.OwnerUserID,
		"conversation_id": conversation.ConversationID,
	}
	_, err := d.collection.ReplaceOne(ctx, filter, conversation)
	return err
}

// Find һỰ
func (d *ConversationDatabase) Find(ctx context.Context, ownerUserID string, conversationIDs []string) ([]*model.Conversation, error) {
	// conversationIDsΪգؿս
	if len(conversationIDs) == 0 {
		return []*model.Conversation{}, nil
	}
	filter := bson.M{
		"owner_user_id":   ownerUserID,
		"conversation_id": bson.M{"$in": conversationIDs},
	}
	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var conversations []*model.Conversation
	if err := cursor.All(ctx, &conversations); err != nil {
		return nil, err
	}
	return conversations, nil
}

// FindUserIDAllConversationID ûлỰID
func (d *ConversationDatabase) FindUserIDAllConversationID(ctx context.Context, userID string) ([]string, error) {
	filter := bson.M{"owner_user_id": userID}
	opts := options.Find().SetProjection(bson.M{"_id": 0, "conversation_id": 1})
	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var conversationIDs []string
	for cursor.Next(ctx) {
		var result struct {
			ConversationID string `bson:"conversation_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		conversationIDs = append(conversationIDs, result.ConversationID)
	}
	return conversationIDs, nil
}

// FindUserIDAllNotNotifyConversationID ûв֪ͨỰID
func (d *ConversationDatabase) FindUserIDAllNotNotifyConversationID(ctx context.Context, userID string) ([]string, error) {
	filter := bson.M{
		"owner_user_id": userID,
		"recv_msg_opt":  1, // ReceiveNotNotifyMessage
	}
	opts := options.Find().SetProjection(bson.M{"_id": 0, "conversation_id": 1})
	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var conversationIDs []string
	for cursor.Next(ctx) {
		var result struct {
			ConversationID string `bson:"conversation_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		conversationIDs = append(conversationIDs, result.ConversationID)
	}
	return conversationIDs, nil
}

// FindUserIDAllPinnedConversationID ûöỰID
func (d *ConversationDatabase) FindUserIDAllPinnedConversationID(ctx context.Context, userID string) ([]string, error) {
	filter := bson.M{
		"owner_user_id": userID,
		"is_pinned":     true,
	}
	opts := options.Find().SetProjection(bson.M{"_id": 0, "conversation_id": 1})
	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var conversationIDs []string
	for cursor.Next(ctx) {
		var result struct {
			ConversationID string `bson:"conversation_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		conversationIDs = append(conversationIDs, result.ConversationID)
	}
	return conversationIDs, nil
}

// Take ȡỰ
func (d *ConversationDatabase) Take(ctx context.Context, userID, conversationID string) (*model.Conversation, error) {
	filter := bson.M{
		"owner_user_id":   userID,
		"conversation_id": conversationID,
	}
	var conversation model.Conversation
	err := d.collection.FindOne(ctx, filter).Decode(&conversation)
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

// FindUserIDAllConversations ûлỰ
func (d *ConversationDatabase) FindUserIDAllConversations(ctx context.Context, userID string) ([]*model.Conversation, error) {
	filter := bson.M{"owner_user_id": userID}
	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var conversations []*model.Conversation
	if err := cursor.All(ctx, &conversations); err != nil {
		return nil, err
	}
	return conversations, nil
}

// GetConversationsByConversationID ݻỰIDһỰ
func (d *ConversationDatabase) GetConversationsByConversationID(ctx context.Context, conversationIDs []string) ([]*model.Conversation, error) {
	// conversationIDsΪջnilؿս
	if len(conversationIDs) == 0 {
		return []*model.Conversation{}, nil
	}
	filter := bson.M{"conversation_id": bson.M{"$in": conversationIDs}}
	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var conversations []*model.Conversation
	if err := cursor.All(ctx, &conversations); err != nil {
		return nil, err
	}
	return conversations, nil
}

// DeleteUsersConversations ɾûĻỰ
func (d *ConversationDatabase) DeleteUsersConversations(ctx context.Context, userID string, conversationIDs []string) error {
	if len(conversationIDs) == 0 {
		return nil
	}
	filter := bson.M{
		"owner_user_id":   userID,
		"conversation_id": bson.M{"$in": conversationIDs},
	}
	_, err := d.collection.DeleteMany(ctx, filter)
	return err
}

// GetAllConversationIDs ȡлỰID
func (d *ConversationDatabase) GetAllConversationIDs(ctx context.Context) ([]string, error) {
	pipeline := []bson.M{
		{"$group": bson.M{"_id": "$conversation_id"}},
		{"$project": bson.M{"_id": 0, "conversation_id": "$_id"}},
	}
	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var conversationIDs []string
	for cursor.Next(ctx) {
		var result struct {
			ConversationID string `bson:"conversation_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		conversationIDs = append(conversationIDs, result.ConversationID)
	}
	return conversationIDs, nil
}

// FindUserID ûID
func (d *ConversationDatabase) FindUserID(ctx context.Context, userIDs []string, conversationIDs []string) ([]string, error) {
	if len(userIDs) == 0 || len(conversationIDs) == 0 {
		return []string{}, nil
	}
	filter := bson.M{
		"owner_user_id":   bson.M{"$in": userIDs},
		"conversation_id": bson.M{"$in": conversationIDs},
	}
	opts := options.Find().SetProjection(bson.M{"_id": 0, "owner_user_id": 1})
	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var resultUserIDs []string
	userIDMap := make(map[string]bool)
	for cursor.Next(ctx) {
		var result struct {
			OwnerUserID string `bson:"owner_user_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		if !userIDMap[result.OwnerUserID] {
			resultUserIDs = append(resultUserIDs, result.OwnerUserID)
			userIDMap[result.OwnerUserID] = true
		}
	}
	return resultUserIDs, nil
}

// FindConversationID һỰID
func (d *ConversationDatabase) FindConversationID(ctx context.Context, userID string, conversationIDs []string) ([]string, error) {
	if len(conversationIDs) == 0 {
		return []string{}, nil
	}
	filter := bson.M{
		"owner_user_id":   userID,
		"conversation_id": bson.M{"$in": conversationIDs},
	}
	opts := options.Find().SetProjection(bson.M{"_id": 0, "conversation_id": 1})
	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var resultConversationIDs []string
	for cursor.Next(ctx) {
		var result struct {
			ConversationID string `bson:"conversation_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		resultConversationIDs = append(resultConversationIDs, result.ConversationID)
	}
	return resultConversationIDs, nil
}

// FindRecvMsgUserIDs ҽϢûID
func (d *ConversationDatabase) FindRecvMsgUserIDs(ctx context.Context, conversationID string, recvOpts []int32) ([]string, error) {
	filter := bson.M{"conversation_id": conversationID}
	if len(recvOpts) > 0 {
		filter["recv_msg_opt"] = bson.M{"$in": recvOpts}
	}
	opts := options.Find().SetProjection(bson.M{"_id": 0, "owner_user_id": 1})
	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userIDs []string
	for cursor.Next(ctx) {
		var result struct {
			OwnerUserID string `bson:"owner_user_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, result.OwnerUserID)
	}
	return userIDs, nil
}

// GetUserRecvMsgOpt ȡûϢѡ
func (d *ConversationDatabase) GetUserRecvMsgOpt(ctx context.Context, ownerUserID, conversationID string) (int32, error) {
	filter := bson.M{
		"owner_user_id":   ownerUserID,
		"conversation_id": conversationID,
	}
	opts := options.FindOne().SetProjection(bson.M{"_id": 0, "recv_msg_opt": 1})
	var result struct {
		RecvMsgOpt int32 `bson:"recv_msg_opt"`
	}
	err := d.collection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		return 0, err
	}
	return result.RecvMsgOpt, nil
}

// GetAllConversationIDsNumber ȡлỰID
func (d *ConversationDatabase) GetAllConversationIDsNumber(ctx context.Context) (int64, error) {
	pipeline := []bson.M{
		{"$group": bson.M{"_id": "$conversation_id"}},
		{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
		{"$project": bson.M{"_id": 0}},
	}
	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result struct {
			Count int64 `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.Count, nil
	}
	return 0, nil
}

// PageConversationIDs ҳȡỰID
func (d *ConversationDatabase) PageConversationIDs(ctx context.Context, offset, limit int64) ([]string, error) {
	pipeline := []bson.M{
		{"$group": bson.M{"_id": "$conversation_id"}},
		{"$project": bson.M{"_id": 0, "conversation_id": "$_id"}},
		{"$skip": offset},
		{"$limit": limit},
	}
	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var conversationIDs []string
	for cursor.Next(ctx) {
		var result struct {
			ConversationID string `bson:"conversation_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		conversationIDs = append(conversationIDs, result.ConversationID)
	}
	return conversationIDs, nil
}

// GetConversationIDsNeedDestruct ȡҪٵĻỰID
func (d *ConversationDatabase) GetConversationIDsNeedDestruct(ctx context.Context) ([]*model.Conversation, error) {
	// is_msg_destruct = 1 && msg_destruct_time != 0 && (UNIX_TIMESTAMP(NOW()) > (msg_destruct_time + UNIX_TIMESTAMP(latest_msg_destruct_time)) || latest_msg_destruct_time is NULL)
	now := time.Now()
	filter := bson.M{
		"is_msg_destruct":   true,
		"msg_destruct_time": bson.M{"$ne": 0},
		"$or": []bson.M{
			{
				"$expr": bson.M{
					"$gt": []interface{}{
						now,
						bson.M{"$add": []interface{}{"$msg_destruct_time", "$latest_msg_destruct_time"}},
					},
				},
			},
			{
				"latest_msg_destruct_time": nil,
			},
		},
	}
	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var conversations []*model.Conversation
	if err := cursor.All(ctx, &conversations); err != nil {
		return nil, err
	}
	return conversations, nil
}

// GetConversationNotReceiveMessageUserIDs ȡϢûID
func (d *ConversationDatabase) GetConversationNotReceiveMessageUserIDs(ctx context.Context, conversationID string) ([]string, error) {
	// recv_msg_opt != ReceiveMessage (0)
	filter := bson.M{
		"conversation_id": conversationID,
		"recv_msg_opt":    bson.M{"$ne": 0}, // ReceiveMessage = 0
	}
	opts := options.Find().SetProjection(bson.M{"_id": 0, "owner_user_id": 1})
	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userIDs []string
	for cursor.Next(ctx) {
		var result struct {
			OwnerUserID string `bson:"owner_user_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, result.OwnerUserID)
	}
	return userIDs, nil
}
