package types

// 用户相关请求和响应类型
// 注意：这些类型需要与 open-im-server 的 API 接口保持一致

type UserRegisterReq struct {
	Secret   string   `json:"secret,optional"`
	UserID   string   `json:"userID,optional"`
	Nickname string   `json:"nickname,optional"`
	FaceURL  string   `json:"faceURL,optional"`
	Password string   `json:"password"` // 密码字段（必填）
}

type UserInfo struct {
	UserID   string `json:"userID,optional"`
	Nickname string `json:"nickname,optional"`
	FaceURL  string `json:"faceURL,optional"`
	Ex       string `json:"ex,optional"`
	Password string `json:"password"` // 添加密码字段（必填）
}

type UserRegisterResp struct {
	BaseResp
}

type UpdateUserInfoReq struct {
	UserInfo UserInfo `json:"userInfo"`
}

type UpdateUserInfoResp struct {
	BaseResp
}

type UpdateUserInfoExReq struct {
	UserInfo UserInfoEx `json:"userInfo"`
}

type UserInfoEx struct {
	UserID   string `json:"userID"`
	Nickname *string `json:"nickname,optional"`
	FaceURL  *string `json:"faceURL,optional"`
	Ex       *string `json:"ex,optional"`
}

type UpdateUserInfoExResp struct {
	BaseResp
}

type GetUsersInfoReq struct {
	UserIDs []string `json:"userIDs"`
}

type GetUsersInfoResp struct {
	BaseResp
	UsersData []UserInfo `json:"usersData"`
}

type AccountCheckReq struct {
	CheckUserIDs []string `json:"checkUserIDs"`
}

type AccountCheckResp struct {
	BaseResp
	Results []AccountCheckResult `json:"results"`
}

type AccountCheckResult struct {
	UserID string `json:"userID"`
	Result int32  `json:"result"` // 0: 不存在, 1: 存在
}

type SetGlobalRecvMessageOptReq struct {
	UserID         string `json:"userID"`
	GlobalRecvMsgOpt int32  `json:"globalRecvMsgOpt"` // 0: 接收, 1: 不接收, 2: 接收但不提醒
}

type SetGlobalRecvMessageOptResp struct {
	BaseResp
}

type GetAllUsersIDReq struct {
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	PageNumber int32 `json:"pageNumber"`
	ShowNumber int32 `json:"showNumber"`
}

type GetAllUsersIDResp struct {
	BaseResp
	Total  int32    `json:"total"`
	UserIDs []string `json:"userIDs"`
}

type GetUsersReq struct {
	Pagination Pagination `json:"pagination"`
	UserID     string     `json:"userID,optional"`
	Nickname   string     `json:"nickname,optional"`
}

type GetUsersResp struct {
	BaseResp
	Total int32      `json:"total"`
	Users []UserInfo `json:"users"`
}

type UserRegisterCountReq struct {
	Start int64 `json:"start"` // 开始时间戳（秒）
	End   int64 `json:"end"`   // 结束时间戳（秒）
}

type UserRegisterCountResp struct {
	BaseResp
	Total  int64 `json:"total"`
	Before int64 `json:"before"`
}

type BaseResp struct {
	ErrCode int32  `json:"-"`
	ErrMsg  string `json:"-"`
	ErrDlt  string `json:"-"`
}

// GetUserStatus 相关类型
type GetUserStatusReq struct {
	UserID  string   `json:"userID,optional"`
	UserIDs []string `json:"userIDs,optional"`
}

type OnlineStatus struct {
	UserID      string  `json:"userID"`
	Status      int32   `json:"status"` // 0: 离线, 1: 在线
	PlatformIDs []int32 `json:"platformIDs,optional"`
}

type GetUserStatusResp struct {
	BaseResp
	StatusList []OnlineStatus `json:"statusList"`
}

// SubscribeOrCancelUsersStatus 相关类型
type SubscribeOrCancelUsersStatusReq struct {
	UserID     string   `json:"userID"`
	UserIDs    []string `json:"userIDs"`
	Genre      int32    `json:"genre"` // 0: 订阅, 1: 取消订阅
}

type SubscribeOrCancelUsersStatusResp struct {
	BaseResp
}

// GetSubscribeUsersStatus 相关类型
type GetSubscribeUsersStatusReq struct {
	UserID string `json:"userID"`
}

