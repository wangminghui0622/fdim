package types

// CreateGroup
type CreateGroupReq struct {
	MemberUserIDs []string    `json:"memberUserIDs"`
	GroupInfo     interface{} `json:"groupInfo"` // ʹ sdkws.GroupInfo
	AdminUserIDs  []string    `json:"adminUserIDs,optional"`
	OwnerUserID   string      `json:"ownerUserID"`
	SendMessage   *bool       `json:"sendMessage,optional"`
}

type CreateGroupResp struct {
	BaseResp
	GroupInfo interface{} `json:"groupInfo"` // ʹ sdkws.GroupInfo
}

// SetGroupInfo
type SetGroupInfoReq struct {
	GroupInfoForSet interface{} `json:"groupInfoForSet"` // ʹ sdkws.GroupInfoForSet
}

type SetGroupInfoResp struct {
	BaseResp
}

// SetGroupInfoEx
type SetGroupInfoExReq struct {
	GroupID           string  `json:"groupID"`
	GroupName         *string `json:"groupName,optional"`
	Notification      *string `json:"notification,optional"`
	Introduction      *string `json:"introduction,optional"`
	FaceURL           *string `json:"faceURL,optional"`
	Ex                *string `json:"ex,optional"`
	NeedVerification  *int32  `json:"needVerification,optional"`
	LookMemberInfo    *int32  `json:"lookMemberInfo,optional"`
	ApplyMemberFriend *int32  `json:"applyMemberFriend,optional"`
}

type SetGroupInfoExResp struct {
	BaseResp
}

// JoinGroup
type JoinGroupReq struct {
	GroupID       string `json:"groupID"`
	ReqMessage    string `json:"reqMessage,optional"`
	JoinSource    int32  `json:"joinSource"`
	InviterUserID string `json:"inviterUserID,optional"`
	Ex            string `json:"ex,optional"`
}

type JoinGroupResp struct {
	BaseResp
}

// QuitGroup
type QuitGroupReq struct {
	GroupID string `json:"groupID"`
	UserID  string `json:"userID"`
}

type QuitGroupResp struct {
	BaseResp
}

// ApplicationGroupResponse
type ApplicationGroupResponseReq struct {
	GroupID      string `json:"groupID"`
	FromUserID   string `json:"fromUserID"`
	HandledMsg   string `json:"handledMsg,optional"`
	HandleResult int32  `json:"handleResult"`
}

type ApplicationGroupResponseResp struct {
	BaseResp
}

// TransferGroupOwner
type TransferGroupOwnerReq struct {
	GroupID        string `json:"groupID"`
	OldOwnerUserID string `json:"oldOwnerUserID"`
	NewOwnerUserID string `json:"newOwnerUserID"`
}

type TransferGroupOwnerResp struct {
	BaseResp
}

// GetRecvGroupApplicationList
type GetRecvGroupApplicationListReq struct {
	Pagination    Pagination `json:"pagination"`
	FromUserID    string     `json:"fromUserID"`
	GroupIDs      []string   `json:"groupIDs,optional"`
	HandleResults []int32    `json:"handleResults,optional"`
}

type GetRecvGroupApplicationListResp struct {
	BaseResp
	Total         uint32        `json:"total"`
	GroupRequests []interface{} `json:"groupRequests"` // ʹ sdkws.GroupRequest
}

// GetUserReqGroupApplicationList
type GetUserReqGroupApplicationListReq struct {
	Pagination    Pagination `json:"pagination"`
	UserID        string     `json:"userID"`
	GroupIDs      []string   `json:"groupIDs,optional"`
	HandleResults []int32    `json:"handleResults,optional"`
}

type GetUserReqGroupApplicationListResp struct {
	BaseResp
	Total         uint32        `json:"total"`
	GroupRequests []interface{} `json:"groupRequests"` // ʹ sdkws.GroupRequest
}

