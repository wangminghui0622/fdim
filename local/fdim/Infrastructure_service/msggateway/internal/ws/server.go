package ws

import (
	"context"
	"encoding/json"
	"fdim/Infrastructure_service/msggateway/internal/config"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"fdim/protocol/auth"
	"github.com/olahol/melody"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

// WsServer WebSocket 服务器
type WsServer struct {
	melody           *melody.Melody
	config           *config.Config
	clients          *ClientManager
	onlineCache      *OnlineCache
	totalConnections atomic.Int64
	authClient       auth.AuthClient
	ready            bool
	mu               sync.RWMutex
}

// ClientManager 客户端管理器
type ClientManager struct {
	clients map[string]map[string]*melody.Session // userID -> platformID -> session
	mu      sync.RWMutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[string]map[string]*melody.Session),
	}
}

func (cm *ClientManager) Add(userID, platformID string, session *melody.Session) *melody.Session {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.clients[userID] == nil {
		cm.clients[userID] = make(map[string]*melody.Session)
	}
	// 返回旧 session（如果存在），用于同平台踢下线
	old := cm.clients[userID][platformID]
	cm.clients[userID][platformID] = session
	return old
}

func (cm *ClientManager) Remove(userID, platformID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.clients[userID] != nil {
		delete(cm.clients[userID], platformID)
		if len(cm.clients[userID]) == 0 {
			delete(cm.clients, userID)
		}
	}
}

func (cm *ClientManager) Get(userID string) map[string]*melody.Session {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.clients[userID]
}

func (cm *ClientManager) GetAll(userID string) []*melody.Session {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	sessions := make([]*melody.Session, 0)
	if cm.clients[userID] != nil {
		for _, session := range cm.clients[userID] {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

// GetTotalConnections 获取总连接数
func (cm *ClientManager) GetTotalConnections() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	count := 0
	for _, platforms := range cm.clients {
		count += len(platforms)
	}
	return count
}

// OnlineCache 在线状态缓存
type OnlineCache struct {
	redis *redis.Client
}

func NewOnlineCache(redis *redis.Client) *OnlineCache {
	return &OnlineCache{
		redis: redis,
	}
}

// SetUserOnline 设置用户在线状态
func (oc *OnlineCache) SetUserOnline(ctx context.Context, userID string, platformID int32) error {
	if oc.redis == nil {
		return nil
	}
	key := fmt.Sprintf("user:online:%s", userID)
	member := strconv.FormatInt(int64(platformID), 10)
	if err := oc.redis.SAdd(ctx, key, member).Err(); err != nil {
		return err
	}
	// 设置过期时间（24小时）
	return oc.redis.Expire(ctx, key, 24*time.Hour).Err()
}

// SetUserOffline 设置用户离线状态
func (oc *OnlineCache) SetUserOffline(ctx context.Context, userID string, platformID int32) error {
	if oc.redis == nil {
		return nil
	}
	key := fmt.Sprintf("user:online:%s", userID)
	member := strconv.FormatInt(int64(platformID), 10)
	return oc.redis.SRem(ctx, key, member).Err()
}

// GetUserOnlinePlatforms 获取用户在线的平台列表
func (oc *OnlineCache) GetUserOnlinePlatforms(ctx context.Context, userID string) ([]int32, error) {
	if oc.redis == nil {
		return nil, nil
	}
	key := fmt.Sprintf("user:online:%s", userID)
	members, err := oc.redis.SMembers(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return []int32{}, nil
		}
		return nil, err
	}
	platforms := make([]int32, 0, len(members))
	for _, m := range members {
		p, err := strconv.ParseInt(m, 10, 32)
		if err != nil {
			continue
		}
		platforms = append(platforms, int32(p))
	}
	return platforms, nil
}

// NewWsServer 创建 WebSocket 服务器
func NewWsServer(cfg *config.Config, authClient auth.AuthClient, redisClient *redis.Client) *WsServer {
	m := melody.New()

	// 配置 Melody
	m.Config.PingPeriod = 25 * time.Second
	m.Config.PongWait = 60 * time.Second
	m.Config.WriteWait = 10 * time.Second
	m.Config.MaxMessageSize = int64(cfg.LongConnServer.WebsocketMaxMsgLen)
	m.Config.MessageBufferSize = 256

	// 设置 Upgrader
	m.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true // 开发环境允许所有来源
	}

	ws := &WsServer{
		melody:      m,
		config:      cfg,
		clients:     NewClientManager(),
		onlineCache: NewOnlineCache(redisClient),
		authClient:  authClient,
		ready:       true,
	}

	// 设置消息处理
	m.HandleMessage(ws.handleMessage)
	m.HandleConnect(ws.handleConnect)
	m.HandleDisconnect(ws.handleDisconnect)
	m.HandleError(ws.handleError)

	return ws
}

