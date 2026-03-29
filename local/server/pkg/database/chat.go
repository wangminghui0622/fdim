package database

import (
	"context"
	"fmt"
	"time"

	"fdim/protocol/sdkws"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ChatDatabase Chatݿӿ
type ChatDatabase interface {
	// û˻
	FindUserAccount(ctx context.Context, accounts []*UserAccount) ([]*UserAccount, error)
	AddUserAccount(ctx context.Context, accounts []*UserAccount) error
	DelUserAccount(ctx context.Context, userIDs []string) error
	UpdateUserAccount(ctx context.Context, userID string, update map[string]interface{}) error
	GetUserAccountByUserID(ctx context.Context, userID string) (*UserAccount, error)
	GetUserAccountByAccount(ctx context.Context, account string) (*UserAccount, error)
	GetUserAccountByPhone(ctx context.Context, areaCode string, phoneNumber string) (*UserAccount, error)
	GetUserAccountByEmail(ctx context.Context, email string) (*UserAccount, error)

	// ûϢ
	FindUserPublicInfo(ctx context.Context, userIDs []string) ([]*UserPublicInfo, error)
	FindUserFullInfo(ctx context.Context, userIDs []string) ([]*UserFullInfo, error)
	SearchUserPublicInfo(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []*UserPublicInfo, error)
	SearchUserFullInfo(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []*UserFullInfo, error)
	AddUserInfo(ctx context.Context, info *UserFullInfo) error
	UpdateUserInfo(ctx context.Context, userID string, update map[string]interface{}) error

	// ֤
	AddVerifyCode(ctx context.Context, code *VerifyCode) error
	FindVerifyCode(ctx context.Context, phoneNumber string, areaCode string, verifyCode string) (*VerifyCode, error)
	DelVerifyCode(ctx context.Context, phoneNumber string, areaCode string) error
	FindVerifyCodeByEmail(ctx context.Context, email string, verifyCode string) (*VerifyCode, error)
	DelVerifyCodeByEmail(ctx context.Context, email string) error
	GetLastVerifyCode(ctx context.Context, phoneNumber string, areaCode string) (*VerifyCode, error)
	GetLastVerifyCodeByEmail(ctx context.Context, email string) (*VerifyCode, error)

	// û¼¼
	AddUserLoginRecord(ctx context.Context, record *UserLoginRecord) error
	UserLoginCount(ctx context.Context, startTime int64, endTime int64) (int64, error)
}

// UserAccount û˻ģ
type UserAccount struct {
	UserID      string    `bson:"user_id"`
	Account     string    `bson:"account"`
	Password    string    `bson:"password"`
	AreaCode    string    `bson:"area_code"`
	PhoneNumber string    `bson:"phone_number"`
	Email       string    `bson:"email"`
	CreateTime  time.Time `bson:"create_time"`
}

// UserPublicInfo ûϢģ
type UserPublicInfo struct {
	UserID   string `bson:"user_id"`
	Account  string `bson:"account"`
	Nickname string `bson:"nickname"`
	FaceURL  string `bson:"face_url"`
	Gender   int32  `bson:"gender"`
	Level    int32  `bson:"level"`
}

// UserFullInfo ûϢģ
type UserFullInfo struct {
	UserID           string    `bson:"user_id"`
	Account          string    `bson:"account"`
	PhoneNumber      string    `bson:"phone_number"`
	AreaCode         string    `bson:"area_code"`
	Email            string    `bson:"email"`
	Nickname         string    `bson:"nickname"`
	FaceURL          string    `bson:"face_url"`
	Gender           int32     `bson:"gender"`
	Level            int32     `bson:"level"`
	Birth            int64     `bson:"birth"`
	AllowAddFriend   int32     `bson:"allow_add_friend"`
	AllowBeep        int32     `bson:"allow_beep"`
	AllowVibration   int32     `bson:"allow_vibration"`
	GlobalRecvMsgOpt int32     `bson:"global_recv_msg_opt"`
	RegisterType     int32     `bson:"register_type"`
	CreateTime       time.Time `bson:"create_time"`
}

// VerifyCode ֤ģ
// ֻ˺ʽ
//   - ֻʹ phone_number + area_code
//   - 䣺ʹ email
//
// ֻ phone_number/area_code ֶΣֶ email ʹomitempty Աּݡ
type VerifyCode struct {
	PhoneNumber string    `bson:"phone_number,omitempty"`
	AreaCode    string    `bson:"area_code,omitempty"`
	Email       string    `bson:"email,omitempty"`
	Code        string    `bson:"code"`
	CreateTime  time.Time `bson:"create_time"`
	ExpireTime  time.Time `bson:"expire_time"`
}

// UserLoginRecord û¼¼ģ
type UserLoginRecord struct {
	UserID    string    `bson:"user_id"`
	LoginTime time.Time `bson:"login_time"`
	IP        string    `bson:"ip"`
	DeviceID  string    `bson:"device_id"`
	Platform  int32     `bson:"platform"`
}

// NewChatDatabase Chatݿʵ
func NewChatDatabase(mongoDB *MongoDB) ChatDatabase {
	return &chatDatabase{
		userAccountCollection: mongoDB.GetCollection("user_account"),
		userInfoCollection:    mongoDB.GetCollection("user_info"),
		verifyCodeCollection:  mongoDB.GetCollection("verify_code"),
		loginRecordCollection: mongoDB.GetCollection("user_login_record"),
	}
}

type chatDatabase struct {
	userAccountCollection *mongo.Collection
	userInfoCollection    *mongo.Collection
	verifyCodeCollection  *mongo.Collection
	loginRecordCollection *mongo.Collection
}

func (d *chatDatabase) FindUserAccount(ctx context.Context, accounts []*UserAccount) ([]*UserAccount, error) {
	if len(accounts) == 0 {
		return nil, nil
	}

	// ѯְ֧accountphone_number+area_codeemailѯ
	var filters []bson.M
	for _, acc := range accounts {
		filter := bson.M{}
		if acc.Account != "" {
			filter["account"] = acc.Account
		}
		if acc.PhoneNumber != "" && acc.AreaCode != "" {
			filter["phone_number"] = acc.PhoneNumber
			filter["area_code"] = acc.AreaCode
		}
		if acc.Email != "" {
			filter["email"] = acc.Email
		}
		if len(filter) > 0 {
			filters = append(filters, filter)
		}
	}

	if len(filters) == 0 {
		return nil, nil
	}

	filter := bson.M{"$or": filters}
	cursor, err := d.userAccountCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*UserAccount
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (d *chatDatabase) AddUserAccount(ctx context.Context, accounts []*UserAccount) error {
	if len(accounts) == 0 {
		return nil
	}

	docs := make([]interface{}, len(accounts))
	for i, acc := range accounts {
		docs[i] = acc
	}

	_, err := d.userAccountCollection.InsertMany(ctx, docs)
	return err
}

func (d *chatDatabase) DelUserAccount(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	_, err := d.userAccountCollection.DeleteMany(ctx, bson.M{"user_id": bson.M{"$in": userIDs}})
	return err
}

func (d *chatDatabase) UpdateUserAccount(ctx context.Context, userID string, update map[string]interface{}) error {
	_, err := d.userAccountCollection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": update})
	return err
}

func (d *chatDatabase) FindUserPublicInfo(ctx context.Context, userIDs []string) ([]*UserPublicInfo, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	filter := bson.M{"user_id": bson.M{"$in": userIDs}}
	cursor, err := d.userInfoCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*UserPublicInfo
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (d *chatDatabase) FindUserFullInfo(ctx context.Context, userIDs []string) ([]*UserFullInfo, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	filter := bson.M{"user_id": bson.M{"$in": userIDs}}
	cursor, err := d.userInfoCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*UserFullInfo
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (d *chatDatabase) SearchUserPublicInfo(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []*UserPublicInfo, error) {
	filter := bson.M{}
	if keyword != "" {
		filter["$or"] = []bson.M{
			{"user_id": bson.M{"$regex": keyword, "$options": "i"}},
			{"account": bson.M{"$regex": keyword, "$options": "i"}},
			{"nickname": bson.M{"$regex": keyword, "$options": "i"}},
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
	total, err := d.userInfoCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "create_time", Value: -1}})

	cursor, err := d.userInfoCollection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var results []*UserPublicInfo
	if err := cursor.All(ctx, &results); err != nil {
		return 0, nil, err
	}

	return total, results, nil
}

func (d *chatDatabase) SearchUserFullInfo(ctx context.Context, keyword string, pagination *sdkws.RequestPagination) (int64, []*UserFullInfo, error) {
	filter := bson.M{}
	if keyword != "" {
		filter["$or"] = []bson.M{
			{"user_id": bson.M{"$regex": keyword, "$options": "i"}},
			{"account": bson.M{"$regex": keyword, "$options": "i"}},
			{"nickname": bson.M{"$regex": keyword, "$options": "i"}},
			{"phone_number": bson.M{"$regex": keyword, "$options": "i"}},
			{"email": bson.M{"$regex": keyword, "$options": "i"}},
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
	total, err := d.userInfoCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	// ҳѯ
	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "create_time", Value: -1}})

	cursor, err := d.userInfoCollection.Find(ctx, filter, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var results []*UserFullInfo
	if err := cursor.All(ctx, &results); err != nil {
		return 0, nil, err
	}

	return total, results, nil
}

func (d *chatDatabase) AddVerifyCode(ctx context.Context, code *VerifyCode) error {
	_, err := d.verifyCodeCollection.InsertOne(ctx, code)
	return err
}

func (d *chatDatabase) FindVerifyCode(ctx context.Context, phoneNumber string, areaCode string, verifyCode string) (*VerifyCode, error) {
	filter := bson.M{
		"phone_number": phoneNumber,
		"area_code":    areaCode,
		"code":         verifyCode,
	}
	var result VerifyCode
	err := d.verifyCodeCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("verify code not found: %w", err)
		}
		return nil, err
	}
	return &result, nil
}

func (d *chatDatabase) DelVerifyCode(ctx context.Context, phoneNumber string, areaCode string) error {
	_, err := d.verifyCodeCollection.DeleteMany(ctx, bson.M{
		"phone_number": phoneNumber,
		"area_code":    areaCode,
	})
	return err
}

func (d *chatDatabase) FindVerifyCodeByEmail(ctx context.Context, email string, verifyCode string) (*VerifyCode, error) {
	filter := bson.M{
		"email": email,
		"code":  verifyCode,
	}
	var result VerifyCode
	err := d.verifyCodeCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("verify code not found: %w", err)
		}
		return nil, err
	}
	return &result, nil
}

