package types

// Group 相关请求和响应类型
// 注意：这些类型需要与 open-im-server 的 API 接口保持一致
// 复杂的嵌套类型（如 GroupInfo）直接使用 protocol 包中的类型

// CreateGroup 相关类型
type CreateGroupReq struct {
	MemberUserIDs []string    `json:"memberUserIDs"`
	GroupInfo     interface{} `json:"groupInfo"` // 使用 sdkws.GroupInfo
	AdminUserIDs  []string    `json:"adminUserIDs,optional"`
	OwnerUserID   string      `json:"ownerUserID"`
	SendMessage   *bool       `json:"sendMessage,optional"`
}

type CreateGroupResp struct {
	BaseResp
	GroupInfo interface{} `json:"groupInfo"` // 使用 sdkws.GroupInfo
}

// SetGroupInfo 相关类型
type SetGroupInfoReq struct {
	GroupInfoForSet interface{} `json:"groupInfoForSet"` // 使用 sdkws.GroupInfoForSet
}

type SetGroupInfoResp struct {
	BaseResp
}

// SetGroupInfoEx 相关类型
type SetGroupInfoExReq struct {
	GroupID         string  `json:"groupID"`
	GroupName       *string `json:"groupName,optional"`
	Notification    *string `json:"notification,optional"`
	Introduction    *string `json:"introduction,optional"`
	FaceURL         *string `json:"faceURL,optional"`
	Ex              *string `json:"ex,optional"`
	NeedVerification *int32  `json:"needVerification,optional"`
	LookMemberInfo  *int32  `json:"lookMemberInfo,optional"`
	ApplyMemberFriend *int32  `json:"applyMemberFriend,optional"`
}

type SetGroupInfoExResp struct {
	BaseResp
}

// JoinGroup 相关类型
type JoinGroupReq struct {
	GroupID      string `json:"groupID"`
	ReqMessage   string `json:"reqMessage,optional"`
	JoinSource   int32  `json:"joinSource"`
	InviterUserID string `json:"inviterUserID,optional"`
	Ex           string `json:"ex,optional"`
}

type JoinGroupResp struct {
	BaseResp
}

// QuitGroup 相关类型
type QuitGroupReq struct {
	GroupID string `json:"groupID"`
	UserID  string `json:"userID"`
}

type QuitGroupResp struct {
	BaseResp
}

// ApplicationGroupResponse 相关类型
type ApplicationGroupResponseReq struct {
	GroupID     string `json:"groupID"`
	FromUserID  string `json:"fromUserID"`
	HandledMsg  string `json:"handledMsg,optional"`
	HandleResult int32  `json:"handleResult"`
}

type ApplicationGroupResponseResp struct {
	BaseResp
}

// TransferGroupOwner 相关类型
type TransferGroupOwnerReq struct {
	GroupID       string `json:"groupID"`
	OldOwnerUserID string `json:"oldOwnerUserID"`
	NewOwnerUserID string `json:"newOwnerUserID"`
}

type TransferGroupOwnerResp struct {
	BaseResp
}

// GetRecvGroupApplicationList 相关类型
type GetRecvGroupApplicationListReq struct {
	Pagination   Pagination `json:"pagination"`
	FromUserID   string     `json:"fromUserID"`
	GroupIDs     []string   `json:"groupIDs,optional"`
	HandleResults []int32    `json:"handleResults,optional"`
}

type GetRecvGroupApplicationListResp struct {
	BaseResp
	Total         uint32        `json:"total"`
	GroupRequests []interface{} `json:"groupRequests"` // 使用 sdkws.GroupRequest
}

// GetUserReqGroupApplicationList 相关类型
type GetUserReqGroupApplicationListReq struct {
	Pagination   Pagination `json:"pagination"`
	UserID       string     `json:"userID"`
	GroupIDs     []string   `json:"groupIDs,optional"`
	HandleResults []int32    `json:"handleResults,optional"`
}

type GetUserReqGroupApplicationListResp struct {
	BaseResp
	Total         uint32        `json:"total"`
	GroupRequests []interface{} `json:"groupRequests"` // 使用 sdkws.GroupRequest
}

