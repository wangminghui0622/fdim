package database

import (
	"context"
	"fmt"
	"time"

	"fdim/protocol/sdkws"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AdminDatabase 管理员数据库接口
type AdminDatabase interface {
	// 管理员账户管理
	GetAdmin(ctx context.Context, account string) (*Admin, error)
	GetAdminUserID(ctx context.Context, userID string) (*Admin, error)
	AddAdminAccount(ctx context.Context, admins []*Admin) error
	DelAdminAccount(ctx context.Context, userIDs []string) error
	UpdateAdmin(ctx context.Context, userID string, update map[string]interface{}) error
	ChangePassword(ctx context.Context, userID string, newPassword string) error
	SearchAdminAccount(ctx context.Context, pagination *sdkws.RequestPagination) (int64, []*Admin, error)

	// 默认好友管理
	FindDefaultFriend(ctx context.Context, userIDs []string) ([]string, error)
	AddDefaultFriend(ctx context.Context, userIDs []string) error
	DelDefaultFriend(ctx context.Context, userIDs []string) error
	SearchDefaultFriend(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []*DefaultFriend, error)

	// 默认群组管理
	FindDefaultGroup(ctx context.Context, groupIDs []string) ([]string, error)
	AddDefaultGroup(ctx context.Context, groupIDs []string) error
	DelDefaultGroup(ctx context.Context, groupIDs []string) error
	SearchDefaultGroup(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []string, error)

	// 邀请码管理
	FindInvitationRegister(ctx context.Context, codes []string) ([]*InvitationRegister, error)
	AddInvitationCode(ctx context.Context, codes []*InvitationRegister) error
	DelInvitationCode(ctx context.Context, codes []string) error
	UseInvitationCode(ctx context.Context, code string, userID string) error
	SearchInvitationCode(ctx context.Context, status int32, userIDs []string, codes []string, keyword string, pagination *sdkws.RequestPagination) (int64, []*InvitationRegister, error)

	// IP禁止管理
	FindIPForbidden(ctx context.Context, ips []string) ([]*IPForbidden, error)
	AddIPForbidden(ctx context.Context, forbiddens []*IPForbidden) error
	DelIPForbidden(ctx context.Context, ips []string) error
	SearchIPForbidden(ctx context.Context, keyword string, status int32, pagination *sdkws.RequestPagination) (int64, []*IPForbidden, error)

	// 用户IP登录限制
	GetLimitUserLoginIP(ctx context.Context, userID string, ip string) (*LimitUserLoginIP, error)
	AddUserIPLimitLogin(ctx context.Context, limits []*LimitUserLoginIP) error
	DelUserIPLimitLogin(ctx context.Context, limits []*LimitUserLoginIP) error
	SearchUserLimitLogin(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []*LimitUserLoginIP, error)
	CountLimitUserLoginIP(ctx context.Context, userID string) (int64, error)

	// 用户封禁管理
	GetBlockInfo(ctx context.Context, userID string) (*BlockUser, error)
	AddBlockUser(ctx context.Context, blocks []*BlockUser) error
	DelBlockUser(ctx context.Context, userIDs []string) error
	SearchBlockUser(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []*BlockUser, error)
	FindUserBlockInfo(ctx context.Context, userIDs []string) ([]*BlockUser, error)

	// Token管理
	CacheToken(ctx context.Context, userID string, token string, expire int64) error
	GetTokens(ctx context.Context, userID string) (map[string]int32, error)
	InvalidateToken(ctx context.Context, userID string) error

	// 小程序管理
	FindApplet(ctx context.Context, appletIDs []string) ([]*Applet, error)
	FindOnShelfApplet(ctx context.Context) ([]*Applet, error)
	AddApplet(ctx context.Context, applets []*Applet) error
	DelApplet(ctx context.Context, appletIDs []string) error
	UpdateApplet(ctx context.Context, appletID string, update map[string]interface{}) error
	SearchApplet(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []*Applet, error)

	// 客户端配置管理
	GetConfig(ctx context.Context) (map[string]string, error)
	SetConfig(ctx context.Context, config map[string]string) error
	DelConfig(ctx context.Context, keys []string) error

	// 应用版本管理
	LatestVersion(ctx context.Context, platform string) (*ApplicationVersion, error)
	AddApplicationVersion(ctx context.Context, versions []*ApplicationVersion) error
	UpdateApplicationVersion(ctx context.Context, id string, update map[string]interface{}) error
	DeleteApplicationVersion(ctx context.Context, ids []string) error
	PageApplicationVersion(ctx context.Context, platforms []string, pagination *sdkws.RequestPagination) (int64, []*ApplicationVersion, error)
}

// Admin 管理员模型
type Admin struct {
	Account    string
	Password   string
	FaceURL    string
	Nickname   string
	UserID     string
	Level      int32
	CreateTime int64
}

// DefaultFriend 默认好友模型
type DefaultFriend struct {
	UserID     string    `bson:"user_id"`
	CreateTime time.Time `bson:"create_time"`
}

// InvitationRegister 邀请码注册模型
type InvitationRegister struct {
	InvitationCode string    `bson:"invitation_code"`
	UsedByUserID   string    `bson:"used_by_user_id"`
	CreateTime     time.Time `bson:"create_time"`
}

// IPForbidden IP禁止模型
type IPForbidden struct {
	IP            string    `bson:"ip"`
	LimitRegister bool      `bson:"limit_register"`
	LimitLogin    bool      `bson:"limit_login"`
	CreateTime    time.Time `bson:"create_time"`
}

// LimitUserLoginIP 用户IP登录限制模型
type LimitUserLoginIP struct {
	UserID     string    `bson:"user_id"`
	IP         string    `bson:"ip"`
	CreateTime time.Time `bson:"create_time"`
}

// BlockUser 封禁用户模型
type BlockUser struct {
	UserID         string    `bson:"user_id"`
	Reason         string    `bson:"reason"`
	OperatorUserID string    `bson:"operator_user_id"`
	CreateTime     time.Time `bson:"create_time"`
}

// Applet 小程序模型
type Applet struct {
	ID         string    `bson:"id"`
	Name       string    `bson:"name"`
	AppID      string    `bson:"app_id"`
	Icon       string    `bson:"icon"`
	URL        string    `bson:"url"`
	MD5        string    `bson:"md5"`
	Size       int64     `bson:"size"`
	Version    string    `bson:"version"`
	Priority   uint32    `bson:"priority"`
	Status     uint32    `bson:"status"`
	CreateTime time.Time `bson:"create_time"`
}

// ApplicationVersion 应用版本模型
type ApplicationVersion struct {
	ID         string    `bson:"id"`
	Platform   string    `bson:"platform"`
	Version    string    `bson:"version"`
	URL        string    `bson:"url"`
	Text       string    `bson:"text"`
	Force      bool      `bson:"force"`
	Latest     bool      `bson:"latest"`
	Hot        bool      `bson:"hot"`
	CreateTime time.Time `bson:"create_time"`
}

// NewAdminDatabase 创建管理员数据库实例
func NewAdminDatabase(mongoDB *MongoDB, redisClient *redis.Client) AdminDatabase {
	return &adminDatabase{
		adminCollection:         mongoDB.GetCollection("admin"),
		defaultFriendCollection: mongoDB.GetCollection("default_friend"),
		defaultGroupCollection:  mongoDB.GetCollection("default_group"),
		invitationCollection:    mongoDB.GetCollection("invitation_code"),
		ipForbiddenCollection:   mongoDB.GetCollection("ip_forbidden"),
		limitLoginCollection:    mongoDB.GetCollection("limit_user_login_ip"),
		blockUserCollection:     mongoDB.GetCollection("block_user"),
		appletCollection:        mongoDB.GetCollection("applet"),
		configCollection:        mongoDB.GetCollection("client_config"),
		versionCollection:       mongoDB.GetCollection("application_version"),
		redis:                   redisClient,
	}
}

type adminDatabase struct {
	adminCollection         *mongo.Collection
	defaultFriendCollection *mongo.Collection
	defaultGroupCollection  *mongo.Collection
	invitationCollection    *mongo.Collection
	ipForbiddenCollection   *mongo.Collection
	limitLoginCollection    *mongo.Collection
	blockUserCollection     *mongo.Collection
	appletCollection        *mongo.Collection
	configCollection        *mongo.Collection
	versionCollection       *mongo.Collection
	redis                   *redis.Client
}

// TODO: 实现所有 AdminDatabase 接口方法
// 这里先提供框架，具体实现需要根据 MongoDB 操作来完成

func (d *adminDatabase) GetAdmin(ctx context.Context, account string) (*Admin, error) {
	var admin Admin
	err := d.adminCollection.FindOne(ctx, bson.M{"account": account}).Decode(&admin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("admin not found: %w", err)
		}
		return nil, err
	}
	return &admin, nil
}