// GetGroupUsersReqApplicationList
type GetGroupUsersReqApplicationListReq struct {
	GroupID string   `json:"groupID"`
	UserIDs []string `json:"userIDs,optional"`
}

type GetGroupUsersReqApplicationListResp struct {
	BaseResp
	Total         uint32        `json:"total"`
	GroupRequests []interface{} `json:"groupRequests"` // ʹ sdkws.GroupRequest
}

// GetSpecifiedUserGroupRequestInfo
type GetSpecifiedUserGroupRequestInfoReq struct {
	GroupID string `json:"groupID"`
	UserID  string `json:"userID"`
}

type GetSpecifiedUserGroupRequestInfoResp struct {
	BaseResp
	Total         uint32        `json:"total"`
	GroupRequests []interface{} `json:"groupRequests"` // ʹ sdkws.GroupRequest
}

// GetGroupsInfo
type GetGroupsInfoReq struct {
	GroupIDs []string `json:"groupIDs"`
}

type GetGroupsInfoResp struct {
	BaseResp
	GroupInfos []interface{} `json:"groupInfos"` // ʹ sdkws.GroupInfo
}

// KickGroupMember
type KickGroupMemberReq struct {
	GroupID       string   `json:"groupID"`
	KickedUserIDs []string `json:"kickedUserIDs"`
	Reason        string   `json:"reason,optional"`
	SendMessage   *bool    `json:"sendMessage,optional"`
}

type KickGroupMemberResp struct {
	BaseResp
}

// GetGroupMembersInfo
type GetGroupMembersInfoReq struct {
	GroupID string   `json:"groupID"`
	UserIDs []string `json:"userIDs"`
}

type GetGroupMembersInfoResp struct {
	BaseResp
	Members []interface{} `json:"members"` // ʹ sdkws.GroupMemberFullInfo
}

// GetGroupMemberList
type GetGroupMemberListReq struct {
	Pagination Pagination `json:"pagination"`
	GroupID    string     `json:"groupID"`
	Filter     int32      `json:"filter,optional"`
	Keyword    string     `json:"keyword,optional"`
}

type GetGroupMemberListResp struct {
	BaseResp
	Total   uint32        `json:"total"`
	Members []interface{} `json:"members"` // ʹ sdkws.GroupMemberFullInfo
}

// InviteUserToGroup
type InviteUserToGroupReq struct {
	GroupID        string   `json:"groupID"`
	Reason         string   `json:"reason,optional"`
	InvitedUserIDs []string `json:"invitedUserIDs"`
	SendMessage    *bool    `json:"sendMessage,optional"`
}

type InviteUserToGroupResp struct {
	BaseResp
}

// GetJoinedGroupList
type GetJoinedGroupListReq struct {
	Pagination Pagination `json:"pagination"`
	FromUserID string     `json:"fromUserID"`
}

type GetJoinedGroupListResp struct {
	BaseResp
	Total  uint32        `json:"total"`
	Groups []interface{} `json:"groups"` // ʹ sdkws.GroupInfo
}

// DismissGroup
type DismissGroupReq struct {
	GroupID string `json:"groupID"`
}

type DismissGroupResp struct {
	BaseResp
}

// MuteGroupMember
type MuteGroupMemberReq struct {
	GroupID      string `json:"groupID"`
	UserID       string `json:"userID"`
	MutedSeconds uint32 `json:"mutedSeconds"`
}

type MuteGroupMemberResp struct {
	BaseResp
}

// CancelMuteGroupMember
type CancelMuteGroupMemberReq struct {
	GroupID string `json:"groupID"`
	UserID  string `json:"userID"`
}

type CancelMuteGroupMemberResp struct {
	BaseResp
}

// MuteGroup
type MuteGroupReq struct {
	GroupID string `json:"groupID"`
}

type MuteGroupResp struct {
	BaseResp
}

// CancelMuteGroup
type CancelMuteGroupReq struct {
	GroupID string `json:"groupID"`
}