type GetSubscribeUsersStatusResp struct {
	BaseResp
	StatusList []OnlineStatus `json:"statusList"`
}

// ProcessUserCommand 相关类型
type ProcessUserCommandAddReq struct {
	UserID string  `json:"userID"`
	Type   int32   `json:"type"`
	Uuid   string  `json:"uuid"`
	Value  *string `json:"value,optional"`
	Ex     *string `json:"ex,optional"`
}

type ProcessUserCommandAddResp struct {
	BaseResp
}

type ProcessUserCommandDeleteReq struct {
	UserID string `json:"userID"`
	Type   int32  `json:"type"`
	Uuid   string `json:"uuid"`
}

type ProcessUserCommandDeleteResp struct {
	BaseResp
}

type ProcessUserCommandUpdateReq struct {
	UserID string  `json:"userID"`
	Type   int32   `json:"type"`
	Uuid   string  `json:"uuid"`
	Value  *string `json:"value,optional"`
	Ex     *string `json:"ex,optional"`
}

type ProcessUserCommandUpdateResp struct {
	BaseResp
}

type ProcessUserCommandGetReq struct {
	UserID string `json:"userID"`
	Type   int32  `json:"type"`
}

type ProcessUserCommandGetResp struct {
	BaseResp
	Commands []UserCommand `json:"commands"`
}

type ProcessUserCommandGetAllReq struct {
	UserID string `json:"userID"`
}

type UserCommand struct {
	Uuid  string `json:"uuid"`
	Value string `json:"value,optional"`
	Ex    string `json:"ex,optional"`
}

type ProcessUserCommandGetAllResp struct {
	BaseResp
	Commands []UserCommand `json:"commands"`
}

// NotificationAccount 相关类型
type AddNotificationAccountReq struct {
	UserID        string `json:"userID"`
	NickName      string `json:"nickName"`
	FaceURL       string `json:"faceURL"`
	AppMangerLevel int32  `json:"appMangerLevel"`
}

type AddNotificationAccountResp struct {
	BaseResp
	UserID        string `json:"userID"`
	FaceURL       string `json:"faceURL"`
	NickName      string `json:"nickName"`
	AppMangerLevel int32  `json:"appMangerLevel"`
}

type UpdateNotificationAccountInfoReq struct {
	UserID   string `json:"userID"`
	FaceURL  string `json:"faceURL,optional"`
	NickName string `json:"nickName,optional"`
}

type UpdateNotificationAccountInfoResp struct {
	BaseResp
}

type SearchNotificationAccountReq struct {
	Keyword        string     `json:"keyword"`
	AppManagerLevel *int32     `json:"appManagerLevel,optional"`
	Pagination     Pagination `json:"pagination"`
}

type NotificationAccountInfo struct {
	UserID        string `json:"userID"`
	FaceURL       string `json:"faceURL"`
	NickName      string `json:"nickName"`
	AppMangerLevel int32  `json:"appMangerLevel"`
}

type SearchNotificationAccountResp struct {
	BaseResp
	Total                int64                   `json:"total"`
	NotificationAccounts []NotificationAccountInfo `json:"notificationAccounts"`
}

// UserClientConfig 相关类型
type GetUserClientConfigReq struct {
	UserID string `json:"userID"`
}

type GetUserClientConfigResp struct {
	BaseResp
	Configs map[string]string `json:"configs"`
}

type SetUserClientConfigReq struct {
	UserID  string            `json:"userID"`
	Configs map[string]string `json:"configs"`
}

type SetUserClientConfigResp struct {
	BaseResp
}

type DelUserClientConfigReq struct {
	UserID string   `json:"userID"`
	Keys   []string `json:"keys"`
}

type DelUserClientConfigResp struct {
	BaseResp
}

type PageUserClientConfigReq struct {
	UserID     string     `json:"userID"`
	Key        string     `json:"key,optional"`
	Pagination Pagination `json:"pagination"`
}

type ClientConfig struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PageUserClientConfigResp struct {
	BaseResp
	Total   int64          `json:"total"`
	Configs []ClientConfig `json:"configs"`
}

// GetUsersOnlineStatus 相关类型（需要 msggateway 服务）
type GetUsersOnlineStatusReq struct {
	UserIDs []string `json:"userIDs"`
}

type UserOnlineStatus struct {
	UserID               string                      `json:"userID"`
	Status               int32                       `json:"status"`
	DetailPlatformStatus []UserOnlinePlatformStatus  `json:"detailPlatformStatus,omitempty"`
}

