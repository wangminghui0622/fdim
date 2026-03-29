package database

import (
	"context"

	"fdim/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MsgDatabase Ϣݿӿ
type MsgDatabase interface {
	// BatchInsertChat2DB Ϣݿ
	BatchInsertChat2DB(ctx context.Context, conversationID string, msgs []*model.MsgDoc, currentMaxSeq int64) error
	// GetMessagesBySeq кŻȡϢ
	GetMessagesBySeq(ctx context.Context, conversationID string, seqs []int64) ([]*model.MsgDoc, error)
	// GetMaxSeq ȡк
	GetMaxSeq(ctx context.Context, conversationID string) (int64, error)
	// GetLastMessageSeqByTime ݻỰʱȡʱ֮ǰһϢ seq
	GetLastMessageSeqByTime(ctx context.Context, conversationID string, timestamp int64) (int64, error)
	// DeleteMessagesBySeq ݻỰ seq бɾϢ
	DeleteMessagesBySeq(ctx context.Context, conversationID string, seqs []int64) error
	// DeleteMessagesByTimeBefore ɾĳʱ֮ǰϢ
	DeleteMessagesByTimeBefore(ctx context.Context, conversationIDs []string, timestamp int64) error
	// RevokeMsg Ϣ content_type  content
	RevokeMsg(ctx context.Context, conversationID string, seq int64, contentType int32, content []byte) error
}

// MsgDocDatabase Ϣĵݿʵ
type MsgDocDatabase struct {
	collection *mongo.Collection
}

// NewMsgDocDatabase Ϣĵݿ
func NewMsgDocDatabase(db *MongoDB) *MsgDocDatabase {
	coll := db.GetCollection("stream_msg")

	// 
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "conversation_id", Value: 1},
				{Key: "seq", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "conversation_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "send_time", Value: -1},
			},
		},
	}
	coll.Indexes().CreateMany(context.Background(), indexes)

	return &MsgDocDatabase{
		collection: coll,
	}
}

// BatchInsertChat2DB Ϣݿ
func (d *MsgDocDatabase) BatchInsertChat2DB(ctx context.Context, conversationID string, msgs []*model.MsgDoc, currentMaxSeq int64) error {
	if len(msgs) == 0 {
		return nil
	}

	docs := make([]interface{}, len(msgs))
	for i, msg := range msgs {
		docs[i] = msg
	}

	_, err := d.collection.InsertMany(ctx, docs)
	return err
}

// GetMessagesBySeq кŻȡϢ
func (d *MsgDocDatabase) GetMessagesBySeq(ctx context.Context, conversationID string, seqs []int64) ([]*model.MsgDoc, error) {
	filter := bson.M{
		"conversation_id": conversationID,
		"seq":             bson.M{"$in": seqs},
	}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var msgs []*model.MsgDoc
	if err := cursor.All(ctx, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

// GetMaxSeq ȡк
func (d *MsgDocDatabase) GetMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	filter := bson.M{"conversation_id": conversationID}
	opts := options.FindOne().SetSort(bson.D{{Key: "seq", Value: -1}})

	var msg model.MsgDoc
	err := d.collection.FindOne(ctx, filter, opts).Decode(&msg)
	if err == mongo.ErrNoDocuments {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return msg.Seq, nil
}

// GetLastMessageSeqByTime ݻỰʱȡʱ֮ǰһϢ seq
func (d *MsgDocDatabase) GetLastMessageSeqByTime(ctx context.Context, conversationID string, timestamp int64) (int64, error) {
	filter := bson.M{
		"conversation_id": conversationID,
		"send_time":       bson.M{"$lte": timestamp},
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "send_time", Value: -1}})

	var msg model.MsgDoc
	err := d.collection.FindOne(ctx, filter, opts).Decode(&msg)
	if err == mongo.ErrNoDocuments {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return msg.Seq, nil
}

// DeleteMessagesBySeq ݻỰ seq бɾϢ
func (d *MsgDocDatabase) DeleteMessagesBySeq(ctx context.Context, conversationID string, seqs []int64) error {
	if len(seqs) == 0 {
		return nil
	}
	filter := bson.M{
		"conversation_id": conversationID,
		"seq":             bson.M{"$in": seqs},
	}
	_, err := d.collection.DeleteMany(ctx, filter)
	return err
}

// DeleteMessagesByTimeBefore ɾĳʱ֮ǰϢ
func (d *MsgDocDatabase) DeleteMessagesByTimeBefore(ctx context.Context, conversationIDs []string, timestamp int64) error {
	if len(conversationIDs) == 0 {
		return nil
	}
	filter := bson.M{
		"conversation_id": bson.M{"$in": conversationIDs},
		"send_time":       bson.M{"$lte": timestamp},
	}
	_, err := d.collection.DeleteMany(ctx, filter)
	return err
}

// RevokeMsg Ϣ content_type  content
func (d *MsgDocDatabase) RevokeMsg(ctx context.Context, conversationID string, seq int64, contentType int32, content []byte) error {
	filter := bson.M{
		"conversation_id": conversationID,
		"seq":             seq,
	}
	update := bson.M{
		"$set": bson.M{
			"content_type": contentType,
			"content":      content,
		},
	}
	_, err := d.collection.UpdateOne(ctx, filter, update)
	return err
}
