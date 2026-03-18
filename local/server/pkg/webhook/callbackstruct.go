package webhook

// 回调命令常量
const (
	CallbackBeforeUpdateUserInfoCommand    = "callbackBeforeUpdateUserInfoCommand"
	CallbackAfterUpdateUserInfoCommand     = "callbackAfterUpdateUserInfoCommand"
	CallbackBeforeUpdateUserInfoExCommand = "callbackBeforeUpdateUserInfoExCommand"
	CallbackAfterUpdateUserInfoExCommand  = "callbackAfterUpdateUserInfoExCommand"
	CallbackBeforeUserRegisterCommand     = "callbackBeforeUserRegisterCommand"
	CallbackAfterUserRegisterCommand       = "callbackAfterUserRegisterCommand"
	CallbackBeforeCreateGroupCommand      = "callbackBeforeCreateGroupCommand"
	CallbackAfterCreateGroupCommand        = "callbackAfterCreateGroupCommand"
	CallbackBeforeMembersJoinGroupCommand = "callbackBeforeMembersJoinGroupCommand"
	CallbackBeforeAddFriendCommand         = "callbackBeforeAddFriendCommand"
	CallbackAfterAddFriendCommand          = "callbackAfterAddFriendCommand"
	CallbackBeforeAddFriendAgreeCommand    = "callbackBeforeAddFriendAgreeCommand"
	CallbackAfterAddFriendAgreeCommand     = "callbackAfterAddFriendAgreeCommand"
	CallbackAfterDeleteFriendCommand       = "callbackAfterDeleteFriendCommand"
	CallbackBeforeSetFriendRemarkCommand   = "callbackBeforeSetFriendRemarkCommand"
	CallbackAfterSetFriendRemarkCommand    = "callbackAfterSetFriendRemarkCommand"
	CallbackBeforeAddBlackCommand          = "callbackBeforeAddBlackCommand"
	CallbackAfterAddBlackCommand           = "callbackAfterAddBlackCommand"
	CallbackAfterRemoveBlackCommand        = "callbackAfterRemoveBlackCommand"
	CallbackBeforeImportFriendsCommand     = "callbackBeforeImportFriendsCommand"
	CallbackAfterImportFriendsCommand      = "callbackAfterImportFriendsCommand"
	CallbackBeforeSetGroupInfoCommand      = "callbackBeforeSetGroupInfoCommand"
	CallbackAfterSetGroupInfoCommand       = "callbackAfterSetGroupInfoCommand"
	CallbackBeforeSetGroupMemberInfoCommand = "callbackBeforeSetGroupMemberInfoCommand"
	CallbackAfterSetGroupMemberInfoCommand  = "callbackAfterSetGroupMemberInfoCommand"
	CallbackAfterQuitGroupCommand          = "callbackAfterQuitGroupCommand"
	CallbackAfterKickGroupMemberCommand    = "callbackAfterKickGroupMemberCommand"
	CallbackAfterDismissGroupCommand       = "callbackAfterDismissGroupCommand"
	CallbackAfterTransferGroupOwnerCommand = "callbackAfterTransferGroupOwnerCommand"
	CallbackBeforeInviteJoinGroupCommand   = "callbackBeforeInviteJoinGroupCommand"
	CallbackAfterJoinGroupCommand          = "callbackAfterJoinGroupCommand"
)

// CallbackBeforeUpdateUserInfoReq 更新用户信息前回调请�?
type CallbackBeforeUpdateUserInfoReq struct {
	CallbackCommand string  `json:"callbackCommand"`
	UserID          string  `json:"userID"`
	FaceURL         *string `json:"faceURL,omitempty"`
	Nickname        *string `json:"nickname,omitempty"`
	Ex              *string `json:"ex,omitempty"`
}

