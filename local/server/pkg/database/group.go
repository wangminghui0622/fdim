package database

import (
	"context"
	"time"

	"fdim/pkg/constant"
	"fdim/pkg/model"
	pt "fdim/protocol/constant"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GroupDatabase struct {
	groupCollection   *mongo.Collection
	memberCollection  *mongo.Collection
	requestCollection *mongo.Collection
}

func NewGroupDatabase(db *MongoDB) *GroupDatabase {
	return &GroupDatabase{
		groupCollection:   db.GetCollection("groups"),
		memberCollection:  db.GetCollection("group_members"),
		requestCollection: db.GetCollection("group_requests"),
	}
}

// TakeGroup ȡȺ
func (d *GroupDatabase) TakeGroup(ctx context.Context, groupID string) (*model.Group, error) {
	var group model.Group
	err := d.groupCollection.FindOne(ctx, bson.M{"group_id": groupID}).Decode(&group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// FindGroup Ⱥ
func (d *GroupDatabase) FindGroup(ctx context.Context, groupIDs []string) ([]*model.Group, error) {
	filter := bson.M{"group_id": bson.M{"$in": groupIDs}}
	cursor, err := d.groupCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []*model.Group
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// CreateGroup Ⱥ
func (d *GroupDatabase) CreateGroup(ctx context.Context, groups []*model.Group, members []*model.GroupMember) error {
	// Ⱥ
	if len(groups) > 0 {
		docs := make([]interface{}, len(groups))
		for i, g := range groups {
			docs[i] = g
		}
		if _, err := d.groupCollection.InsertMany(ctx, docs); err != nil {
			return err
		}
	}

	// Ⱥ??
	if len(members) > 0 {
		docs := make([]interface{}, len(members))
		for i, m := range members {
			docs[i] = m
		}
		if _, err := d.memberCollection.InsertMany(ctx, docs); err != nil {
			return err
		}
	}

	return nil
}

// UpdateGroupMap ȺϢ
func (d *GroupDatabase) UpdateGroupMap(ctx context.Context, groupID string, data map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}
	filter := bson.M{"group_id": groupID}
	update := bson.M{"$set": data}
	_, err := d.groupCollection.UpdateOne(ctx, filter, update)
	return err
}

// FindGroupMembers Ⱥ??
func (d *GroupDatabase) FindGroupMembers(ctx context.Context, groupID string, userIDs []string) ([]*model.GroupMember, error) {
	filter := bson.M{
		"group_id": groupID,
		"user_id":  bson.M{"$in": userIDs},
	}
	cursor, err := d.memberCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var members []*model.GroupMember
	if err := cursor.All(ctx, &members); err != nil {
		return nil, err
	}
	return members, nil
}

// TakeGroupMember ȡȺ??
func (d *GroupDatabase) TakeGroupMember(ctx context.Context, groupID, userID string) (*model.GroupMember, error) {
	var member model.GroupMember
	err := d.memberCollection.FindOne(ctx, bson.M{
		"group_id": groupID,
		"user_id":  userID,
	}).Decode(&member)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// FindGroupMemberAll ȺԱ
func (d *GroupDatabase) FindGroupMemberAll(ctx context.Context, groupID string) ([]*model.GroupMember, error) {
	filter := bson.M{"group_id": groupID}
	cursor, err := d.memberCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var members []*model.GroupMember
	if err := cursor.All(ctx, &members); err != nil {
		return nil, err
	}
	return members, nil
}

// DeleteGroupMember ɾȺ??
func (d *GroupDatabase) DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	filter := bson.M{
		"group_id": groupID,
		"user_id":  bson.M{"$in": userIDs},
	}
	_, err := d.memberCollection.DeleteMany(ctx, filter)
	return err
}

// UpdateGroupMemberMap ȺԱ??
func (d *GroupDatabase) UpdateGroupMemberMap(ctx context.Context, groupID, userID string, data map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}
	filter := bson.M{
		"group_id": groupID,
		"user_id":  userID,
	}
	update := bson.M{"$set": data}
	_, err := d.memberCollection.UpdateOne(ctx, filter, update)
	return err
}

// TakeGroupOwner ȡȺ
func (d *GroupDatabase) TakeGroupOwner(ctx context.Context, groupID string) (*model.GroupMember, error) {
	filter := bson.M{
		"group_id":   groupID,
		"role_level": 100, // GroupOwner
	}
	var member model.GroupMember
	err := d.memberCollection.FindOne(ctx, filter).Decode(&member)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// FindGroupsOwner ȡȺȺ??
func (d *GroupDatabase) FindGroupsOwner(ctx context.Context, groupIDs []string) ([]*model.GroupMember, error) {
	if len(groupIDs) == 0 {
		return nil, nil
	}
	filter := bson.M{
		"group_id":   bson.M{"$in": groupIDs},
		"role_level": 100, // GroupOwner
	}
	cursor, err := d.memberCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var owners []*model.GroupMember
	if err := cursor.All(ctx, &owners); err != nil {
		return nil, err
	}
	return owners, nil
}

// CreateGroupRequest Ⱥ
func (d *GroupDatabase) CreateGroupRequest(ctx context.Context, requests []*model.GroupRequest) error {
	if len(requests) == 0 {
		return nil
	}
	docs := make([]interface{}, len(requests))
	for i, r := range requests {
		docs[i] = r
	}
	_, err := d.requestCollection.InsertMany(ctx, docs)
	return err
}

// FindGroupRequests Ⱥ
func (d *GroupDatabase) FindGroupRequests(ctx context.Context, groupID string, userIDs []string) ([]*model.GroupRequest, error) {
	filter := bson.M{
		"group_id": groupID,
		"user_id":  bson.M{"$in": userIDs},
	}
	cursor, err := d.requestCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var requests []*model.GroupRequest
	if err := cursor.All(ctx, &requests); err != nil {
		return nil, err
	}
	return requests, nil
}

// PageGroupRequest ҳȡȺ루ûӽǣ
func (d *GroupDatabase) PageGroupRequest(ctx context.Context, userID string, groupIDs []string, handleResults []int32, offset, limit int32) (int64, []*model.GroupRequest, error) {
	filter := bson.M{"user_id": userID}
	if len(groupIDs) > 0 {
		filter["group_id"] = bson.M{"$in": groupIDs}
	}
	if len(handleResults) > 0 {
		filter["handle_result"] = bson.M{"$in": handleResults}
	}

	// ȡ
	total, err := d.requestCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "req_time", Value: -1}})
	cursor, err := d.requestCollection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var requests []*model.GroupRequest
	if err := cursor.All(ctx, &requests); err != nil {
		return 0, nil, err
	}
	return total, requests, nil
}