func (d *adminDatabase) GetAdminUserID(ctx context.Context, userID string) (*Admin, error) {
	var admin Admin
	err := d.adminCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&admin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("admin not found: %w", err)
		}
		return nil, err
	}
	return &admin, nil
}

func (d *adminDatabase) AddAdminAccount(ctx context.Context, admins []*Admin) error {
	docs := make([]interface{}, len(admins))
	for i, admin := range admins {
		docs[i] = admin
	}
	_, err := d.adminCollection.InsertMany(ctx, docs)
	return err
}

func (d *adminDatabase) DelAdminAccount(ctx context.Context, userIDs []string) error {
	_, err := d.adminCollection.DeleteMany(ctx, bson.M{"user_id": bson.M{"$in": userIDs}})
	return err
}

func (d *adminDatabase) UpdateAdmin(ctx context.Context, userID string, update map[string]interface{}) error {
	_, err := d.adminCollection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": update})
	return err
}

func (d *adminDatabase) ChangePassword(ctx context.Context, userID string, newPassword string) error {
	_, err := d.adminCollection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": bson.M{"password": newPassword}})
	return err
}

func (d *adminDatabase) SearchAdminAccount(ctx context.Context, pagination *sdkws.RequestPagination) (int64, []*Admin, error) {
	// 计算分页参数
	pageNumber := int32(1)
	showNumber := int32(10)
	if pagination != nil {
		if pagination.PageNumber > 0 {
			pageNumber = pagination.PageNumber
		}
		if pagination.ShowNumber > 0 {
			showNumber = pagination.ShowNumber
		}
	}

	offset := int64((pageNumber - 1) * showNumber)
	limit := int64(showNumber)

	// 构建查询条件（只查询普通管理员，级别为 80）
	filter := bson.M{"level": 80} // NormalAdmin

	// 获取总数
	total, err := d.adminCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// 分页查询
	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "create_time", Value: -1}})

	cursor, err := d.adminCollection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var admins []*Admin
	if err := cursor.All(ctx, &admins); err != nil {
		return 0, nil, err
	}

	return total, admins, nil
}