func (c *CallbackBeforeUpdateUserInfoReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeUpdateUserInfoResp 更新用户信息前回调响�?
type CallbackBeforeUpdateUserInfoResp struct {
	CommonCallbackResp
	FaceURL  *string `json:"faceURL,omitempty"`
	Nickname *string `json:"nickname,omitempty"`
	Ex       *string `json:"ex,omitempty"`
}

// CallbackAfterUpdateUserInfoReq 更新用户信息后回调请�?
type CallbackAfterUpdateUserInfoReq struct {
	CallbackCommand string `json:"callbackCommand"`
	UserID          string `json:"userID"`
	FaceURL         string `json:"faceURL"`
	Nickname        string `json:"nickname"`
	Ex              string `json:"ex"`
}

func (c *CallbackAfterUpdateUserInfoReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterUpdateUserInfoResp 更新用户信息后回调响�?
type CallbackAfterUpdateUserInfoResp struct {
	CommonCallbackResp
}

// CallbackBeforeUpdateUserInfoExReq 更新用户信息Ex前回调请�?
type CallbackBeforeUpdateUserInfoExReq struct {
	CallbackCommand string  `json:"callbackCommand"`
	UserID          string  `json:"userID"`
	FaceURL         *string `json:"faceURL,omitempty"`
	Nickname        *string `json:"nickname,omitempty"`
	Ex              *string `json:"ex,omitempty"`
}

func (c *CallbackBeforeUpdateUserInfoExReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeUpdateUserInfoExResp 更新用户信息Ex前回调响�?
type CallbackBeforeUpdateUserInfoExResp struct {
	CommonCallbackResp
	FaceURL  *string `json:"faceURL,omitempty"`
	Nickname *string `json:"nickname,omitempty"`
	Ex       *string `json:"ex,omitempty"`
}

// CallbackAfterUpdateUserInfoExReq 更新用户信息Ex后回调请�?
type CallbackAfterUpdateUserInfoExReq struct {
	CallbackCommand string `json:"callbackCommand"`
	UserID          string `json:"userID"`
	FaceURL         string `json:"faceURL"`
	Nickname        string `json:"nickname"`
	Ex              string `json:"ex"`
}

func (c *CallbackAfterUpdateUserInfoExReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterUpdateUserInfoExResp 更新用户信息Ex后回调响�?
type CallbackAfterUpdateUserInfoExResp struct {
	CommonCallbackResp
}

// CallbackBeforeUserRegisterReq 用户注册前回调请�?
type CallbackBeforeUserRegisterReq struct {
	CallbackCommand string                   `json:"callbackCommand"`
	Users           []map[string]interface{} `json:"users"`
}

func (c *CallbackBeforeUserRegisterReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeUserRegisterResp 用户注册前回调响�?
type CallbackBeforeUserRegisterResp struct {
	CommonCallbackResp
	Users []map[string]interface{} `json:"users,omitempty"`
}

// CallbackAfterUserRegisterReq 用户注册后回调请�?
type CallbackAfterUserRegisterReq struct {
	CallbackCommand string                   `json:"callbackCommand"`
	Users           []map[string]interface{} `json:"users"`
}

func (c *CallbackAfterUserRegisterReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterUserRegisterResp 用户注册后回调响�?
type CallbackAfterUserRegisterResp struct {
	CommonCallbackResp
}

// CallbackBeforeCreateGroupReq 创建群组前回调请�?
type CallbackBeforeCreateGroupReq struct {
	CallbackCommand string                   `json:"callbackCommand"`
	GroupID         string                   `json:"groupID"`
	GroupName       string                   `json:"groupName"`
	GroupType       int32                    `json:"groupType"`
	OwnerUserID     string                   `json:"ownerUserID"`
	MemberUserIDs   []string                 `json:"memberUserIDs"`
	AdminUserIDs    []string                 `json:"adminUserIDs"`
	Ex              string                   `json:"ex"`
}

func (c *CallbackBeforeCreateGroupReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeCreateGroupResp 创建群组前回调响�?
type CallbackBeforeCreateGroupResp struct {
	CommonCallbackResp
	GroupName     *string   `json:"groupName,omitempty"`
	GroupType     *int32    `json:"groupType,omitempty"`
	OwnerUserID   *string   `json:"ownerUserID,omitempty"`
	MemberUserIDs []string  `json:"memberUserIDs,omitempty"`
	AdminUserIDs  []string  `json:"adminUserIDs,omitempty"`
	Ex            *string   `json:"ex,omitempty"`
}

// CallbackAfterCreateGroupReq 创建群组后回调请�?
type CallbackAfterCreateGroupReq struct {
	CallbackCommand string                   `json:"callbackCommand"`
	GroupID         string                   `json:"groupID"`
	GroupName       string                   `json:"groupName"`
	GroupType       int32                    `json:"groupType"`
	OwnerUserID     string                   `json:"ownerUserID"`
	MemberUserIDs   []string                 `json:"memberUserIDs"`
	AdminUserIDs    []string                 `json:"adminUserIDs"`
	Ex              string                   `json:"ex"`
}

func (c *CallbackAfterCreateGroupReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterCreateGroupResp 创建群组后回调响�?
type CallbackAfterCreateGroupResp struct {
	CommonCallbackResp
}

// CallbackBeforeMembersJoinGroupReq 成员加入群组前回调请�?
type CallbackBeforeMembersJoinGroupReq struct {
	CallbackCommand string                   `json:"callbackCommand"`
	GroupID         string                   `json:"groupID"`
	GroupType       int32                    `json:"groupType"`
	MemberUserIDs   []string                 `json:"memberUserIDs"`
	Ex              string                   `json:"ex"`
}

func (c *CallbackBeforeMembersJoinGroupReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeMembersJoinGroupResp 成员加入群组前回调响�?
type CallbackBeforeMembersJoinGroupResp struct {
	CommonCallbackResp
	MemberUserIDs []string `json:"memberUserIDs,omitempty"`
}

// CallbackBeforeAddFriendReq 添加好友前回调请�?
type CallbackBeforeAddFriendReq struct {
	CallbackCommand string `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID"`
	ToUserID        string `json:"toUserID"`
	ReqMsg          string `json:"reqMsg"`
	Ex              string `json:"ex"`
}

func (c *CallbackBeforeAddFriendReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeAddFriendResp 添加好友前回调响�?
type CallbackBeforeAddFriendResp struct {
	CommonCallbackResp
}

// CallbackAfterAddFriendReq 添加好友后回调请�?
type CallbackAfterAddFriendReq struct {
	CallbackCommand string `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID"`
	ToUserID        string `json:"toUserID"`
	ReqMsg          string `json:"reqMsg"`
}

func (c *CallbackAfterAddFriendReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterAddFriendResp 添加好友后回调响�?
type CallbackAfterAddFriendResp struct {
	CommonCallbackResp
}

// CallbackBeforeAddFriendAgreeReq 同意好友申请前回调请�?
type CallbackBeforeAddFriendAgreeReq struct {
	CallbackCommand string `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID"`
	ToUserID        string `json:"toUserID"`
	HandleMsg       string `json:"handleMsg"`
	HandleResult    int32  `json:"handleResult"`
}

func (c *CallbackBeforeAddFriendAgreeReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeAddFriendAgreeResp 同意好友申请前回调响�?
type CallbackBeforeAddFriendAgreeResp struct {
	CommonCallbackResp
}

// CallbackAfterAddFriendAgreeReq 同意好友申请后回调请�?
type CallbackAfterAddFriendAgreeReq struct {
	CallbackCommand string `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID"`
	ToUserID        string `json:"toUserID"`
	HandleMsg       string `json:"handleMsg"`
}

func (c *CallbackAfterAddFriendAgreeReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterAddFriendAgreeResp 同意好友申请后回调响�?
type CallbackAfterAddFriendAgreeResp struct {
	CommonCallbackResp
}

// CallbackAfterDeleteFriendReq 删除好友后回调请�?
type CallbackAfterDeleteFriendReq struct {
	CallbackCommand string `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID"`
	FriendUserID   string `json:"friendUserID"`
}

func (c *CallbackAfterDeleteFriendReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterDeleteFriendResp 删除好友后回调响�?
type CallbackAfterDeleteFriendResp struct {
	CommonCallbackResp
}

// CallbackBeforeSetFriendRemarkReq 设置好友备注前回调请�?
type CallbackBeforeSetFriendRemarkReq struct {
	CallbackCommand string `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID"`
	FriendUserID    string `json:"friendUserID"`
	Remark          string `json:"remark"`
}

func (c *CallbackBeforeSetFriendRemarkReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeSetFriendRemarkResp 设置好友备注前回调响�?
type CallbackBeforeSetFriendRemarkResp struct {
	CommonCallbackResp
	Remark *string `json:"remark,omitempty"`
}

// CallbackAfterSetFriendRemarkReq 设置好友备注后回调请�?
type CallbackAfterSetFriendRemarkReq struct {
	CallbackCommand string `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID"`
	FriendUserID    string `json:"friendUserID"`
	Remark          string `json:"remark"`
}

func (c *CallbackAfterSetFriendRemarkReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterSetFriendRemarkResp 设置好友备注后回调响�?
type CallbackAfterSetFriendRemarkResp struct {
	CommonCallbackResp
}

// CallbackBeforeAddBlackReq 添加黑名单前回调请求
type CallbackBeforeAddBlackReq struct {
	CallbackCommand string `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID"`
	BlackUserID     string `json:"blackUserID"`
}

func (c *CallbackBeforeAddBlackReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeAddBlackResp 添加黑名单前回调响应
type CallbackBeforeAddBlackResp struct {
	CommonCallbackResp
}

// CallbackAfterAddBlackReq 添加黑名单后回调请求
type CallbackAfterAddBlackReq struct {
	CallbackCommand string `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID"`
	BlackUserID     string `json:"blackUserID"`
}

func (c *CallbackAfterAddBlackReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterAddBlackResp 添加黑名单后回调响应
type CallbackAfterAddBlackResp struct {
	CommonCallbackResp
}

// CallbackAfterRemoveBlackReq 移除黑名单后回调请求
type CallbackAfterRemoveBlackReq struct {
	CallbackCommand string `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID"`
	BlackUserID     string `json:"blackUserID"`
}

func (c *CallbackAfterRemoveBlackReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterRemoveBlackResp 移除黑名单后回调响应
type CallbackAfterRemoveBlackResp struct {
	CommonCallbackResp
}

// CallbackBeforeImportFriendsReq 导入好友前回调请�?
type CallbackBeforeImportFriendsReq struct {
	CallbackCommand string   `json:"callbackCommand"`
	OwnerUserID     string   `json:"ownerUserID"`
	FriendUserIDs   []string `json:"friendUserIDs"`
}

func (c *CallbackBeforeImportFriendsReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeImportFriendsResp 导入好友前回调响�?
type CallbackBeforeImportFriendsResp struct {
	CommonCallbackResp
	FriendUserIDs []string `json:"friendUserIDs,omitempty"`
}

// CallbackAfterImportFriendsReq 导入好友后回调请�?
type CallbackAfterImportFriendsReq struct {
	CallbackCommand string   `json:"callbackCommand"`
	OwnerUserID     string   `json:"ownerUserID"`
	FriendUserIDs   []string `json:"friendUserIDs"`
}

func (c *CallbackAfterImportFriendsReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterImportFriendsResp 导入好友后回调响�?
type CallbackAfterImportFriendsResp struct {
	CommonCallbackResp
}

// CallbackBeforeSetGroupInfoReq 设置群组信息前回调请�?
type CallbackBeforeSetGroupInfoReq struct {
	CallbackCommand string  `json:"callbackCommand"`
	GroupID         string  `json:"groupID"`
	GroupName       *string `json:"groupName,omitempty"`
	Introduction    *string `json:"introduction,omitempty"`
	FaceURL         *string `json:"faceURL,omitempty"`
	Ex              *string `json:"ex,omitempty"`
}

func (c *CallbackBeforeSetGroupInfoReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeSetGroupInfoResp 设置群组信息前回调响�?
type CallbackBeforeSetGroupInfoResp struct {
	CommonCallbackResp
	GroupName    *string `json:"groupName,omitempty"`
	Introduction *string `json:"introduction,omitempty"`
	FaceURL      *string `json:"faceURL,omitempty"`
	Ex           *string `json:"ex,omitempty"`
}

// CallbackAfterSetGroupInfoReq 设置群组信息后回调请�?
type CallbackAfterSetGroupInfoReq struct {
	CallbackCommand string `json:"callbackCommand"`
	GroupID         string `json:"groupID"`
	GroupName       string `json:"groupName"`
	Introduction    string `json:"introduction"`
	FaceURL         string `json:"faceURL"`
	Ex              string `json:"ex"`
}

func (c *CallbackAfterSetGroupInfoReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterSetGroupInfoResp 设置群组信息后回调响�?
type CallbackAfterSetGroupInfoResp struct {
	CommonCallbackResp
}

// CallbackBeforeSetGroupMemberInfoReq 设置群成员信息前回调请求
type CallbackBeforeSetGroupMemberInfoReq struct {
	CallbackCommand string  `json:"callbackCommand"`
	GroupID         string  `json:"groupID"`
	UserID          string  `json:"userID"`
	Nickname        *string `json:"nickname,omitempty"`
	FaceURL         *string `json:"faceURL,omitempty"`
	Ex              *string `json:"ex,omitempty"`
}

func (c *CallbackBeforeSetGroupMemberInfoReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeSetGroupMemberInfoResp 设置群成员信息前回调响应
type CallbackBeforeSetGroupMemberInfoResp struct {
	CommonCallbackResp
	Nickname *string `json:"nickname,omitempty"`
	FaceURL  *string `json:"faceURL,omitempty"`
	Ex       *string `json:"ex,omitempty"`
}

// CallbackAfterSetGroupMemberInfoReq 设置群成员信息后回调请求
type CallbackAfterSetGroupMemberInfoReq struct {
	CallbackCommand string `json:"callbackCommand"`
	GroupID         string `json:"groupID"`
	UserID          string `json:"userID"`
	Nickname        string `json:"nickname"`
	FaceURL         string `json:"faceURL"`
	Ex              string `json:"ex"`
}

func (c *CallbackAfterSetGroupMemberInfoReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterSetGroupMemberInfoResp 设置群成员信息后回调响应
type CallbackAfterSetGroupMemberInfoResp struct {
	CommonCallbackResp
}

// CallbackAfterQuitGroupReq 退出群组后回调请求
type CallbackAfterQuitGroupReq struct {
	CallbackCommand string `json:"callbackCommand"`
	GroupID         string `json:"groupID"`
	UserID          string `json:"userID"`
}

func (c *CallbackAfterQuitGroupReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterQuitGroupResp 退出群组后回调响应
type CallbackAfterQuitGroupResp struct {
	CommonCallbackResp
}

// CallbackAfterKickGroupMemberReq 踢出群成员后回调请求
type CallbackAfterKickGroupMemberReq struct {
	CallbackCommand string   `json:"callbackCommand"`
	GroupID         string   `json:"groupID"`
	KickedUserIDs   []string `json:"kickedUserIDs"`
	OpUserID        string   `json:"opUserID"`
}

func (c *CallbackAfterKickGroupMemberReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterKickGroupMemberResp 踢出群成员后回调响应
type CallbackAfterKickGroupMemberResp struct {
	CommonCallbackResp
}

// CallbackAfterDismissGroupReq 解散群组后回调请�?
type CallbackAfterDismissGroupReq struct {
	CallbackCommand string `json:"callbackCommand"`
	GroupID         string `json:"groupID"`
	OpUserID        string `json:"opUserID"`
}

func (c *CallbackAfterDismissGroupReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterDismissGroupResp 解散群组后回调响�?
type CallbackAfterDismissGroupResp struct {
	CommonCallbackResp
}

// CallbackAfterTransferGroupOwnerReq 转让群主后回调请�?
type CallbackAfterTransferGroupOwnerReq struct {
	CallbackCommand string `json:"callbackCommand"`
	GroupID         string `json:"groupID"`
	OldOwnerUserID  string `json:"oldOwnerUserID"`
	NewOwnerUserID  string `json:"newOwnerUserID"`
}

func (c *CallbackAfterTransferGroupOwnerReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterTransferGroupOwnerResp 转让群主后回调响�?
type CallbackAfterTransferGroupOwnerResp struct {
	CommonCallbackResp
}

// CallbackBeforeInviteJoinGroupReq 邀请加入群组前回调请求
type CallbackBeforeInviteJoinGroupReq struct {
	CallbackCommand string   `json:"callbackCommand"`
	GroupID         string   `json:"groupID"`
	InvitedUserIDs  []string `json:"invitedUserIDs"`
	Reason          string   `json:"reason"`
}

func (c *CallbackBeforeInviteJoinGroupReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeInviteJoinGroupResp 邀请加入群组前回调响应
type CallbackBeforeInviteJoinGroupResp struct {
	CommonCallbackResp
	RefusedMembersAccount []string `json:"refusedMembersAccount,omitempty"`
}

// CallbackAfterJoinGroupReq 加入群组后回调请�?
type CallbackAfterJoinGroupReq struct {
	CallbackCommand string `json:"callbackCommand"`
	GroupID         string `json:"groupID"`
	UserID          string `json:"userID"`
	ReqMessage      string `json:"reqMessage"`
	JoinSource      int32  `json:"joinSource"`
	InviterUserID   string `json:"inviterUserID"`
}

func (c *CallbackAfterJoinGroupReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterJoinGroupResp 加入群组后回调响�?
type CallbackAfterJoinGroupResp struct {
	CommonCallbackResp
}
