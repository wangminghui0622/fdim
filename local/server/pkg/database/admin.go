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
	// 管理员账户
	GetAdmin(ctx context.Context, account string) (*Admin, error)
	GetAdminUserID(ctx context.Context, userID string) (*Admin, error)
	AddAdminAccount(ctx context.Context, admins []*Admin) error
	DelAdminAccount(ctx context.Context, userIDs []string) error
	UpdateAdmin(ctx context.Context, userID string, update map[string]interface{}) error
	ChangePassword(ctx context.Context, userID string, newPassword string) error
	SearchAdminAccount(ctx context.Context, pagination *sdkws.RequestPagination) (int64, []*Admin, error)

	// 默认好友关系
	FindDefaultFriend(ctx context.Context, userIDs []string) ([]string, error)
	AddDefaultFriend(ctx context.Context, userIDs []string) error
	DelDefaultFriend(ctx context.Context, userIDs []string) error
	SearchDefaultFriend(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []*DefaultFriend, error)

	// 默认群组
	FindDefaultGroup(ctx context.Context, groupIDs []string) ([]string, error)
	AddDefaultGroup(ctx context.Context, groupIDs []string) error
	DelDefaultGroup(ctx context.Context, groupIDs []string) error
	SearchDefaultGroup(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []string, error)

	// 邀请码
	FindInvitationRegister(ctx context.Context, codes []string) ([]*InvitationRegister, error)
	AddInvitationCode(ctx context.Context, codes []*InvitationRegister) error
	DelInvitationCode(ctx context.Context, codes []string) error
	UseInvitationCode(ctx context.Context, code string, userID string) error
	SearchInvitationCode(ctx context.Context, status int32, userIDs []string, codes []string, keyword string, pagination *sdkws.RequestPagination) (int64, []*InvitationRegister, error)

	// IP禁止
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

	// 用户封禁
	GetBlockInfo(ctx context.Context, userID string) (*BlockUser, error)
	AddBlockUser(ctx context.Context, blocks []*BlockUser) error
	DelBlockUser(ctx context.Context, userIDs []string) error
	SearchBlockUser(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []*BlockUser, error)
	FindUserBlockInfo(ctx context.Context, userIDs []string) ([]*BlockUser, error)

	// Token
	CacheToken(ctx context.Context, userID string, token string, expire int64) error
	GetTokens(ctx context.Context, userID string) (map[string]int32, error)
	InvalidateToken(ctx context.Context, userID string) error

	// 小程序
	FindApplet(ctx context.Context, appletIDs []string) ([]*Applet, error)
	FindOnShelfApplet(ctx context.Context) ([]*Applet, error)
	AddApplet(ctx context.Context, applets []*Applet) error
	DelApplet(ctx context.Context, appletIDs []string) error
	UpdateApplet(ctx context.Context, appletID string, update map[string]interface{}) error
	SearchApplet(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []*Applet, error)

	// 客户端配置
	GetConfig(ctx context.Context) (map[string]string, error)
	SetConfig(ctx context.Context, config map[string]string) error
	DelConfig(ctx context.Context, keys []string) error

	// 应用版本
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

// InvitationRegister 邀请注册模型
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

// BlockUser ûģ
type BlockUser struct {
	UserID         string    `bson:"user_id"`
	Reason         string    `bson:"reason"`
	OperatorUserID string    `bson:"operator_user_id"`
	CreateTime     time.Time `bson:"create_time"`
}

// Applet Сģ
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

// ApplicationVersion Ӧð汾ģ
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

// NewAdminDatabase Աݿʵ
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

// TODO: ʵ AdminDatabase ӿڷ
// ṩܣʵҪ MongoDB 

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
	// ҳ
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

	// ѯֻѯͨԱΪ 80
	filter := bson.M{"level": 80} // NormalAdmin

	// ȡ
	total, err := d.adminCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
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
	total, err := d.defaultFriendCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
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
	total, err := d.defaultGroupCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
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
	// 룬Ϊʹ
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

	// ״̬ˣ0=δʹ(used_by_user_idΪ), 1=ʹ(used_by_user_idΪ)
	if status == 0 {
		filter["used_by_user_id"] = bson.M{"$in": []interface{}{nil, ""}}
	} else if status == 1 {
		filter["used_by_user_id"] = bson.M{"$ne": ""}
	}

	// ûID
	if len(userIDs) > 0 {
		filter["used_by_user_id"] = bson.M{"$in": userIDs}
	}

	// 
	if len(codes) > 0 {
		filter["invitation_code"] = bson.M{"$in": codes}
	}

	// ؼ
	if keyword != "" {
		filter["invitation_code"] = bson.M{"$regex": keyword, "$options": "i"}
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
	total, err := d.invitationCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
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

	// ؼ
	if keyword != "" {
		filter["ip"] = bson.M{"$regex": keyword, "$options": "i"}
	}

	// ״̬ˣlimit_registerlimit_login
	// status: 0=ȫ, 1=ֹע, 2=ֹ¼, 3=ֹע͵¼
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
	total, err := d.ipForbiddenCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
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

	// ɾ
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
		// ְ֧userIDIP
		filter["$or"] = []bson.M{
			{"user_id": bson.M{"$regex": keyword, "$options": "i"}},
			{"ip": bson.M{"$regex": keyword, "$options": "i"}},
		}
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
	total, err := d.limitLoginCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
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
		// ְ֧userID
		filter["user_id"] = bson.M{"$regex": keyword, "$options": "i"}
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
	total, err := d.blockUserCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
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
	// ʹ Hash ṹ洢 token: ADMIN_UID_TOKEN_STATUS:{userID}
	key := fmt.Sprintf("ADMIN_UID_TOKEN_STATUS:%s", userID)

	// ʹ HSet 洢 tokenֵΪ 1 ʾtoken
	if err := d.redis.HSet(ctx, key, token, 1).Err(); err != nil {
		return err
	}

	//  Hash Ĺʱ
	expireDuration := time.Duration(expire) * time.Second
	return d.redis.Expire(ctx, key, expireDuration).Err()
}

func (d *adminDatabase) GetTokens(ctx context.Context, userID string) (map[string]int32, error) {
	if d.redis == nil {
		return nil, fmt.Errorf("redis client not initialized")
	}
	// ʹ Hash ṹ: ADMIN_UID_TOKEN_STATUS:{userID}
	key := fmt.Sprintf("ADMIN_UID_TOKEN_STATUS:%s", userID)

	// ȡ Hash еֶκֵ
	m, err := d.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// תΪ map[string]int32
	result := make(map[string]int32)
	for token, value := range m {
		// value ַ "1"ҪתΪ int32
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
	// ʹ Hash ṹ: ADMIN_UID_TOKEN_STATUS:{userID}
	key := fmt.Sprintf("ADMIN_UID_TOKEN_STATUS:%s", userID)

	// ɾ Hashʹû token ʧЧ
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
	total, err := d.appletCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
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
	// ϼܵСstatus=1
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
	// ʹ upsert »
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
	total, err := d.versionCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
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