func (d *chatDatabase) DelVerifyCodeByEmail(ctx context.Context, email string) error {
	_, err := d.verifyCodeCollection.DeleteMany(ctx, bson.M{
		"email": email,
	})
	return err
}

func (d *chatDatabase) GetLastVerifyCodeByEmail(ctx context.Context, email string) (*VerifyCode, error) {
	filter := bson.M{
		"email": email,
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "create_time", Value: -1}})
	var result VerifyCode
	err := d.verifyCodeCollection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("verify code not found: %w", err)
		}
		return nil, err
	}
	return &result, nil
}

func (d *chatDatabase) AddUserLoginRecord(ctx context.Context, record *UserLoginRecord) error {
	_, err := d.loginRecordCollection.InsertOne(ctx, record)
	return err
}

func (d *chatDatabase) UserLoginCount(ctx context.Context, startTime int64, endTime int64) (int64, error) {
	filter := bson.M{
		"login_time": bson.M{
			"$gte": time.Unix(startTime, 0),
			"$lte": time.Unix(endTime, 0),
		},
	}
	count, err := d.loginRecordCollection.CountDocuments(ctx, filter)
	return count, err
}

func (d *chatDatabase) GetUserAccountByUserID(ctx context.Context, userID string) (*UserAccount, error) {
	var result UserAccount
	err := d.userAccountCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user account not found: %w", err)
		}
		return nil, err
	}
	return &result, nil
}