func (d *adminDatabase) FindDefaultFriend(ctx context.Context, userIDs []string) ([]string, error) {
	var filter bson.M
	if len(userIDs) > 0 {
		filter = bson.M{"user_id": bson.M{"$in": userIDs}}
	} else {
		filter = bson.M{}
	}

	cursor, err := d.defaultFriendCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*DefaultFriend
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	userIDList := make([]string, 0, len(results))
	for _, r := range results {
		userIDList = append(userIDList, r.UserID)
	}
	return userIDList, nil
}

func (d *adminDatabase) AddDefaultFriend(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	docs := make([]interface{}, len(userIDs))
	now := time.Now()
	for i, userID := range userIDs {
		docs[i] = &DefaultFriend{
			UserID:     userID,
			CreateTime: now,
		}
	}

	_, err := d.defaultFriendCollection.InsertMany(ctx, docs)
	return err
}

func (d *adminDatabase) DelDefaultFriend(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	_, err := d.defaultFriendCollection.DeleteMany(ctx, bson.M{"user_id": bson.M{"$in": userIDs}})
	return err
}

func (d *adminDatabase) SearchDefaultFriend(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []*DefaultFriend, error) {
	filter := bson.M{}
	if keyword != "" {
		filter["user_id"] = bson.M{"$regex": keyword, "$options": "i"}
	}

	// 计算分页
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

	// 获取总数
	total, err := d.defaultFriendCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// 分页查询
	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "create_time", Value: -1}})

	cursor, err := d.defaultFriendCollection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var results []*DefaultFriend
	if err := cursor.All(ctx, &results); err != nil {
		return 0, nil, err
	}

	return total, results, nil
}

