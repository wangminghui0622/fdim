package constant

// 用户级别常量
const (
	// IMOrdinaryUser 普通用户级?
	IMOrdinaryUser int32 = 1
	// AppOrdinaryUsers 普通应用用户级?
	AppOrdinaryUsers int32 = 1
	// AppAdmin 应用管理员级?
	AppAdmin int32 = 100
	// AppNotificationAdmin 通知管理员级?
	AppNotificationAdmin int32 = 80
)

// 群组角色级别常量
const (
	// GroupOwner 群主
	GroupOwner int32 = 100
	// GroupAdmin 群管理员
	GroupAdmin int32 = 60
	// GroupOrdinaryUsers 普通群成员
	GroupOrdinaryUsers int32 = 20
)

// 群组状态常?
const (
	// GroupStatusNormal 正常状?
	GroupStatusNormal int32 = 0
	// GroupStatusMuted 禁言状?
	GroupStatusMuted int32 = 1
	// GroupStatusDismissed 已解?
	GroupStatusDismissed int32 = 3
)

// 好友申请处理结果常量
const (
	// FriendRequestUnhandled 未处?
	FriendRequestUnhandled int32 = 0
	// FriendRequestAgree 同意
	FriendRequestAgree int32 = 1
	// FriendRequestRefuse 拒绝
	FriendRequestRefuse int32 = -1
)

// 群组申请处理结果常量
const (
	// GroupRequestUnhandled 未处?
	GroupRequestUnhandled int32 = 0
	// GroupRequestAgree 同意
	GroupRequestAgree int32 = 1
	// GroupRequestRefuse 拒绝
	GroupRequestRefuse int32 = -1
)

// 好友添加来源常量
const (
	// BecomeFriendByImport 通过导入添加好友
	BecomeFriendByImport int32 = 1
	// BecomeFriendByApply 通过申请添加好友
	BecomeFriendByApply int32 = 2
)

// 群组类型常量
const (
	// WorkingGroup 工作群组
	WorkingGroup int32 = 2
)

// 会话类型常量
const (
	// SingleChatType 单聊
	SingleChatType int32 = 1
	// WriteGroupChatType 可写群聊（暂未启用）
	WriteGroupChatType int32 = 2
	// ReadGroupChatType 只读群聊（超级群?
	ReadGroupChatType int32 = 3
	// NotificationChatType 通知会话
	NotificationChatType int32 = 4
	// SuperGroupChatType 超级群聊（等同于 ReadGroupChatType?
	SuperGroupChatType int32 = 3
)

// 加入群组来源常量
const (
	// JoinByInvitation 通过邀请加??
	JoinByInvitation int32 = 1
	// JoinBySearch 通过搜索加入
	JoinBySearch int32 = 2
	// JoinByQRCode 通过二维码加??
	JoinByQRCode int32 = 3
)