type UserOnlinePlatformStatus struct {
	Platform string `json:"platform"`
	Status   int32  `json:"status"`
}

type GetUsersOnlineStatusResp struct {
	BaseResp
	SuccessResult []UserOnlineStatus `json:"successResult"`
	FailedResult  []interface{}      `json:"failedResult,omitempty"`
}

// GetUsersOnlineTokenDetail 相关类型（需要 msggateway 服务）
type GetUsersOnlineTokenDetailReq struct {
	UserIDs []string `json:"userIDs"`
}

type GetUsersOnlineTokenDetailResp struct {
	BaseResp
	Details []interface{} `json:"details"` // 使用 SingleDetail
}

// SearchUserFullInfo 相关类型
type SearchUserFullInfoReq struct {
	Keyword    string     `json:"keyword"`
	Pagination Pagination `json:"pagination,optional"`
}

type UserFullInfo struct {
	UserID           string `json:"userID"`
	Nickname         string `json:"nickname"`
	FaceURL          string `json:"faceURL"`
	PhoneNumber      string `json:"phoneNumber,optional"`
	Email            string `json:"email,optional"`
	Gender           int32  `json:"gender,optional"`
	Birth            int64  `json:"birth,optional"`
	Ex               string `json:"ex,optional"`
	CreateTime       int64  `json:"createTime,optional"`
	AppMangerLevel   int32  `json:"appMangerLevel,optional"`
	GlobalRecvMsgOpt int32  `json:"globalRecvMsgOpt,optional"`
}

type SearchUserFullInfoResp struct {
	Total int32          `json:"total"`
	Users []UserFullInfo `json:"users"`
}

// ========== Account 登录/注册 ==========

type AccountLoginReq struct {
	PhoneNumber string `json:"phoneNumber,optional"`
	AreaCode    string `json:"areaCode,optional"`
	Email       string `json:"email,optional"`
	Account     string `json:"account,optional"`
	Password    string `json:"password,optional"`
	VerifyCode  string `json:"verifyCode,optional"`
	Platform    int32  `json:"platform"`
	DeviceID    string `json:"deviceID,optional"`
	Ip          string `json:"ip,optional"`
}

type AccountLoginResp struct {
	UserID    string `json:"userID"`
	ImToken   string `json:"imToken"`
	ChatToken string `json:"chatToken"`
}

type AccountRegisterReq struct {
	VerifyCode     string            `json:"verifyCode,optional"`
	Platform       int32             `json:"platform"`
	DeviceID       string            `json:"deviceID,optional"`
	InvitationCode string            `json:"invitationCode,optional"`
	AutoLogin      bool              `json:"autoLogin,optional"`
	User           *RegisterUserInfo `json:"user"`
}

type RegisterUserInfo struct {
	Nickname    string `json:"nickname"`
	FaceURL     string `json:"faceURL,optional"`
	AreaCode    string `json:"areaCode,optional"`
	PhoneNumber string `json:"phoneNumber,optional"`
	Email       string `json:"email,optional"`
	Account     string `json:"account,optional"`
	Password    string `json:"password"`
	Birth       int64  `json:"birth,optional"`
	Gender      int32  `json:"gender,optional"`
}

type AccountRegisterResp struct {
	UserID    string `json:"userID"`
	ImToken   string `json:"imToken"`
	ChatToken string `json:"chatToken"`
}

// ========== Account 修改/重置密码 ==========

type ChangePasswordReq struct {
	UserID          string `json:"userID"`
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
	Platform        int32  `json:"platform,optional"`
}

type ChangePasswordResp struct{}

type ResetPasswordReq struct {
	PhoneNumber string `json:"phoneNumber,optional"`
	AreaCode    string `json:"areaCode,optional"`
	Email       string `json:"email,optional"`
	VerifyCode  string `json:"verifyCode"`
	NewPassword string `json:"newPassword"`
	Platform    int32  `json:"platform,optional"`
}

type ResetPasswordResp struct{}

// ========== Friend 搜索 ==========

type SearchFriendReq struct {
	UserID  string `json:"userID"`
	Keyword string `json:"keyword"`
}

type SearchFriendResp struct {
	Friends interface{} `json:"friends"`
}

// ========== 验证码 ==========