func (d *adminDatabase) FindDefaultGroup(ctx context.Context, groupIDs []string) ([]string, error) {
	var filter bson.M
	if len(groupIDs) > 0 {
		filter = bson.M{"group_id": bson.M{"$in": groupIDs}}
	} else {
		filter = bson.M{}
	}

	cursor, err := d.defaultGroupCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	groupIDList := make([]string, 0, len(results))
	for _, r := range results {
		if groupID, ok := r["group_id"].(string); ok {
			groupIDList = append(groupIDList, groupID)
		}
	}
	return groupIDList, nil
}

func (d *adminDatabase) AddDefaultGroup(ctx context.Context, groupIDs []string) error {
	if len(groupIDs) == 0 {
		return nil
	}

	docs := make([]interface{}, len(groupIDs))
	now := time.Now()
	for i, groupID := range groupIDs {
		docs[i] = bson.M{
			"group_id":    groupID,
			"create_time": now,
		}
	}

	_, err := d.defaultGroupCollection.InsertMany(ctx, docs)
	return err
}

func (d *adminDatabase) DelDefaultGroup(ctx context.Context, groupIDs []string) error {
	if len(groupIDs) == 0 {
		return nil
	}

	_, err := d.defaultGroupCollection.DeleteMany(ctx, bson.M{"group_id": bson.M{"$in": groupIDs}})
	return err
}

func (d *adminDatabase) SearchDefaultGroup(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []string, error) {
	filter := bson.M{}
	if keyword != "" {
		filter["group_id"] = bson.M{"$regex": keyword, "$options": "i"}
	}

	// 计算分页
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

	// 获取总数
	total, err := d.defaultGroupCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// 分页查询
	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "create_time", Value: -1}})

	cursor, err := d.defaultGroupCollection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return 0, nil, err
	}

	groupIDList := make([]string, 0, len(results))
	for _, r := range results {
		if groupID, ok := r["group_id"].(string); ok {
			groupIDList = append(groupIDList, groupID)
		}
	}

	return total, groupIDList, nil
}

func (d *adminDatabase) FindInvitationRegister(ctx context.Context, codes []string) ([]*InvitationRegister, error) {
	if len(codes) == 0 {
		return nil, nil
	}

	filter := bson.M{"invitation_code": bson.M{"$in": codes}}
	cursor, err := d.invitationCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*InvitationRegister
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (d *adminDatabase) AddInvitationCode(ctx context.Context, codes []*InvitationRegister) error {
	if len(codes) == 0 {
		return nil
	}

	docs := make([]interface{}, len(codes))
	for i, code := range codes {
		docs[i] = code
	}

	_, err := d.invitationCollection.InsertMany(ctx, docs)
	return err
}

func (d *adminDatabase) DelInvitationCode(ctx context.Context, codes []string) error {
	if len(codes) == 0 {
		return nil
	}

	_, err := d.invitationCollection.DeleteMany(ctx, bson.M{"invitation_code": bson.M{"$in": codes}})
	return err
}

func (d *adminDatabase) UseInvitationCode(ctx context.Context, code string, userID string) error {
	// 更新邀请码，标记为已使用
	update := bson.M{
		"$set": bson.M{
			"used_by_user_id": userID,
		},
	}
	_, err := d.invitationCollection.UpdateOne(ctx, bson.M{"invitation_code": code}, update)
	return err
}

func (d *adminDatabase) SearchInvitationCode(ctx context.Context, status int32, userIDs []string, codes []string, keyword string, pagination *sdkws.RequestPagination) (int64, []*InvitationRegister, error) {
	filter := bson.M{}

	// 状态过滤：0=未使用(used_by_user_id为空), 1=已使用(used_by_user_id不为空)
	if status == 0 {
		filter["used_by_user_id"] = bson.M{"$in": []interface{}{nil, ""}}
	} else if status == 1 {
		filter["used_by_user_id"] = bson.M{"$ne": ""}
	}

	// 用户ID过滤
	if len(userIDs) > 0 {
		filter["used_by_user_id"] = bson.M{"$in": userIDs}
	}

	// 邀请码过滤
	if len(codes) > 0 {
		filter["invitation_code"] = bson.M{"$in": codes}
	}

	// 关键词搜索
	if keyword != "" {
		filter["invitation_code"] = bson.M{"$regex": keyword, "$options": "i"}
	}

	// 计算分页
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

	// 获取总数
	total, err := d.invitationCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// 分页查询
	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "create_time", Value: -1}})

	cursor, err := d.invitationCollection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var results []*InvitationRegister
	if err := cursor.All(ctx, &results); err != nil {
		return 0, nil, err
	}

	return total, results, nil
}