// PageGroupRequestByGroup ҳȡȺ루Ⱥӽǣ
func (d *GroupDatabase) PageGroupRequestByGroup(ctx context.Context, groupIDs []string, handleResults []int32, offset, limit int32) (int64, []*model.GroupRequest, error) {
	if len(groupIDs) == 0 {
		return 0, nil, nil
	}
	filter := bson.M{"group_id": bson.M{"$in": groupIDs}}
	if len(handleResults) > 0 {
		filter["handle_result"] = bson.M{"$in": handleResults}
	}

	// ȡ
	total, err := d.requestCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "req_time", Value: -1}})
	cursor, err := d.requestCollection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var requests []*model.GroupRequest
	if err := cursor.All(ctx, &requests); err != nil {
		return 0, nil, err
	}
	return total, requests, nil
}

// GetUnhandledGroupRequestCount ȡδȺ
func (d *GroupDatabase) GetUnhandledGroupRequestCount(ctx context.Context, groupIDs []string, ts int64) (int64, error) {
	if len(groupIDs) == 0 {
		return 0, nil
	}
	filter := bson.M{
		"group_id":      bson.M{"$in": groupIDs},
		"handle_result": 0, // δ??
	}
	if ts != 0 {
		filter["req_time"] = bson.M{"$gt": time.UnixMilli(ts)}
	}
	return d.requestCollection.CountDocuments(ctx, filter)
}

// FindGroupMemberByUserID ûȺ??
func (d *GroupDatabase) FindGroupMemberByUserID(ctx context.Context, userID string) ([]*model.GroupMember, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := d.memberCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var members []*model.GroupMember
	if err := cursor.All(ctx, &members); err != nil {
		return nil, err
	}
	return members, nil
}

// FindGroupMemberNum ȡȺԱ
func (d *GroupDatabase) FindGroupMemberNum(ctx context.Context, groupID string) (uint32, error) {
	count, err := d.memberCollection.CountDocuments(ctx, bson.M{"group_id": groupID})
	if err != nil {
		return 0, err
	}
	return uint32(count), nil
}

