package types

// Friend 相关请求和响应类?

// ApplyToAddFriend 相关类型
type ApplyToAddFriendReq struct {
	FromUserID string `json:"fromUserID"`
	ToUserID   string `json:"toUserID"`
	ReqMsg     string `json:"reqMsg,optional"`
	Ex         string `json:"ex,optional"`
}

type ApplyToAddFriendResp struct {
	BaseResp
}

// RespondFriendApply 相关类型
type RespondFriendApplyReq struct {
	FromUserID   string `json:"fromUserID"`
	ToUserID     string `json:"toUserID"`
	HandleResult int32  `json:"handleResult"` // 0: 拒绝, 1: 同意
	HandleMsg    string `json:"handleMsg,optional"`
}

type RespondFriendApplyResp struct {
	BaseResp
}

// DeleteFriend 相关类型
type DeleteFriendReq struct {
	OwnerUserID  string `json:"ownerUserID"`
	FriendUserID string `json:"friendUserID"`
}

type DeleteFriendResp struct {
	BaseResp
}

// GetFriendApplyList 相关类型
type GetFriendApplyListReq struct {
	UserID     string     `json:"userID"`
	Pagination Pagination `json:"pagination"`
}

type FriendRequest struct {
	FromUserID   string `json:"fromUserID"`
	ToUserID     string `json:"toUserID"`
	HandleResult int32  `json:"handleResult"`
	ReqMsg       string `json:"reqMsg,optional"`
	HandleMsg    string `json:"handleMsg,optional"`
	CreateTime   int64  `json:"createTime"`
	HandleTime   int64  `json:"handleTime"`
	Ex           string `json:"ex,optional"`
}

type GetFriendApplyListResp struct {
	BaseResp
	FriendRequests []FriendRequest `json:"friendRequests"`
	Total          int32           `json:"total"`
}

// GetDesignatedFriendsApply 相关类型
type GetDesignatedFriendsApplyReq struct {
	FromUserID string `json:"fromUserID"`
	ToUserID   string `json:"toUserID"`
}

type GetDesignatedFriendsApplyResp struct {
	BaseResp
	FriendRequests []FriendRequest `json:"friendRequests"`
}

// GetSelfApplyList 相关类型
type GetSelfApplyListReq struct {
	UserID     string     `json:"userID"`
	Pagination Pagination `json:"pagination"`
}

type GetSelfApplyListResp struct {
	BaseResp
	FriendRequests []FriendRequest `json:"friendRequests"`
	Total          int32           `json:"total"`
}

// GetFriendList 相关类型
type GetFriendListReq struct {
	UserID     string     `json:"userID"`
	Pagination Pagination `json:"pagination"`
}

type FriendInfo struct {
	OwnerUserID       string  `json:"ownerUserID"`
	FriendUserID      string  `json:"friendUserID"`
	Remark            string  `json:"remark,optional"`
	CreateTime        int64   `json:"createTime"`
	AddSource         int32   `json:"addSource"`
	Ex                string  `json:"ex,optional"`
	Nickname          string  `json:"nickname,optional"`
	FaceURL           string  `json:"faceURL,optional"`
	OnlinePlatformIDs []int32 `json:"onlinePlatformIDs,optional"` // 在线平台ID列表
}

type GetFriendListResp struct {
	BaseResp
	FriendsInfo []FriendInfo `json:"friendsInfo"`
	Total       int32        `json:"total"`
}

// GetDesignatedFriends 相关类型
type GetDesignatedFriendsReq struct {
	OwnerUserID   string   `json:"ownerUserID"`
	FriendUserIDs []string `json:"friendUserIDs"`
}

type GetDesignatedFriendsResp struct {
	BaseResp
	FriendsInfo []FriendInfo `json:"friendsInfo"`
}

// SetFriendRemark 相关类型
type SetFriendRemarkReq struct {
	OwnerUserID  string `json:"ownerUserID"`
	FriendUserID string `json:"friendUserID"`
	Remark       string `json:"remark,optional"`
}

type SetFriendRemarkResp struct {
	BaseResp
}

// AddBlack 相关类型
type AddBlackReq struct {
	OwnerUserID string `json:"ownerUserID"`
	BlackUserID string `json:"blackUserID"`
}

type AddBlackResp struct {
	BaseResp
}