type SendVerifyCodeReq struct {
	UsedFor     int32  `json:"usedFor"`
	PhoneNumber string `json:"phoneNumber,optional"`
	AreaCode    string `json:"areaCode,optional"`
	Email       string `json:"email,optional"`
	Ip          string `json:"ip,optional"`
	Platform    int32  `json:"platform,optional"`
	DeviceID    string `json:"deviceID,optional"`
}

type SendVerifyCodeResp struct{}

type VerifyCodeReq struct {
	PhoneNumber string `json:"phoneNumber,optional"`
	AreaCode    string `json:"areaCode,optional"`
	Email       string `json:"email,optional"`
	VerifyCode  string `json:"verifyCode"`
}

type VerifyCodeResp struct{}

// ========== Chat 层用户接口 ==========

type FindUserFullInfoReq struct {
	UserIDs []string `json:"userIDs"`
}

type FindUserFullInfoResp struct {
	Users interface{} `json:"users"`
}

type ChatUpdateUserInfoReq struct {
	UserID      string  `json:"userID"`
	Account     *string `json:"account,optional"`
	PhoneNumber *string `json:"phoneNumber,optional"`
	AreaCode    *string `json:"areaCode,optional"`
	Email       *string `json:"email,optional"`
	Nickname    *string `json:"nickname,optional"`
	FaceURL     *string `json:"faceURL,optional"`
	Gender      *int32  `json:"gender,optional"`
	Birth       *int64  `json:"birth,optional"`
}

type ChatUpdateUserInfoResp struct{}

// ========== 客户端配置 ==========

type GetClientConfigReq struct{}

type GetClientConfigResp struct {
	Config interface{} `json:"config"`
}

// ========== Chat 层: user/find/public, user/search/public ==========

type FindUserPublicInfoReq struct {
	UserIDs []string `json:"userIDs"`
}

type FindUserPublicInfoResp struct {
	Users interface{} `json:"users"`
}

type SearchUserPublicInfoReq struct {
	Keyword    string     `json:"keyword"`
	Pagination Pagination `json:"pagination,optional"`
}

type SearchUserPublicInfoResp struct {
	Total int32       `json:"total"`
	Users interface{} `json:"users"`
}

// ========== RTC Token ==========

type GetTokenForVideoMeetingReq struct {
	Room     string `json:"room"`
	Identity string `json:"identity"`
}

type GetTokenForVideoMeetingResp struct {
	ServerUrl string `json:"serverUrl"`
	Token     string `json:"token"`
}

// ========== Applet ==========

type FindAppletReq struct{}

type AppletInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	AppID      string `json:"appID"`
	Icon       string `json:"icon"`
	Url        string `json:"url"`
	Md5        string `json:"md5"`
	Size       int64  `json:"size"`
	Version    string `json:"version"`
	Priority   int32  `json:"priority"`
	Status     int32  `json:"status"`
	CreateTime int64  `json:"createTime"`
}

type FindAppletResp struct {
	Applets []AppletInfo `json:"applets"`
}

type AddAppletReq struct {
	Name     string `json:"name"`
	AppID    string `json:"appID"`
	Icon     string `json:"icon"`
	Url      string `json:"url"`
	Md5      string `json:"md5"`
	Size     int64  `json:"size"`
	Version  string `json:"version"`
	Priority int32  `json:"priority"`
}

type AddAppletResp struct{}

type DelAppletReq struct {
	AppletIDs []string `json:"appletIDs"`
}

type DelAppletResp struct{}

type UpdateAppletReq struct {
	ID       string  `json:"id"`
	Name     *string `json:"name,optional"`
	AppID    *string `json:"appID,optional"`
	Icon     *string `json:"icon,optional"`
	Url      *string `json:"url,optional"`
	Md5      *string `json:"md5,optional"`
	Size     *int64  `json:"size,optional"`
	Version  *string `json:"version,optional"`
	Priority *int32  `json:"priority,optional"`
}

type UpdateAppletResp struct{}

type SearchAppletReq struct {
	Keyword    string     `json:"keyword,optional"`
	Pagination Pagination `json:"pagination"`
}

type SearchAppletResp struct {
	Total   int64        `json:"total"`
	Applets []AppletInfo `json:"applets"`
}

// ========== Application Version ==========