// TakeGroupRequest ȡȺ
func (d *GroupDatabase) TakeGroupRequest(ctx context.Context, groupID, userID string) (*model.GroupRequest, error) {
	var request model.GroupRequest
	err := d.requestCollection.FindOne(ctx, bson.M{
		"group_id": groupID,
		"user_id":  userID,
	}).Decode(&request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// UpdateGroupRequest Ⱥ
func (d *GroupDatabase) UpdateGroupRequest(ctx context.Context, groupID, userID string, data map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}
	filter := bson.M{
		"group_id": groupID,
		"user_id":  userID,
	}
	update := bson.M{"$set": data}
	_, err := d.requestCollection.UpdateOne(ctx, filter, update)
	return err
}

// HandlerGroupRequest Ⱥ루״̬ͬ򴴽Ա
func (d *GroupDatabase) HandlerGroupRequest(ctx context.Context, groupID, userID, handledMsg string, handleResult int32, member *model.GroupMember) error {
	// ״??
	updateData := map[string]interface{}{
		"handle_result":  handleResult,
		"handled_msg":    handledMsg,
		"handled_time":   time.Now(),
		"handle_user_id": "", // TODO: contextȡ
	}
	if err := d.UpdateGroupRequest(ctx, groupID, userID, updateData); err != nil {
		return err
	}

	// ͬҳԱΪգȺԱ
	if handleResult == constant.GroupRequestAgree && member != nil {
		// ǷѾȺ??
		_, err := d.TakeGroupMember(ctx, groupID, userID)
		if err == nil {
			// ѾǳԱҪٴδ??
			return nil
		}
		// Ⱥ??
		if err := d.CreateGroup(ctx, nil, []*model.GroupMember{member}); err != nil {
			return err
		}
	}

	return nil
}

// MapGroupMemberNum ȡȺԱӳ
func (d *GroupDatabase) MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"group_id": bson.M{"$in": groupIDs}}},
		{"$group": bson.M{
			"_id":   "$group_id",
			"count": bson.M{"$sum": 1},
		}},
	}
	cursor, err := d.memberCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[string]uint32)
	for cursor.Next(ctx) {
		var item struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		result[item.ID] = uint32(item.Count)
	}
	return result, nil
}

// CountGroups ͳȺ
func (d *GroupDatabase) CountGroups(ctx context.Context) (int64, error) {
	return d.groupCollection.CountDocuments(ctx, bson.M{})
}

// CountGroupsBeforeTime ͳָʱ֮ǰȺ
func (d *GroupDatabase) CountGroupsBeforeTime(ctx context.Context, timestamp int64) (int64, error) {
	filter := bson.M{
		"create_time": bson.M{"$lt": time.UnixMilli(timestamp)},
	}
	return d.groupCollection.CountDocuments(ctx, filter)
}

// FindGroupMemberRoleLevels ݽɫȺ??
func (d *GroupDatabase) FindGroupMemberRoleLevels(ctx context.Context, groupID string, roleLevels []int32) ([]*model.GroupMember, error) {
	if len(roleLevels) == 0 {
		return nil, nil
	}
	filter := bson.M{
		"group_id":   groupID,
		"role_level": bson.M{"$in": roleLevels},
	}
	cursor, err := d.memberCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var members []*model.GroupMember
	if err := cursor.All(ctx, &members); err != nil {
		return nil, err
	}
	return members, nil
}

// FindRoleLevelUserIDs ݽɫûIDб
func (d *GroupDatabase) FindRoleLevelUserIDs(ctx context.Context, groupID string, roleLevel int32) ([]string, error) {
	filter := bson.M{
		"group_id":   groupID,
		"role_level": roleLevel,
	}
	opts := options.Find().SetProjection(bson.M{"user_id": 1, "_id": 0})
	cursor, err := d.memberCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userIDs []string
	for cursor.Next(ctx) {
		var result struct {
			UserID string `bson:"user_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, result.UserID)
	}
	return userIDs, nil
}

// IsGroupAdmin ûǷȺԱ
func (d *GroupDatabase) IsGroupAdmin(ctx context.Context, groupID, userID string) (bool, error) {
	member, err := d.TakeGroupMember(ctx, groupID, userID)
	if err != nil {
		return false, err
	}
	return member.RoleLevel == pt.GroupOwner || member.RoleLevel == pt.GroupAdmin, nil
}