func (d *chatDatabase) GetUserAccountByAccount(ctx context.Context, account string) (*UserAccount, error) {
	var result UserAccount
	err := d.userAccountCollection.FindOne(ctx, bson.M{"account": account}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user account not found: %w", err)
		}
		return nil, err
	}
	return &result, nil
}

func (d *chatDatabase) GetUserAccountByPhone(ctx context.Context, areaCode string, phoneNumber string) (*UserAccount, error) {
	var result UserAccount
	err := d.userAccountCollection.FindOne(ctx, bson.M{
		"area_code":    areaCode,
		"phone_number": phoneNumber,
	}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user account not found: %w", err)
		}
		return nil, err
	}
	return &result, nil
}

func (d *chatDatabase) GetUserAccountByEmail(ctx context.Context, email string) (*UserAccount, error) {
	var result UserAccount
	err := d.userAccountCollection.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user account not found: %w", err)
		}
		return nil, err
	}
	return &result, nil
}

func (d *chatDatabase) AddUserInfo(ctx context.Context, info *UserFullInfo) error {
	_, err := d.userInfoCollection.InsertOne(ctx, info)
	return err
}

func (d *chatDatabase) UpdateUserInfo(ctx context.Context, userID string, update map[string]interface{}) error {
	_, err := d.userInfoCollection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": update})
	return err
}

func (d *chatDatabase) GetLastVerifyCode(ctx context.Context, phoneNumber string, areaCode string) (*VerifyCode, error) {
	filter := bson.M{
		"phone_number": phoneNumber,
		"area_code":    areaCode,
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "create_time", Value: -1}})
	var result VerifyCode
	err := d.verifyCodeCollection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("verify code not found: %w", err)
		}
		return nil, err
	}
	return &result, nil
}
