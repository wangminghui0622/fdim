package webhook

// ص
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

// CallbackBeforeUpdateUserInfoReq ûϢǰص??
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

// CallbackBeforeUpdateUserInfoResp ûϢǰص??
type CallbackBeforeUpdateUserInfoResp struct {
	CommonCallbackResp
	FaceURL  *string `json:"faceURL,omitempty"`
	Nickname *string `json:"nickname,omitempty"`
	Ex       *string `json:"ex,omitempty"`
}

// CallbackAfterUpdateUserInfoReq ûϢص??
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

// CallbackAfterUpdateUserInfoResp ûϢص??
type CallbackAfterUpdateUserInfoResp struct {
	CommonCallbackResp
}

// CallbackBeforeUpdateUserInfoExReq ûϢExǰص??
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

// CallbackBeforeUpdateUserInfoExResp ûϢExǰص??
type CallbackBeforeUpdateUserInfoExResp struct {
	CommonCallbackResp
	FaceURL  *string `json:"faceURL,omitempty"`
	Nickname *string `json:"nickname,omitempty"`
	Ex       *string `json:"ex,omitempty"`
}

// CallbackAfterUpdateUserInfoExReq ûϢExص??
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

// CallbackAfterUpdateUserInfoExResp ûϢExص??
type CallbackAfterUpdateUserInfoExResp struct {
	CommonCallbackResp
}

// CallbackBeforeUserRegisterReq ûעǰص??
type CallbackBeforeUserRegisterReq struct {
	CallbackCommand string                   `json:"callbackCommand"`
	Users           []map[string]interface{} `json:"users"`
}

func (c *CallbackBeforeUserRegisterReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeUserRegisterResp ûעǰص??
type CallbackBeforeUserRegisterResp struct {
	CommonCallbackResp
	Users []map[string]interface{} `json:"users,omitempty"`
}

// CallbackAfterUserRegisterReq ûעص??
type CallbackAfterUserRegisterReq struct {
	CallbackCommand string                   `json:"callbackCommand"`
	Users           []map[string]interface{} `json:"users"`
}

func (c *CallbackAfterUserRegisterReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterUserRegisterResp ûעص??
type CallbackAfterUserRegisterResp struct {
	CommonCallbackResp
}

// CallbackBeforeCreateGroupReq Ⱥǰص??
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

// CallbackBeforeCreateGroupResp Ⱥǰص??
type CallbackBeforeCreateGroupResp struct {
	CommonCallbackResp
	GroupName     *string   `json:"groupName,omitempty"`
	GroupType     *int32    `json:"groupType,omitempty"`
	OwnerUserID   *string   `json:"ownerUserID,omitempty"`
	MemberUserIDs []string  `json:"memberUserIDs,omitempty"`
	AdminUserIDs  []string  `json:"adminUserIDs,omitempty"`
	Ex            *string   `json:"ex,omitempty"`
}

// CallbackAfterCreateGroupReq Ⱥص??
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

// CallbackAfterCreateGroupResp Ⱥص??
type CallbackAfterCreateGroupResp struct {
	CommonCallbackResp
}

// CallbackBeforeMembersJoinGroupReq ԱȺǰص??
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

// CallbackBeforeMembersJoinGroupResp ԱȺǰص??
type CallbackBeforeMembersJoinGroupResp struct {
	CommonCallbackResp
	MemberUserIDs []string `json:"memberUserIDs,omitempty"`
}

// CallbackBeforeAddFriendReq Ӻǰص??
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

// CallbackBeforeAddFriendResp Ӻǰص??
type CallbackBeforeAddFriendResp struct {
	CommonCallbackResp
}

// CallbackAfterAddFriendReq ӺѺص??
type CallbackAfterAddFriendReq struct {
	CallbackCommand string `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID"`
	ToUserID        string `json:"toUserID"`
	ReqMsg          string `json:"reqMsg"`
}

func (c *CallbackAfterAddFriendReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterAddFriendResp ӺѺص??
type CallbackAfterAddFriendResp struct {
	CommonCallbackResp
}

// CallbackBeforeAddFriendAgreeReq ͬǰص??
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

// CallbackBeforeAddFriendAgreeResp ͬǰص??
type CallbackBeforeAddFriendAgreeResp struct {
	CommonCallbackResp
}

// CallbackAfterAddFriendAgreeReq ͬص??
type CallbackAfterAddFriendAgreeReq struct {
	CallbackCommand string `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID"`
	ToUserID        string `json:"toUserID"`
	HandleMsg       string `json:"handleMsg"`
}

func (c *CallbackAfterAddFriendAgreeReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterAddFriendAgreeResp ͬص??
type CallbackAfterAddFriendAgreeResp struct {
	CommonCallbackResp
}

// CallbackAfterDeleteFriendReq ɾѺص??
type CallbackAfterDeleteFriendReq struct {
	CallbackCommand string `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID"`
	FriendUserID   string `json:"friendUserID"`
}

func (c *CallbackAfterDeleteFriendReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterDeleteFriendResp ɾѺص??
type CallbackAfterDeleteFriendResp struct {
	CommonCallbackResp
}

// CallbackBeforeSetFriendRemarkReq úѱעǰص??
type CallbackBeforeSetFriendRemarkReq struct {
	CallbackCommand string `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID"`
	FriendUserID    string `json:"friendUserID"`
	Remark          string `json:"remark"`
}

func (c *CallbackBeforeSetFriendRemarkReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeSetFriendRemarkResp úѱעǰص??
type CallbackBeforeSetFriendRemarkResp struct {
	CommonCallbackResp
	Remark *string `json:"remark,omitempty"`
}

// CallbackAfterSetFriendRemarkReq úѱעص??
type CallbackAfterSetFriendRemarkReq struct {
	CallbackCommand string `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID"`
	FriendUserID    string `json:"friendUserID"`
	Remark          string `json:"remark"`
}

func (c *CallbackAfterSetFriendRemarkReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterSetFriendRemarkResp úѱעص??
type CallbackAfterSetFriendRemarkResp struct {
	CommonCallbackResp
}

// CallbackBeforeAddBlackReq Ӻǰص
type CallbackBeforeAddBlackReq struct {
	CallbackCommand string `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID"`
	BlackUserID     string `json:"blackUserID"`
}

func (c *CallbackBeforeAddBlackReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeAddBlackResp ӺǰصӦ
type CallbackBeforeAddBlackResp struct {
	CommonCallbackResp
}

// CallbackAfterAddBlackReq Ӻص
type CallbackAfterAddBlackReq struct {
	CallbackCommand string `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID"`
	BlackUserID     string `json:"blackUserID"`
}

func (c *CallbackAfterAddBlackReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterAddBlackResp ӺصӦ
type CallbackAfterAddBlackResp struct {
	CommonCallbackResp
}

// CallbackAfterRemoveBlackReq Ƴص
type CallbackAfterRemoveBlackReq struct {
	CallbackCommand string `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID"`
	BlackUserID     string `json:"blackUserID"`
}

func (c *CallbackAfterRemoveBlackReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterRemoveBlackResp ƳصӦ
type CallbackAfterRemoveBlackResp struct {
	CommonCallbackResp
}

// CallbackBeforeImportFriendsReq ǰص??
type CallbackBeforeImportFriendsReq struct {
	CallbackCommand string   `json:"callbackCommand"`
	OwnerUserID     string   `json:"ownerUserID"`
	FriendUserIDs   []string `json:"friendUserIDs"`
}

func (c *CallbackBeforeImportFriendsReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeImportFriendsResp ǰص??
type CallbackBeforeImportFriendsResp struct {
	CommonCallbackResp
	FriendUserIDs []string `json:"friendUserIDs,omitempty"`
}

// CallbackAfterImportFriendsReq Ѻص??
type CallbackAfterImportFriendsReq struct {
	CallbackCommand string   `json:"callbackCommand"`
	OwnerUserID     string   `json:"ownerUserID"`
	FriendUserIDs   []string `json:"friendUserIDs"`
}

func (c *CallbackAfterImportFriendsReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterImportFriendsResp Ѻص??
type CallbackAfterImportFriendsResp struct {
	CommonCallbackResp
}

// CallbackBeforeSetGroupInfoReq ȺϢǰص??
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

// CallbackBeforeSetGroupInfoResp ȺϢǰص??
type CallbackBeforeSetGroupInfoResp struct {
	CommonCallbackResp
	GroupName    *string `json:"groupName,omitempty"`
	Introduction *string `json:"introduction,omitempty"`
	FaceURL      *string `json:"faceURL,omitempty"`
	Ex           *string `json:"ex,omitempty"`
}

// CallbackAfterSetGroupInfoReq ȺϢص??
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

// CallbackAfterSetGroupInfoResp ȺϢص??
type CallbackAfterSetGroupInfoResp struct {
	CommonCallbackResp
}

// CallbackBeforeSetGroupMemberInfoReq ȺԱϢǰص
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

// CallbackBeforeSetGroupMemberInfoResp ȺԱϢǰصӦ
type CallbackBeforeSetGroupMemberInfoResp struct {
	CommonCallbackResp
	Nickname *string `json:"nickname,omitempty"`
	FaceURL  *string `json:"faceURL,omitempty"`
	Ex       *string `json:"ex,omitempty"`
}

// CallbackAfterSetGroupMemberInfoReq ȺԱϢص
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

// CallbackAfterSetGroupMemberInfoResp ȺԱϢصӦ
type CallbackAfterSetGroupMemberInfoResp struct {
	CommonCallbackResp
}

// CallbackAfterQuitGroupReq ˳Ⱥص
type CallbackAfterQuitGroupReq struct {
	CallbackCommand string `json:"callbackCommand"`
	GroupID         string `json:"groupID"`
	UserID          string `json:"userID"`
}

func (c *CallbackAfterQuitGroupReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterQuitGroupResp ˳ȺصӦ
type CallbackAfterQuitGroupResp struct {
	CommonCallbackResp
}

// CallbackAfterKickGroupMemberReq ߳ȺԱص
type CallbackAfterKickGroupMemberReq struct {
	CallbackCommand string   `json:"callbackCommand"`
	GroupID         string   `json:"groupID"`
	KickedUserIDs   []string `json:"kickedUserIDs"`
	OpUserID        string   `json:"opUserID"`
}

func (c *CallbackAfterKickGroupMemberReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterKickGroupMemberResp ߳ȺԱصӦ
type CallbackAfterKickGroupMemberResp struct {
	CommonCallbackResp
}

// CallbackAfterDismissGroupReq ɢȺص??
type CallbackAfterDismissGroupReq struct {
	CallbackCommand string `json:"callbackCommand"`
	GroupID         string `json:"groupID"`
	OpUserID        string `json:"opUserID"`
}

func (c *CallbackAfterDismissGroupReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterDismissGroupResp ɢȺص??
type CallbackAfterDismissGroupResp struct {
	CommonCallbackResp
}

// CallbackAfterTransferGroupOwnerReq תȺص??
type CallbackAfterTransferGroupOwnerReq struct {
	CallbackCommand string `json:"callbackCommand"`
	GroupID         string `json:"groupID"`
	OldOwnerUserID  string `json:"oldOwnerUserID"`
	NewOwnerUserID  string `json:"newOwnerUserID"`
}

func (c *CallbackAfterTransferGroupOwnerReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackAfterTransferGroupOwnerResp תȺص??
type CallbackAfterTransferGroupOwnerResp struct {
	CommonCallbackResp
}

// CallbackBeforeInviteJoinGroupReq Ⱥǰص
type CallbackBeforeInviteJoinGroupReq struct {
	CallbackCommand string   `json:"callbackCommand"`
	GroupID         string   `json:"groupID"`
	InvitedUserIDs  []string `json:"invitedUserIDs"`
	Reason          string   `json:"reason"`
}

func (c *CallbackBeforeInviteJoinGroupReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

// CallbackBeforeInviteJoinGroupResp ȺǰصӦ
type CallbackBeforeInviteJoinGroupResp struct {
	CommonCallbackResp
	RefusedMembersAccount []string `json:"refusedMembersAccount,omitempty"`
}

// CallbackAfterJoinGroupReq Ⱥص??
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

// CallbackAfterJoinGroupResp Ⱥص??
type CallbackAfterJoinGroupResp struct {
	CommonCallbackResp
}