// handleConnect 处理连接
func (ws *WsServer) handleConnect(session *melody.Session) {
	ws.totalConnections.Add(1)
	logx.Infof("[WebSocket Connected] New connection established - remoteAddr=%s, totalConnections=%d",
		session.RemoteAddr(), ws.totalConnections.Load())

	// 从 session 中获取用户信息
	userID, _ := session.Get("userID")
	platformID, _ := session.Get("platformID")

	logx.Infof("[WebSocket Connect] Retrieved from session - userID=%v, platformID=%v", userID, platformID)

	if userIDStr, ok := userID.(string); ok {
		if platformIDStr, ok := platformID.(string); ok {
			logx.Infof("[WebSocket Connect] Adding to clients map - userID=%s, platformID=%s", userIDStr, platformIDStr)
			// 同平台踢下线（官方行为：同一平台只允许一个连接）
			oldSession := ws.clients.Add(userIDStr, platformIDStr, session)
			if oldSession != nil {
				kickResp := WsResp{
					ReqIdentifier: WSKickOnlineMsg,
					ErrCode:       0,
				}
				kickBytes, _ := json.Marshal(kickResp)
				oldSession.Write(kickBytes)
				oldSession.Close()
				logx.Infof("[WebSocket Connect] Kicked old session for same platform - userID=%s, platformID=%s", userIDStr, platformIDStr)
			}
			logx.Infof("[WebSocket Connect] Successfully added to clients map - userID=%s, platformID=%s", userIDStr, platformIDStr)
			
			// 同步在线状态到 Redis
			if platformIDInt, err := strconv.Atoi(platformIDStr); err == nil {
				if err := ws.onlineCache.SetUserOnline(context.Background(), userIDStr, int32(platformIDInt)); err != nil {
					logx.Errorf("[WebSocket Login] Failed to sync online status to Redis - userID=%s, platformID=%s, error=%v",
						userIDStr, platformIDStr, err)
				} else {
					logx.Infof("[WebSocket Login] Online status synced to Redis - userID=%s, platformID=%s",
						userIDStr, platformIDStr)
				}
			}
			
			logx.Infof("[WebSocket Login Success] User logged in - userID=%s, platformID=%s, remoteAddr=%s, totalConnections=%d",
				userIDStr, platformIDStr, session.RemoteAddr(), ws.totalConnections.Load())
		} else {
			logx.Errorf("[WebSocket Connect] platformID type assertion failed - userID=%s, platformID=%v (type=%T)",
				userIDStr, platformID, platformID)
		}
	} else {
		logx.Errorf("[WebSocket Connect] userID type assertion failed - userID=%v (type=%T)", userID, userID)
	}
}