// GetPaginationBlacks 相关类型
type GetPaginationBlacksReq struct {
	UserID     string     `json:"userID"`
	Pagination Pagination `json:"pagination"`
}

type BlackInfo struct {
	OwnerUserID string `json:"ownerUserID"`
	BlackUserID string `json:"blackUserID"`
	CreateTime  int64  `json:"createTime"`
	AddSource   int32  `json:"addSource"`
	Ex          string `json:"ex,optional"`
}

type GetPaginationBlacksResp struct {
	BaseResp
	Blacks []BlackInfo `json:"blacks"`
	Total  int32       `json:"total"`
}

// GetSpecifiedBlacks 相关类型
type GetSpecifiedBlacksReq struct {
	OwnerUserID string   `json:"ownerUserID"`
	UserIDList  []string `json:"userIDList"`
}

type GetSpecifiedBlacksResp struct {
	BaseResp
	Blacks []BlackInfo `json:"blacks"`
}

// RemoveBlack 相关类型
type RemoveBlackReq struct {
	OwnerUserID string `json:"ownerUserID"`
	BlackUserID string `json:"blackUserID"`
}

type RemoveBlackResp struct {
	BaseResp
}

// GetIncrementalBlacks 相关类型
type GetIncrementalBlacksReq struct {
	UserID    string `json:"userID"`
	VersionID string `json:"versionID,optional"`
	Version   uint64 `json:"version"`
}

type GetIncrementalBlacksResp struct {
	BaseResp
	Insert    []BlackInfo `json:"insert"`
	Update    []BlackInfo `json:"update"`
	Delete    []string    `json:"delete"`
	Version   uint64      `json:"version"`
	VersionID string      `json:"versionID"`
	Full      bool        `json:"full"`
}

// ImportFriends 相关类型
type ImportFriendsReq struct {
	OwnerUserID   string   `json:"ownerUserID"`
	FriendUserIDs []string `json:"friendUserIDs"`
}

type ImportFriendsResp struct {
	BaseResp
}

// IsFriend 相关类型
type IsFriendReq struct {
	UserID1 string `json:"userID1"`
	UserID2 string `json:"userID2"`
}

type IsFriendResp struct {
	BaseResp
	InUser1Friends bool `json:"inUser1Friends"`
	InUser2Friends bool `json:"inUser2Friends"`
}

// GetFriendIDs 相关类型
type GetFriendIDsReq struct {
	UserID string `json:"userID"`
}

type GetFriendIDsResp struct {
	BaseResp
	FriendIDs []string `json:"friendIDs"`
}

// GetSpecifiedFriendsInfo 相关类型
type GetSpecifiedFriendsInfoReq struct {
	OwnerUserID string   `json:"ownerUserID"`
	UserIDList  []string `json:"userIDList"`
}

type GetSpecifiedFriendsInfoResp struct {
	BaseResp
	FriendsInfo []FriendInfo `json:"friendsInfo"`
}

// UpdateFriends 相关类型
type UpdateFriendsReq struct {
	OwnerUserID   string   `json:"ownerUserID"`
	FriendUserIDs []string `json:"friendUserIDs"`
}

type UpdateFriendsResp struct {
	BaseResp
}

// GetIncrementalFriends 相关类型
type GetIncrementalFriendsReq struct {
	UserID    string `json:"userID"`
	VersionID string `json:"versionID,optional"`
	Version   uint64 `json:"version"`
}

type GetIncrementalFriendsResp struct {
	BaseResp
	Insert    []FriendInfo `json:"insert"`
	Update    []FriendInfo `json:"update"`
	Delete    []string     `json:"delete"`
	Version   uint64       `json:"version"`
	VersionID string       `json:"versionID"`
	Full      bool         `json:"full"`
}

// GetFullFriendUserIDs 相关类型
type GetFullFriendUserIDsReq struct {
	UserID string `json:"userID"`
}

type GetFullFriendUserIDsResp struct {
	BaseResp
	UserIDs []string `json:"userIDs"`
}

// GetSelfUnhandledApplyCount 相关类型
type GetSelfUnhandledApplyCountReq struct {
	UserID string `json:"userID"`
}

type GetSelfUnhandledApplyCountResp struct {
	BaseResp
	Count int64 `json:"count"`
}