type ApplicationVersion struct {
	ID         string `json:"id"`
	Platform   string `json:"platform"`
	Version    string `json:"version"`
	Url        string `json:"url"`
	Text       string `json:"text"`
	ForceID    int64  `json:"forceID"`
	IsForce    bool   `json:"isForce"`
	Latest     bool   `json:"latest"`
	Hot        bool   `json:"hot"`
	CreateTime int64  `json:"createTime"`
}

type LatestApplicationVersionReq struct {
	Platform string `json:"platform"`
}

type LatestApplicationVersionResp struct {
	Version *ApplicationVersion `json:"version"`
}

type PageApplicationVersionReq struct {
	Platform   string     `json:"platform,optional"`
	Pagination Pagination `json:"pagination"`
}

type PageApplicationVersionResp struct {
	Total    int64                `json:"total"`
	Versions []ApplicationVersion `json:"versions"`
}

type AddApplicationVersionReq struct {
	Platform string `json:"platform"`
	Version  string `json:"version"`
	Url      string `json:"url"`
	Text     string `json:"text"`
	IsForce  bool   `json:"isForce"`
	Latest   bool   `json:"latest"`
	Hot      bool   `json:"hot"`
}

type AddApplicationVersionResp struct{}

type UpdateApplicationVersionReq struct {
	ID       string  `json:"id"`
	Platform *string `json:"platform,optional"`
	Version  *string `json:"version,optional"`
	Url      *string `json:"url,optional"`
	Text     *string `json:"text,optional"`
	IsForce  *bool   `json:"isForce,optional"`
	Latest   *bool   `json:"latest,optional"`
	Hot      *bool   `json:"hot,optional"`
}

type UpdateApplicationVersionResp struct{}

type DeleteApplicationVersionReq struct {
	ID string `json:"id"`
}

type DeleteApplicationVersionResp struct{}

// ========== OpenIM Callback ==========

type OpenIMCallbackReq struct {
	Command     string      `json:"command"`
	Body        interface{} `json:"body"`
	CallbackKey string      `json:"callbackKey,optional"`
}

type OpenIMCallbackResp struct{}

// ========== Admin ==========

type AdminLoginReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type AdminLoginResp struct {
	AdminToken  string `json:"adminToken"`
	AdminUserID string `json:"adminUserID"`
	Nickname    string `json:"nickname"`
	FaceURL     string `json:"faceURL"`
	Level       int32  `json:"level"`
}

type AdminUpdateInfoReq struct {
	Account  *string `json:"account,optional"`
	Password *string `json:"password,optional"`
	FaceURL  *string `json:"faceURL,optional"`
	Nickname *string `json:"nickname,optional"`
	Level    *int32  `json:"level,optional"`
}

type AdminUpdateInfoResp struct{}

type AdminInfoReq struct{}

type AdminInfoResp struct {
	Account  string `json:"account"`
	Nickname string `json:"nickname"`
	FaceURL  string `json:"faceURL"`
	Level    int32  `json:"level"`
}

type ChangeAdminPasswordReq struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type ChangeAdminPasswordResp struct{}

type AddAdminAccountReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	FaceURL  string `json:"faceURL,optional"`
	Nickname string `json:"nickname,optional"`
	Level    int32  `json:"level,optional"`
}

type AddAdminAccountResp struct{}

type DelAdminAccountReq struct {
	UserIDs []string `json:"userIDs"`
}

type DelAdminAccountResp struct{}

type SearchAdminAccountReq struct {
	Keyword    string     `json:"keyword,optional"`
	Pagination Pagination `json:"pagination"`
}

type AdminAccountInfo struct {
	UserID   string `json:"userID"`
	Account  string `json:"account"`
	Nickname string `json:"nickname"`
	FaceURL  string `json:"faceURL"`
	Level    int32  `json:"level"`
}

type SearchAdminAccountResp struct {
	Total  int64              `json:"total"`
	Admins []AdminAccountInfo `json:"admins"`
}

type AddUserAccountReq struct {
	User RegisterUserInfo `json:"user"`
}

type AddUserAccountResp struct {
	UserID string `json:"userID"`
}

// ========== Admin: User Import ==========

type ImportUserByJsonReq struct {
	Users []RegisterUserInfo `json:"users"`
}

type ImportUserByJsonResp struct {
	UserIDs []string `json:"userIDs"`
}

// ========== Admin: Allow Register ==========

type GetAllowRegisterReq struct{}

type GetAllowRegisterResp struct {
	AllowRegister bool `json:"allowRegister"`
}