// handleDisconnect 处理断开连接
func (ws *WsServer) handleDisconnect(session *melody.Session) {
	ws.totalConnections.Add(-1)
	logx.Infof("[WebSocket Disconnected] Connection closed - remoteAddr=%s, totalConnections=%d",
		session.RemoteAddr(), ws.totalConnections.Load())

	userID, _ := session.Get("userID")
	platformID, _ := session.Get("platformID")

	logx.Infof("[WebSocket Disconnect] Retrieved from session - userID=%v, platformID=%v", userID, platformID)

	if userIDStr, ok := userID.(string); ok {
		if platformIDStr, ok := platformID.(string); ok {
			logx.Infof("[WebSocket Disconnect] Removing from clients map - userID=%s, platformID=%s", userIDStr, platformIDStr)
			ws.clients.Remove(userIDStr, platformIDStr)
			logx.Infof("[WebSocket Disconnect] Successfully removed from clients map - userID=%s, platformID=%s", userIDStr, platformIDStr)
			
			// 同步离线状态到 Redis
			if platformIDInt, err := strconv.Atoi(platformIDStr); err == nil {
				if err := ws.onlineCache.SetUserOffline(context.Background(), userIDStr, int32(platformIDInt)); err != nil {
					logx.Errorf("[WebSocket Logout] Failed to sync offline status to Redis - userID=%s, platformID=%s, error=%v",
						userIDStr, platformIDStr, err)
				} else {
					logx.Infof("[WebSocket Logout] Offline status synced to Redis - userID=%s, platformID=%s",
						userIDStr, platformIDStr)
				}
			}
			
			logx.Infof("[WebSocket Logout] User disconnected - userID=%s, platformID=%s, remoteAddr=%s, totalConnections=%d",
				userIDStr, platformIDStr, session.RemoteAddr(), ws.totalConnections.Load())
		}
	}
}

// handleMessage 处理消息 (official OpenIM WsReq/WsResp protocol)
func (ws *WsServer) handleMessage(session *melody.Session, msg []byte) {
	logx.Debugf("Received message: %s", string(msg))

	var req WsReq
	if err := json.Unmarshal(msg, &req); err != nil {
		logx.Errorf("Failed to unmarshal WsReq: %v", err)
		return
	}

	switch req.ReqIdentifier {
	case WSSetBackgroundStatus:
		// 心跳/后台状态设置
		resp := WsResp{
			ReqIdentifier: req.ReqIdentifier,
			OperationID:   req.OperationID,
			MsgIncr:       req.MsgIncr,
			ErrCode:       0,
		}
		respBytes, _ := json.Marshal(resp)
		if err := session.Write(respBytes); err != nil {
			logx.Errorf("Failed to send SetBackgroundStatus resp: %v", err)
		}
	case WSGetNewestSeq:
		// 获取最新 seq - 通过 REST API 实现，WS 返回空
		resp := WsResp{
			ReqIdentifier: req.ReqIdentifier,
			OperationID:   req.OperationID,
			MsgIncr:       req.MsgIncr,
			ErrCode:       0,
		}
		respBytes, _ := json.Marshal(resp)
		session.Write(respBytes)
	case WSSendMsg:
		// 发送消息 - 通过 REST API 实现，WS 返回不支持
		resp := WsResp{
			ReqIdentifier: req.ReqIdentifier,
			OperationID:   req.OperationID,
			MsgIncr:       req.MsgIncr,
			ErrCode:       1,
			ErrMsg:        "please use REST API for sending messages",
		}
		respBytes, _ := json.Marshal(resp)
		session.Write(respBytes)
	case WSPullMsgBySeq:
		// 拉取消息 - 通过 REST API 实现
		resp := WsResp{
			ReqIdentifier: req.ReqIdentifier,
			OperationID:   req.OperationID,
			MsgIncr:       req.MsgIncr,
			ErrCode:       1,
			ErrMsg:        "please use REST API for pulling messages",
		}
		respBytes, _ := json.Marshal(resp)
		session.Write(respBytes)
	default:
		logx.Debugf("Unhandled reqIdentifier: %d", req.ReqIdentifier)
	}
}

// handleError 处理错误
func (ws *WsServer) handleError(session *melody.Session, err error) {
	logx.Errorf("WebSocket error: %v", err)
}

