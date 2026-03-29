package rtc

import (
	"encoding/json"
	"sync"
)

// CallType 通话类型
type CallType int

const (
	CallTypeVoice CallType = 1 // 语音通话
	CallTypeVideo CallType = 2 // 视频通话
)

// SignalType 信令类型
type SignalType string

const (
	SignalTypeOffer     SignalType = "offer"     // SDP Offer
	SignalTypeAnswer    SignalType = "answer"    // SDP Answer
	SignalTypeCandidate SignalType = "candidate" // ICE Candidate
	SignalTypeInvite    SignalType = "invite"    // 邀请通话
	SignalTypeAccept    SignalType = "accept"    // 接受通话
	SignalTypeReject    SignalType = "reject"    // 拒绝通话
	SignalTypeCancel    SignalType = "cancel"    // 取消通话
	SignalTypeHangup    SignalType = "hangup"    // 挂断通话
	SignalTypeBusy      SignalType = "busy"      // 忙线
)

// SignalMessage 信令消息
type SignalMessage struct {
	Type       SignalType `json:"type"`
	FromUserID string     `json:"fromUserID"`
	ToUserID   string     `json:"toUserID"`
	CallType   CallType   `json:"callType"`
	RoomID     string     `json:"roomID"`
	SDP        string     `json:"sdp,omitempty"`
	Candidate  string     `json:"candidate,omitempty"`
	Timestamp  int64      `json:"timestamp"`
}

// CallSession 通话会话
type CallSession struct {
	RoomID     string   `json:"roomID"`
	CallerID   string   `json:"callerID"`
	CalleeID   string   `json:"calleeID"`
	CallType   CallType `json:"callType"`
	Status     string   `json:"status"` // waiting, connected, ended
	StartTime  int64    `json:"startTime"`
	AnswerTime int64    `json:"answerTime"`
	EndTime    int64    `json:"endTime"`
}

// SignalingServer 信令服务器
type SignalingServer struct {
	sessions map[string]*CallSession // roomID -> session
	userRoom map[string]string       // userID -> roomID
	mutex    sync.RWMutex
}

// NewSignalingServer 创建信令服务器
func NewSignalingServer() *SignalingServer {
	return &SignalingServer{
		sessions: make(map[string]*CallSession),
		userRoom: make(map[string]string),
	}
}

// CreateSession 创建通话会话
func (s *SignalingServer) CreateSession(callerID, calleeID string, callType CallType) *CallSession {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	roomID := generateRoomID(callerID, calleeID)
	session := &CallSession{
		RoomID:   roomID,
		CallerID: callerID,
		CalleeID: calleeID,
		CallType: callType,
		Status:   "waiting",
	}
	s.sessions[roomID] = session
	s.userRoom[callerID] = roomID
	s.userRoom[calleeID] = roomID
	return session
}

// GetSession 获取通话会话
func (s *SignalingServer) GetSession(roomID string) *CallSession {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.sessions[roomID]
}

// GetUserSession 获取用户当前的通话会话
func (s *SignalingServer) GetUserSession(userID string) *CallSession {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if roomID, ok := s.userRoom[userID]; ok {
		return s.sessions[roomID]
	}
	return nil
}

// UpdateSessionStatus 更新会话状态
func (s *SignalingServer) UpdateSessionStatus(roomID, status string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if session, ok := s.sessions[roomID]; ok {
		session.Status = status
	}
}

// EndSession 结束通话会话
func (s *SignalingServer) EndSession(roomID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if session, ok := s.sessions[roomID]; ok {
		delete(s.userRoom, session.CallerID)
		delete(s.userRoom, session.CalleeID)
		delete(s.sessions, roomID)
	}
}

// IsUserBusy 检查用户是否忙线
func (s *SignalingServer) IsUserBusy(userID string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	_, ok := s.userRoom[userID]
	return ok
}

// generateRoomID 生成房间ID
func generateRoomID(userID1, userID2 string) string {
	if userID1 > userID2 {
		userID1, userID2 = userID2, userID1
	}
	return userID1 + "_" + userID2
}

// ParseSignalMessage 解析信令消息
func ParseSignalMessage(data []byte) (*SignalMessage, error) {
	var msg SignalMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// ToJSON 转换为JSON
func (m *SignalMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// DefaultSignalingServer 默认信令服务器实
var DefaultSignalingServer = NewSignalingServer()