func (d *adminDatabase) FindIPForbidden(ctx context.Context, ips []string) ([]*IPForbidden, error) {
	if len(ips) == 0 {
		return nil, nil
	}

	filter := bson.M{"ip": bson.M{"$in": ips}}
	cursor, err := d.ipForbiddenCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*IPForbidden
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (d *adminDatabase) AddIPForbidden(ctx context.Context, forbiddens []*IPForbidden) error {
	if len(forbiddens) == 0 {
		return nil
	}

	docs := make([]interface{}, len(forbiddens))
	for i, fb := range forbiddens {
		docs[i] = fb
	}

	_, err := d.ipForbiddenCollection.InsertMany(ctx, docs)
	return err
}

func (d *adminDatabase) DelIPForbidden(ctx context.Context, ips []string) error {
	if len(ips) == 0 {
		return nil
	}

	_, err := d.ipForbiddenCollection.DeleteMany(ctx, bson.M{"ip": bson.M{"$in": ips}})
	return err
}

func (d *adminDatabase) SearchIPForbidden(ctx context.Context, keyword string, status int32, pagination *sdkws.RequestPagination) (int64, []*IPForbidden, error) {
	filter := bson.M{}

	// 关键词搜索
	if keyword != "" {
		filter["ip"] = bson.M{"$regex": keyword, "$options": "i"}
	}

	// 状态过滤：根据limit_register和limit_login
	// status: 0=全部, 1=仅禁止注册, 2=仅禁止登录, 3=禁止注册和登录
	if status > 0 {
		if status == 1 {
			filter["limit_register"] = true
		} else if status == 2 {
			filter["limit_login"] = true
		} else if status == 3 {
			filter["limit_register"] = true
			filter["limit_login"] = true
		}
	}

	// 计算分页
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

	// 获取总数
	total, err := d.ipForbiddenCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// 分页查询
	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "create_time", Value: -1}})

	cursor, err := d.ipForbiddenCollection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var results []*IPForbidden
	if err := cursor.All(ctx, &results); err != nil {
		return 0, nil, err
	}

	return total, results, nil
}

func (d *adminDatabase) GetLimitUserLoginIP(ctx context.Context, userID string, ip string) (*LimitUserLoginIP, error) {
	filter := bson.M{
		"user_id": userID,
		"ip":      ip,
	}
	var result LimitUserLoginIP
	err := d.limitLoginCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("limit not found: %w", err)
		}
		return nil, err
	}
	return &result, nil
}

func (d *adminDatabase) AddUserIPLimitLogin(ctx context.Context, limits []*LimitUserLoginIP) error {
	if len(limits) == 0 {
		return nil
	}

	docs := make([]interface{}, len(limits))
	for i, limit := range limits {
		docs[i] = limit
	}

	_, err := d.limitLoginCollection.InsertMany(ctx, docs)
	return err
}

func (d *adminDatabase) DelUserIPLimitLogin(ctx context.Context, limits []*LimitUserLoginIP) error {
	if len(limits) == 0 {
		return nil
	}

	// 构建删除条件
	filters := make([]bson.M, len(limits))
	for i, limit := range limits {
		filters[i] = bson.M{
			"user_id": limit.UserID,
			"ip":      limit.IP,
		}
	}

	_, err := d.limitLoginCollection.DeleteMany(ctx, bson.M{"$or": filters})
	return err
}

func (d *adminDatabase) SearchUserLimitLogin(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []*LimitUserLoginIP, error) {
	filter := bson.M{}
	if keyword != "" {
		// 支持按userID或IP搜索
		filter["$or"] = []bson.M{
			{"user_id": bson.M{"$regex": keyword, "$options": "i"}},
			{"ip": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}

	// 计算分页
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

	// 获取总数
	total, err := d.limitLoginCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// 分页查询
	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "create_time", Value: -1}})

	cursor, err := d.limitLoginCollection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var results []*LimitUserLoginIP
	if err := cursor.All(ctx, &results); err != nil {
		return 0, nil, err
	}

	return total, results, nil
}