// HandleRequest 处理 HTTP 请求并升级为 WebSocket
func (ws *WsServer) HandleRequest(w http.ResponseWriter, r *http.Request) error {
	// 检查连接数限制
	if int64(ws.clients.GetTotalConnections()) >= ws.config.LongConnServer.WebsocketMaxConnNum {
		http.Error(w, "Connection limit exceeded", http.StatusServiceUnavailable)
		return fmt.Errorf("connection limit exceeded")
	}

	// 从请求中解析参数
	userID := r.URL.Query().Get("sendID")
	if userID == "" {
		userID = r.URL.Query().Get("userID")
	}
	platformIDStr := r.URL.Query().Get("platformID")
	token := r.URL.Query().Get("token")
	operationID := r.URL.Query().Get("operationID")

	logx.Infof("[WebSocket Login] Connection request received - userID=%s, platformID=%s, remoteAddr=%s, operationID=%s",
		userID, platformIDStr, r.RemoteAddr, operationID)

	if userID == "" || platformIDStr == "" || token == "" {
		logx.Errorf("[WebSocket Login Failed] Missing parameters - userID=%s, platformID=%s, hasToken=%v, remoteAddr=%s",
			userID, platformIDStr, token != "", r.RemoteAddr)
		http.Error(w, "Missing required parameters: sendID/userID, platformID, token", http.StatusBadRequest)
		return fmt.Errorf("missing required parameters")
	}

	// 解析并验证 platformID
	platformID, err := strconv.Atoi(platformIDStr)
	if err != nil {
		logx.Errorf("[WebSocket Login Failed] Invalid platformID - userID=%s, platformID=%s, remoteAddr=%s, error=%v",
			userID, platformIDStr, r.RemoteAddr, err)
		http.Error(w, "Invalid platformID", http.StatusBadRequest)
		return fmt.Errorf("invalid platformID: %w", err)
	}

	// 验证 Token
	if ws.authClient != nil {
		resp, err := ws.authClient.ParseToken(r.Context(), &auth.ParseTokenReq{
			Token: token,
		})
		if err != nil {
			logx.Errorf("[WebSocket Login Failed] Token validation failed - userID=%s, platformID=%d, remoteAddr=%s, error=%v",
				userID, platformID, r.RemoteAddr, err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return fmt.Errorf("invalid token: %w", err)
		}

		// 验证 token 中的 userID 是否与请求一致
		if resp.UserID != userID {
			logx.Errorf("[WebSocket Login Failed] Token userID mismatch - tokenUserID=%s, requestUserID=%s, platformID=%d, remoteAddr=%s",
				resp.UserID, userID, platformID, r.RemoteAddr)
			http.Error(w, "Token userID mismatch", http.StatusUnauthorized)
			return fmt.Errorf("token userID mismatch: token has %s, request has %s", resp.UserID, userID)
		}
		// 不严格验证 platformID，允许跨平台使用 token
		if resp.PlatformID != int32(platformID) {
			logx.Infof("[WebSocket Login] Token platformID differs: token=%d, request=%d (allowed)", resp.PlatformID, platformID)
		}

		logx.Infof("[WebSocket Login] Token validated successfully - userID=%s, platformID=%d, tokenPlatformID=%d, remoteAddr=%s",
			userID, platformID, resp.PlatformID, r.RemoteAddr)
	} else {
		logx.Infof("AuthClient is nil, skipping token validation")
	}

	// 从 session 中存储用户信息
	keys := make(map[string]interface{})
	keys["userID"] = userID
	keys["platformID"] = platformIDStr
	keys["token"] = token
	if operationID != "" {
		keys["operationID"] = operationID
	}

	return ws.melody.HandleRequestWithKeys(w, r, keys)
}

// BroadcastToUser 向指定用户推送消息
func (ws *WsServer) BroadcastToUser(userID string, message []byte) error {
	logx.Infof("[BroadcastToUser] Attempting to broadcast to user: %s, message size: %d bytes", userID, len(message))
	platformSessions := ws.clients.Get(userID)
	if len(platformSessions) == 0 {
		logx.Infof("[BroadcastToUser] No sessions found for user: %s", userID)
		return fmt.Errorf("user %s not online", userID)
	}

	platformIDs := make([]string, 0, len(platformSessions))
	sessions := make([]*melody.Session, 0, len(platformSessions))
	for platformID, session := range platformSessions {
		platformIDs = append(platformIDs, platformID)
		sessions = append(sessions, session)
	}
	logx.Infof("[BroadcastToUser] Found %d sessions for user: %s, platforms=%v", len(sessions), userID, platformIDs)

	for _, session := range sessions {
		logx.Infof("[BroadcastToUser] Writing to session: %s", session.RemoteAddr())
		logx.Infof("[BroadcastToUser] Message content: %s", string(message))
		if err := session.Write(message); err != nil {
			logx.Errorf("Failed to write message to user %s session %s: %v", userID, session.RemoteAddr(), err)
		} else {
			logx.Infof("[BroadcastToUser] Successfully wrote message to user %s session %s", userID, session.RemoteAddr())
		}
	}
	return nil
}

// BroadcastToUsers 向多个用户推送消息
func (ws *WsServer) BroadcastToUsers(userIDs []string, message []byte) error {
	for _, userID := range userIDs {
		if err := ws.BroadcastToUser(userID, message); err != nil {
			logx.Errorf("Failed to broadcast to user %s: %v", userID, err)
		}
	}
	return nil
}

// BroadcastFilter 按条件广播消息
func (ws *WsServer) BroadcastFilter(message []byte, filter func(*melody.Session) bool) error {
	return ws.melody.BroadcastFilter(message, filter)
}

// KickUserOffline 踢用户下线
func (ws *WsServer) KickUserOffline(userID string, platformID string) error {
	logx.Infof("[WebSocket Kick] Attempting to kick user offline - userID=%s, platformID=%s", userID, platformID)
	
	sessions := ws.clients.Get(userID)
	if sessions == nil {
		logx.Infof("[WebSocket Kick Failed] User not online - userID=%s, platformID=%s", userID, platformID)
		return fmt.Errorf("user %s not online", userID)
	}

	if session, ok := sessions[platformID]; ok {
		// 发送踢下线通知 (official: WSKickOnlineMsg = 2002)
		kickResp := WsResp{
			ReqIdentifier: WSKickOnlineMsg,
			ErrCode:       0,
		}
		kickBytes, _ := json.Marshal(kickResp)
		session.Write(kickBytes)
		logx.Infof("[WebSocket Kick Success] User kicked offline - userID=%s, platformID=%s",
			userID, platformID)
		session.Close()
		ws.clients.Remove(userID, platformID)
		return nil
	}

	logx.Infof("[WebSocket Kick Failed] User not online on specified platform - userID=%s, platformID=%s",
		userID, platformID)
	return fmt.Errorf("user %s not online on platform %s", userID, platformID)
}

// PushMessage 推送消息 (official: WSPushMsg = 2001, all push events use this)
func (ws *WsServer) PushMessage(userID string, data []byte) error {
	resp := WsResp{
		ReqIdentifier: WSPushMsg,
		ErrCode:       0,
		Data:          data,
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return ws.BroadcastToUser(userID, respBytes)
}

// PushMessageToUsers 推送消息给多个用户
func (ws *WsServer) PushMessageToUsers(userIDs []string, data []byte) error {
	resp := WsResp{
		ReqIdentifier: WSPushMsg,
		ErrCode:       0,
		Data:          data,
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return ws.BroadcastToUsers(userIDs, respBytes)
}

// GetOnlineUsers 获取在线用户列表
func (ws *WsServer) GetOnlineUsers(userIDs []string) map[string]bool {
	result := make(map[string]bool)
	for _, userID := range userIDs {
		sessions := ws.clients.Get(userID)
		result[userID] = len(sessions) > 0
	}
	return result
}

// GetUserSessions 获取用户的所有会话
func (ws *WsServer) GetUserSessions(userID string) map[string]*melody.Session {
	return ws.clients.Get(userID)
}

// GetTotalConnections 获取当前总连接数
func (ws *WsServer) GetTotalConnections() int64 {
	return ws.totalConnections.Load()
}

// Run 启动 WebSocket 服务器（阻塞）
func (ws *WsServer) Run(ctx context.Context, port int) error {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		err := ws.HandleRequest(w, r)
		if err != nil {
			logx.Errorf("WebSocket handle request error: %v", err)
		}
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: nil,
	}

	// 监听关闭信号
	go func() {
		<-ctx.Done()
		logx.Info("WebSocket server shutting down...")
		server.Shutdown(context.Background())
	}()

	logx.Infof("WebSocket server listening on port %d", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logx.Errorf("WebSocket server error: %v", err)
		return err
	}
	return nil
}