type CancelMuteGroupResp struct {
	BaseResp
}

// SetGroupMemberInfo
type SetGroupMemberInfoReq struct {
	GroupID string      `json:"groupID"`
	UserID  string      `json:"userID"`
	Info    interface{} `json:"info,optional"` // ʹ sdkws.GroupMemberInfoForSet
}

type SetGroupMemberInfoResp struct {
	BaseResp
}

// GetGroupAbstractInfo
type GetGroupAbstractInfoReq struct {
	GroupIDs []string `json:"groupIDs"`
}

type GetGroupAbstractInfoResp struct {
	BaseResp
	GroupAbstractInfos []interface{} `json:"groupAbstractInfos"` // ʹ user.GroupAbstractInfo
}

// GetGroups
type GetGroupsReq struct {
	Pagination Pagination `json:"pagination"`
	GroupName  string     `json:"groupName,optional"`
}

type GetGroupsResp struct {
	BaseResp
	Total  uint32        `json:"total"`
	Groups []interface{} `json:"groups"` // ʹ sdkws.GroupInfo
}

// GetGroupMemberUserIDs
type GetGroupMemberUserIDsReq struct {
	GroupID string `json:"groupID"`
}

type GetGroupMemberUserIDsResp struct {
	BaseResp
	UserIDs []string `json:"userIDs"`
}

// GetIncrementalJoinGroup
type GetIncrementalJoinGroupReq struct {
	UserID    string `json:"userID"`
	VersionID string `json:"versionID,optional"`
	Version   uint64 `json:"version"`
}

type GetIncrementalJoinGroupResp struct {
	BaseResp
	Version   uint64        `json:"version"`
	VersionID string        `json:"versionID"`
	Full      bool          `json:"full"`
	Delete    []string      `json:"delete"`
	Insert    []interface{} `json:"insert"` // ʹ sdkws.GroupInfo
	Update    []interface{} `json:"update"` // ʹ sdkws.GroupInfo
	Groups    []interface{} `json:"groups"` // Insert + Update ϣڼ
}

// GetIncrementalGroupMember
type GetIncrementalGroupMemberReq struct {
	GroupID   string `json:"groupID"`
	VersionID string `json:"versionID,optional"`
	Version   uint64 `json:"version"`
}

type GetIncrementalGroupMemberResp struct {
	BaseResp
	Version   uint64        `json:"version"`
	VersionID string        `json:"versionID"`
	Full      bool          `json:"full"`
	Delete    []string      `json:"delete"`
	Insert    []interface{} `json:"insert"` // ʹ sdkws.GroupMemberFullInfo
	Update    []interface{} `json:"update"` // ʹ sdkws.GroupMemberFullInfo
}

// GetIncrementalGroupMemberBatch
type GetIncrementalGroupMemberBatchReq struct {
	ReqList []GetIncrementalGroupMemberReq `json:"reqList"`
}

type GetIncrementalGroupMemberBatchResp struct {
	BaseResp
	RespList map[string]GetIncrementalGroupMemberResp `json:"respList"`
}

// GetFullGroupMemberUserIDs
type GetFullGroupMemberUserIDsReq struct {
	GroupID string `json:"groupID"`
}

type GetFullGroupMemberUserIDsResp struct {
	BaseResp
	UserIDs []string `json:"userIDs"`
}

// GetFullJoinGroupIDs
type GetFullJoinGroupIDsReq struct {
	UserID string `json:"userID"`
}

type GetFullJoinGroupIDsResp struct {
	BaseResp
	GroupIDs []string `json:"groupIDs"`
}

// GetGroupApplicationUnhandledCount
type GetGroupApplicationUnhandledCountReq struct {
	UserID string `json:"userID"`
	Time   int64  `json:"time"`
}

type GetGroupApplicationUnhandledCountResp struct {
	BaseResp
	Count int64 `json:"count"`
}