// GetGroupUsersReqApplicationList 相关类型
type GetGroupUsersReqApplicationListReq struct {
	GroupID      string     `json:"groupID"`
	UserIDs      []string   `json:"userIDs,optional"`
}

type GetGroupUsersReqApplicationListResp struct {
	BaseResp
	Total         uint32        `json:"total"`
	GroupRequests []interface{} `json:"groupRequests"` // 使用 sdkws.GroupRequest
}

// GetSpecifiedUserGroupRequestInfo 相关类型
type GetSpecifiedUserGroupRequestInfoReq struct {
	GroupID string `json:"groupID"`
	UserID  string `json:"userID"`
}

type GetSpecifiedUserGroupRequestInfoResp struct {
	BaseResp
	Total         uint32        `json:"total"`
	GroupRequests []interface{} `json:"groupRequests"` // 使用 sdkws.GroupRequest
}

// GetGroupsInfo 相关类型
type GetGroupsInfoReq struct {
	GroupIDs []string `json:"groupIDs"`
}

type GetGroupsInfoResp struct {
	BaseResp
	GroupInfos []interface{} `json:"groupInfos"` // 使用 sdkws.GroupInfo
}

// KickGroupMember 相关类型
type KickGroupMemberReq struct {
	GroupID      string   `json:"groupID"`
	KickedUserIDs []string `json:"kickedUserIDs"`
	Reason       string   `json:"reason,optional"`
	SendMessage  *bool    `json:"sendMessage,optional"`
}

type KickGroupMemberResp struct {
	BaseResp
}

// GetGroupMembersInfo 相关类型
type GetGroupMembersInfoReq struct {
	GroupID string   `json:"groupID"`
	UserIDs []string `json:"userIDs"`
}

type GetGroupMembersInfoResp struct {
	BaseResp
	Members []interface{} `json:"members"` // 使用 sdkws.GroupMemberFullInfo
}

// GetGroupMemberList 相关类型
type GetGroupMemberListReq struct {
	Pagination Pagination `json:"pagination"`
	GroupID    string     `json:"groupID"`
	Filter     int32      `json:"filter,optional"`
	Keyword    string     `json:"keyword,optional"`
}

type GetGroupMemberListResp struct {
	BaseResp
	Total   uint32        `json:"total"`
	Members []interface{} `json:"members"` // 使用 sdkws.GroupMemberFullInfo
}

// InviteUserToGroup 相关类型
type InviteUserToGroupReq struct {
	GroupID       string   `json:"groupID"`
	Reason        string   `json:"reason,optional"`
	InvitedUserIDs []string `json:"invitedUserIDs"`
	SendMessage   *bool    `json:"sendMessage,optional"`
}

type InviteUserToGroupResp struct {
	BaseResp
}

// GetJoinedGroupList 相关类型
type GetJoinedGroupListReq struct {
	Pagination Pagination `json:"pagination"`
	FromUserID string     `json:"fromUserID"`
}

type GetJoinedGroupListResp struct {
	BaseResp
	Total  uint32        `json:"total"`
	Groups []interface{} `json:"groups"` // 使用 sdkws.GroupInfo
}

// DismissGroup 相关类型
type DismissGroupReq struct {
	GroupID string `json:"groupID"`
}

type DismissGroupResp struct {
	BaseResp
}

// MuteGroupMember 相关类型
type MuteGroupMemberReq struct {
	GroupID    string `json:"groupID"`
	UserID     string `json:"userID"`
	MutedSeconds uint32 `json:"mutedSeconds"`
}

type MuteGroupMemberResp struct {
	BaseResp
}

// CancelMuteGroupMember 相关类型
type CancelMuteGroupMemberReq struct {
	GroupID string `json:"groupID"`
	UserID  string `json:"userID"`
}

type CancelMuteGroupMemberResp struct {
	BaseResp
}

// MuteGroup 相关类型
type MuteGroupReq struct {
	GroupID string `json:"groupID"`
}

type MuteGroupResp struct {
	BaseResp
}