type SetAllowRegisterReq struct {
	AllowRegister bool `json:"allowRegister"`
}

type SetAllowRegisterResp struct{}

// ========== Admin: Default Friend ==========

type AddDefaultFriendReq struct {
	UserIDs []string `json:"userIDs"`
}

type AddDefaultFriendResp struct{}

type DelDefaultFriendReq struct {
	UserIDs []string `json:"userIDs"`
}

type DelDefaultFriendResp struct{}

type FindDefaultFriendReq struct{}

type FindDefaultFriendResp struct {
	UserIDs []string `json:"userIDs"`
}

type SearchDefaultFriendReq struct {
	Keyword    string     `json:"keyword,optional"`
	Pagination Pagination `json:"pagination"`
}

type SearchDefaultFriendResp struct {
	Total   int64    `json:"total"`
	UserIDs []string `json:"userIDs"`
}

// ========== Admin: Default Group ==========

type AddDefaultGroupReq struct {
	GroupIDs []string `json:"groupIDs"`
}

type AddDefaultGroupResp struct{}

type DelDefaultGroupReq struct {
	GroupIDs []string `json:"groupIDs"`
}

type DelDefaultGroupResp struct{}

type FindDefaultGroupReq struct{}

type FindDefaultGroupResp struct {
	GroupIDs []string `json:"groupIDs"`
}

type SearchDefaultGroupReq struct {
	Keyword    string     `json:"keyword,optional"`
	Pagination Pagination `json:"pagination"`
}

type SearchDefaultGroupResp struct {
	Total    int64    `json:"total"`
	GroupIDs []string `json:"groupIDs"`
}

// ========== Admin: Invitation Code ==========

type AddInvitationCodeReq struct {
	Codes []string `json:"codes"`
}

type AddInvitationCodeResp struct{}

type GenInvitationCodeReq struct {
	Len int32 `json:"len"`
	Num int32 `json:"num"`
}

type GenInvitationCodeResp struct {
	Codes []string `json:"codes"`
}

type DelInvitationCodeReq struct {
	Codes []string `json:"codes"`
}

type DelInvitationCodeResp struct{}

type SearchInvitationCodeReq struct {
	Keyword    string     `json:"keyword,optional"`
	Status     *int32     `json:"status,optional"`
	Pagination Pagination `json:"pagination"`
}

type InvitationCodeInfo struct {
	Code       string `json:"code"`
	CreateTime int64  `json:"createTime"`
	UsedUserID string `json:"usedUserID"`
	Status     int32  `json:"status"`
}

type SearchInvitationCodeResp struct {
	Total int64                `json:"total"`
	Codes []InvitationCodeInfo `json:"codes"`
}

// ========== Admin: IP Forbidden ==========

type AddIPForbiddenReq struct {
	Ips []string `json:"ips"`
}

type AddIPForbiddenResp struct{}

type DelIPForbiddenReq struct {
	Ips []string `json:"ips"`
}

type DelIPForbiddenResp struct{}

type SearchIPForbiddenReq struct {
	Keyword    string     `json:"keyword,optional"`
	Pagination Pagination `json:"pagination"`
}

type IPForbiddenInfo struct {
	Ip         string `json:"ip"`
	LimitLogin bool   `json:"limitLogin"`
	LimitRegister bool `json:"limitRegister"`
	CreateTime int64  `json:"createTime"`
}

type SearchIPForbiddenResp struct {
	Total int64             `json:"total"`
	Ips   []IPForbiddenInfo `json:"ips"`
}

// ========== Admin: User IP Limit ==========

type AddUserIPLimitLoginReq struct {
	UserID string   `json:"userID"`
	Ips    []string `json:"ips"`
}

type AddUserIPLimitLoginResp struct{}

type DelUserIPLimitLoginReq struct {
	UserID string   `json:"userID"`
	Ips    []string `json:"ips"`
}

type DelUserIPLimitLoginResp struct{}

type SearchUserIPLimitLoginReq struct {
	Keyword    string     `json:"keyword,optional"`
	Pagination Pagination `json:"pagination"`
}

type UserIPLimitInfo struct {
	UserID     string `json:"userID"`
	Ip         string `json:"ip"`
	CreateTime int64  `json:"createTime"`
}

type SearchUserIPLimitLoginResp struct {
	Total  int64             `json:"total"`
	Limits []UserIPLimitInfo `json:"limits"`
}