func (d *adminDatabase) CountLimitUserLoginIP(ctx context.Context, userID string) (int64, error) {
	count, err := d.limitLoginCollection.CountDocuments(ctx, bson.M{"user_id": userID})
	return count, err
}

func (d *adminDatabase) GetBlockInfo(ctx context.Context, userID string) (*BlockUser, error) {
	var result BlockUser
	err := d.blockUserCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("block info not found: %w", err)
		}
		return nil, err
	}
	return &result, nil
}

func (d *adminDatabase) AddBlockUser(ctx context.Context, blocks []*BlockUser) error {
	if len(blocks) == 0 {
		return nil
	}

	docs := make([]interface{}, len(blocks))
	for i, block := range blocks {
		docs[i] = block
	}

	_, err := d.blockUserCollection.InsertMany(ctx, docs)
	return err
}

func (d *adminDatabase) DelBlockUser(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	_, err := d.blockUserCollection.DeleteMany(ctx, bson.M{"user_id": bson.M{"$in": userIDs}})
	return err
}

func (d *adminDatabase) SearchBlockUser(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []*BlockUser, error) {
	filter := bson.M{}
	if keyword != "" {
		// 支持按userID搜索
		filter["user_id"] = bson.M{"$regex": keyword, "$options": "i"}
	}

	// 计算分页
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

	// 获取总数
	total, err := d.blockUserCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// 分页查询
	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "create_time", Value: -1}})

	cursor, err := d.blockUserCollection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var results []*BlockUser
	if err := cursor.All(ctx, &results); err != nil {
		return 0, nil, err
	}

	return total, results, nil
}

func (d *adminDatabase) FindUserBlockInfo(ctx context.Context, userIDs []string) ([]*BlockUser, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	filter := bson.M{"user_id": bson.M{"$in": userIDs}}
	cursor, err := d.blockUserCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*BlockUser
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (d *adminDatabase) CacheToken(ctx context.Context, userID string, token string, expire int64) error {
	if d.redis == nil {
		return fmt.Errorf("redis client not initialized")
	}
	// 使用 Hash 结构存储 token: ADMIN_UID_TOKEN_STATUS:{userID}
	key := fmt.Sprintf("ADMIN_UID_TOKEN_STATUS:%s", userID)

	// 使用 HSet 存储 token，值为 1 表示正常token
	if err := d.redis.HSet(ctx, key, token, 1).Err(); err != nil {
		return err
	}

	// 设置整个 Hash 的过期时间
	expireDuration := time.Duration(expire) * time.Second
	return d.redis.Expire(ctx, key, expireDuration).Err()
}

func (d *adminDatabase) GetTokens(ctx context.Context, userID string) (map[string]int32, error) {
	if d.redis == nil {
		return nil, fmt.Errorf("redis client not initialized")
	}
	// 使用 Hash 结构: ADMIN_UID_TOKEN_STATUS:{userID}
	key := fmt.Sprintf("ADMIN_UID_TOKEN_STATUS:%s", userID)

	// 获取 Hash 中的所有字段和值
	m, err := d.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// 转换为 map[string]int32
	result := make(map[string]int32)
	for token, value := range m {
		// value 是字符串 "1"，需要转换为 int32
		var flag int32
		if value == "1" {
			flag = 1
		}
		result[token] = flag
	}

	return result, nil
}

func (d *adminDatabase) InvalidateToken(ctx context.Context, userID string) error {
	if d.redis == nil {
		return fmt.Errorf("redis client not initialized")
	}
	// 使用 Hash 结构: ADMIN_UID_TOKEN_STATUS:{userID}
	key := fmt.Sprintf("ADMIN_UID_TOKEN_STATUS:%s", userID)

	// 删除整个 Hash（使该用户的所有 token 失效）
	return d.redis.Del(ctx, key).Err()
}

func (d *adminDatabase) FindApplet(ctx context.Context, appletIDs []string) ([]*Applet, error) {
	if len(appletIDs) == 0 {
		return nil, nil
	}

	filter := bson.M{"id": bson.M{"$in": appletIDs}}
	cursor, err := d.appletCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*Applet
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (d *adminDatabase) AddApplet(ctx context.Context, applets []*Applet) error {
	if len(applets) == 0 {
		return nil
	}

	docs := make([]interface{}, len(applets))
	for i, applet := range applets {
		docs[i] = applet
	}

	_, err := d.appletCollection.InsertMany(ctx, docs)
	return err
}

func (d *adminDatabase) DelApplet(ctx context.Context, appletIDs []string) error {
	if len(appletIDs) == 0 {
		return nil
	}

	_, err := d.appletCollection.DeleteMany(ctx, bson.M{"id": bson.M{"$in": appletIDs}})
	return err
}

func (d *adminDatabase) UpdateApplet(ctx context.Context, appletID string, update map[string]interface{}) error {
	_, err := d.appletCollection.UpdateOne(ctx, bson.M{"id": appletID}, bson.M{"$set": update})
	return err
}

func (d *adminDatabase) SearchApplet(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []*Applet, error) {
	filter := bson.M{}
	if keyword != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": keyword, "$options": "i"}},
			{"app_id": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}

	// 计算分页
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

	// 获取总数
	total, err := d.appletCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// 分页查询
	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "create_time", Value: -1}})

	cursor, err := d.appletCollection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var results []*Applet
	if err := cursor.All(ctx, &results); err != nil {
		return 0, nil, err
	}

	return total, results, nil
}