// CancelMuteGroup 相关类型
type CancelMuteGroupReq struct {
	GroupID string `json:"groupID"`
}

type CancelMuteGroupResp struct {
	BaseResp
}

// SetGroupMemberInfo 相关类型
type SetGroupMemberInfoReq struct {
	GroupID string      `json:"groupID"`
	UserID  string      `json:"userID"`
	Info    interface{} `json:"info,optional"` // 使用 sdkws.GroupMemberInfoForSet
}

type SetGroupMemberInfoResp struct {
	BaseResp
}

// GetGroupAbstractInfo 相关类型
type GetGroupAbstractInfoReq struct {
	GroupIDs []string `json:"groupIDs"`
}

type GetGroupAbstractInfoResp struct {
	BaseResp
	GroupAbstractInfos []interface{} `json:"groupAbstractInfos"` // 使用 user.GroupAbstractInfo
}

// GetGroups 相关类型
type GetGroupsReq struct {
	Pagination Pagination `json:"pagination"`
	GroupName  string     `json:"groupName,optional"`
}

type GetGroupsResp struct {
	BaseResp
	Total  uint32        `json:"total"`
	Groups []interface{} `json:"groups"` // 使用 sdkws.GroupInfo
}

// GetGroupMemberUserIDs 相关类型
type GetGroupMemberUserIDsReq struct {
	GroupID string `json:"groupID"`
}

type GetGroupMemberUserIDsResp struct {
	BaseResp
	UserIDs []string `json:"userIDs"`
}

// GetIncrementalJoinGroup 相关类型
type GetIncrementalJoinGroupReq struct {
	UserID    string `json:"userID"`
	VersionID string `json:"versionID,optional"`
	Version   uint64 `json:"version"`
}

type GetIncrementalJoinGroupResp struct {
	BaseResp
	Version     uint64        `json:"version"`
	VersionID   string        `json:"versionID"`
	Full        bool          `json:"full"`
	Delete      []string      `json:"delete"`
	Insert      []interface{} `json:"insert"` // 使用 sdkws.GroupInfo
	Update      []interface{} `json:"update"` // 使用 sdkws.GroupInfo
	Groups      []interface{} `json:"groups"` // Insert + Update 的组合，用于兼容
}

// GetIncrementalGroupMember 相关类型
type GetIncrementalGroupMemberReq struct {
	GroupID    string `json:"groupID"`
	VersionID string `json:"versionID,optional"`
	Version    uint64 `json:"version"`
}

type GetIncrementalGroupMemberResp struct {
	BaseResp
	Version     uint64        `json:"version"`
	VersionID   string        `json:"versionID"`
	Full        bool          `json:"full"`
	Delete      []string      `json:"delete"`
	Insert      []interface{} `json:"insert"` // 使用 sdkws.GroupMemberFullInfo
	Update      []interface{} `json:"update"` // 使用 sdkws.GroupMemberFullInfo
}

// GetIncrementalGroupMemberBatch 相关类型
type GetIncrementalGroupMemberBatchReq struct {
	ReqList []GetIncrementalGroupMemberReq `json:"reqList"`
}

type GetIncrementalGroupMemberBatchResp struct {
	BaseResp
	RespList map[string]GetIncrementalGroupMemberResp `json:"respList"`
}

// GetFullGroupMemberUserIDs 相关类型
type GetFullGroupMemberUserIDsReq struct {
	GroupID string `json:"groupID"`
}

type GetFullGroupMemberUserIDsResp struct {
	BaseResp
	UserIDs []string `json:"userIDs"`
}

// GetFullJoinGroupIDs 相关类型
type GetFullJoinGroupIDsReq struct {
	UserID string `json:"userID"`
}

type GetFullJoinGroupIDsResp struct {
	BaseResp
	GroupIDs []string `json:"groupIDs"`
}

// GetGroupApplicationUnhandledCount 相关类型
type GetGroupApplicationUnhandledCountReq struct {
	UserID string `json:"userID"`
	Time   int64  `json:"time"`
}

type GetGroupApplicationUnhandledCountResp struct {
	BaseResp
	Count int64 `json:"count"`
}