// ========== Admin: Block User ==========

type BlockUserReq struct {
	UserID string `json:"userID"`
	Reason string `json:"reason,optional"`
}

type BlockUserResp struct{}

type UnblockUserReq struct {
	UserIDs []string `json:"userIDs"`
}

type UnblockUserResp struct{}

type SearchBlockUserReq struct {
	Keyword    string     `json:"keyword,optional"`
	Pagination Pagination `json:"pagination"`
}

type BlockUserInfo struct {
	UserID     string `json:"userID"`
	Reason     string `json:"reason"`
	OpUserID   string `json:"opUserID"`
	CreateTime int64  `json:"createTime"`
}

type SearchBlockUserResp struct {
	Total int64           `json:"total"`
	Users []BlockUserInfo `json:"users"`
}

// ========== Admin: Reset User Password ==========

type AdminResetUserPasswordReq struct {
	UserID      string `json:"userID"`
	NewPassword string `json:"newPassword"`
}

type AdminResetUserPasswordResp struct{}

// ========== Admin: Client Config Set/Del ==========

type SetClientConfigReq struct {
	Config map[string]string `json:"config"`
}

type SetClientConfigResp struct{}

type DelClientConfigReq struct {
	Keys []string `json:"keys"`
}

type DelClientConfigResp struct{}

// ========== Admin: Statistics ==========

type NewUserCountReq struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type DateCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type NewUserCountResp struct {
	Total      int64       `json:"total"`
	DateCounts []DateCount `json:"dateCounts"`
}

type LoginUserCountReq struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type LoginUserCountResp struct {
	Total      int64       `json:"total"`
	DateCounts []DateCount `json:"dateCounts"`
}

// ========== Statistics (open-im-server) ==========

type GetActiveUserReq struct {
	Start int64  `json:"start"`
	End   int64  `json:"end"`
	Ase   bool   `json:"ase,optional"`
}

type ActiveUser struct {
	UserID       string `json:"userID"`
	Nickname     string `json:"nickname"`
	MessageCount int64  `json:"messageCount"`
}

type GetActiveUserResp struct {
	BaseResp
	MsgCount int64        `json:"msgCount"`
	Users    []ActiveUser `json:"users"`
}

type GroupCreateCountReq struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type GroupCreateCountResp struct {
	BaseResp
	Total  int64       `json:"total"`
	Before int64       `json:"before"`
	Counts []DateCount `json:"counts"`
}

type GetActiveGroupReq struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
	Ase   bool  `json:"ase,optional"`
}

type ActiveGroup struct {
	GroupID      string `json:"groupID"`
	GroupName    string `json:"groupName"`
	MessageCount int64  `json:"messageCount"`
}

type GetActiveGroupResp struct {
	BaseResp
	MsgCount int64         `json:"msgCount"`
	Groups   []ActiveGroup `json:"groups"`
}

// ========== JSSDK ==========

type JSSdkGetConversationsReq struct {
	UserID         string     `json:"userID"`
	ConversationIDs []string  `json:"conversationIDs,optional"`
	Pagination     Pagination `json:"pagination,optional"`
}

type JSSdkGetConversationsResp struct {
	Total         int64       `json:"total"`
	Conversations interface{} `json:"conversations"`
}

type JSSdkGetActiveConversationsReq struct {
	UserID string `json:"userID"`
	Count  int32  `json:"count"`
}

type JSSdkGetActiveConversationsResp struct {
	Conversations interface{} `json:"conversations"`
}

// ========== Config Management ==========

type GetConfigListReq struct{}

type ConfigItem struct {
	Name    string `json:"name"`
	Version int64  `json:"version"`
}

type GetConfigListResp struct {
	Configs []ConfigItem `json:"configs"`
}

type GetConfigReq struct {
	Name string `json:"name"`
}

type GetConfigResp struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Version int64  `json:"version"`
}

type SetConfigReq struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type SetConfigResp struct{}

type ResetConfigReq struct {
	Name string `json:"name"`
}

type ResetConfigResp struct{}

type SetEnableConfigManagerReq struct {
	Enable bool `json:"enable"`
}

type SetEnableConfigManagerResp struct{}

type GetEnableConfigManagerReq struct{}

type GetEnableConfigManagerResp struct {
	Enable bool `json:"enable"`
}