func (d *adminDatabase) FindOnShelfApplet(ctx context.Context) ([]*Applet, error) {
	// 查找所有已上架的小程序（status=1）
	filter := bson.M{"status": 1}
	cursor, err := d.appletCollection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "priority", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*Applet
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (d *adminDatabase) GetConfig(ctx context.Context) (map[string]string, error) {
	cursor, err := d.configCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	config := make(map[string]string)
	for _, r := range results {
		if key, ok := r["key"].(string); ok {
			if value, ok := r["value"].(string); ok {
				config[key] = value
			}
		}
	}
	return config, nil
}

func (d *adminDatabase) SetConfig(ctx context.Context, config map[string]string) error {
	// 使用 upsert 更新或插入配置
	for key, value := range config {
		_, err := d.configCollection.UpdateOne(
			ctx,
			bson.M{"key": key},
			bson.M{"$set": bson.M{"key": key, "value": value}},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *adminDatabase) DelConfig(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	_, err := d.configCollection.DeleteMany(ctx, bson.M{"key": bson.M{"$in": keys}})
	return err
}

func (d *adminDatabase) LatestVersion(ctx context.Context, platform string) (*ApplicationVersion, error) {
	filter := bson.M{"platform": platform, "latest": true}
	opts := options.FindOne().SetSort(bson.D{{Key: "create_time", Value: -1}})

	var result ApplicationVersion
	err := d.versionCollection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("latest version not found for platform %s: %w", platform, err)
		}
		return nil, err
	}
	return &result, nil
}

func (d *adminDatabase) AddApplicationVersion(ctx context.Context, versions []*ApplicationVersion) error {
	if len(versions) == 0 {
		return nil
	}

	docs := make([]interface{}, len(versions))
	for i, version := range versions {
		docs[i] = version
	}

	_, err := d.versionCollection.InsertMany(ctx, docs)
	return err
}

func (d *adminDatabase) UpdateApplicationVersion(ctx context.Context, id string, update map[string]interface{}) error {
	_, err := d.versionCollection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": update})
	return err
}

func (d *adminDatabase) DeleteApplicationVersion(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := d.versionCollection.DeleteMany(ctx, bson.M{"id": bson.M{"$in": ids}})
	return err
}

func (d *adminDatabase) PageApplicationVersion(ctx context.Context, platforms []string, pagination *sdkws.RequestPagination) (int64, []*ApplicationVersion, error) {
	filter := bson.M{}
	if len(platforms) > 0 {
		filter["platform"] = bson.M{"$in": platforms}
	}

	// 计算分页
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

	// 获取总数
	total, err := d.versionCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// 分页查询
	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "create_time", Value: -1}})

	cursor, err := d.versionCollection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var results []*ApplicationVersion
	if err := cursor.All(ctx, &results); err != nil {
		return 0, nil, err
	}

	return total, results, nil
}
