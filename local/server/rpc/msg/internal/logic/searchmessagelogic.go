package logic

import (
	"context"
	"strconv"
	"time"

	"fdim/pkg/errs"
	"fdim/pkg/model"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SearchMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchMessageLogic {
	return &SearchMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SearchMessage 在消息存储中按简单条件搜索聊天记录?
// 为了保持实现可控，这里实现一个简化版?
// - 支持?sendID / recvID / contentType / sessionType 过滤
// - 支持?sendTime（时间戳字符串）作为起始时间过滤
// - 支持分页（RequestPagination）和?send_time 倒序排序
func (l *SearchMessageLogic) SearchMessage(req *msg.SearchMessageReq) (*msg.SearchMessageResp, error) {
	coll := l.svcCtx.MongoDB.GetCollection("stream_msg")
	if coll == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message collection not initialized")
	}

	// 构建查询条件
	filter := bson.M{}
	if req.SendID != "" {
		filter["send_id"] = req.SendID
	}
	if req.RecvID != "" {
		filter["recv_id"] = req.RecvID
	}
	if req.ContentType != 0 {
		filter["content_type"] = req.ContentType
	}
	if req.SessionType != 0 {
		filter["session_type"] = req.SessionType
	}
	if req.SendTime != "" {
		// SendTime 约定为毫秒时间戳字符?
		if ts, err := timeFromMillisString(req.SendTime); err == nil {
			filter["send_time"] = bson.M{"$gte": ts.UnixMilli()}
		}
	}

	// 分页参数
	var (
		skip  int64 = 0
		limit int64 = 20
	)
	if req.Pagination != nil {
		if req.Pagination.PageNumber > 0 && req.Pagination.ShowNumber > 0 {
			skip = int64((req.Pagination.PageNumber - 1) * req.Pagination.ShowNumber)
			limit = int64(req.Pagination.ShowNumber)
		}
	}

	// 查询总数
	total, err := coll.CountDocuments(l.ctx, filter)
	if err != nil {
		l.Errorw("CountDocuments failed", logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to count messages")
	}

	// 查询数据，按 send_time 倒序
	findOpts := options.Find().
		SetSort(bson.D{{Key: "send_time", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := coll.Find(l.ctx, filter, findOpts)
	if err != nil {
		l.Errorw("Find failed", logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to search messages")
	}
	defer cursor.Close(l.ctx)

	var docs []*model.MsgDoc
	if err := cursor.All(l.ctx, &docs); err != nil {
		l.Errorw("Decode messages failed", logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to decode messages")
	}

	// 转为 SearchChatLog 列表
	resp := &msg.SearchMessageResp{
		ChatLogs:    make([]*msg.SearchChatLog, 0, len(docs)),
		ChatLogsNum: int32(total),
	}

	for _, d := range docs {
		cl := &msg.ChatLog{
			ServerMsgID:      d.ServerMsgID,
			ClientMsgID:      d.ClientMsgID,
			SendID:           d.SendID,
			RecvID:           d.RecvID,
			GroupID:          d.GroupID,
			RecvNickname:     "", // 简化：不反查昵?
			SenderPlatformID: d.SenderPlatformID,
			SenderNickname:   d.SenderNickname,
			SenderFaceURL:    d.SenderFaceURL,
			GroupName:        "",
			SessionType:      d.SessionType,
			MsgFrom:          d.MsgFrom,
			ContentType:      d.ContentType,
			Content:          string(d.Content),
			Status:           d.Status,
			SendTime:         d.SendTime,
			CreateTime:       d.CreateTime,  // 已经?int64 时间?
			Ex:               d.Ex,
			GroupFaceURL:     "",
			GroupMemberCount: 0,
			Seq:              d.Seq,
			GroupOwner:       "",
			GroupType:        0,
		}
		resp.ChatLogs = append(resp.ChatLogs, &msg.SearchChatLog{
			ChatLog:   cl,
			IsRevoked: false,
		})
	}

	return resp, nil
}

// timeFromMillisString 将毫秒时间戳的字符串转换?time.Time
func timeFromMillisString(s string) (time.Time, error) {
	// 允许空字符串在外层忽略错?
	if s == "" {
		return time.Time{}, errs.ErrArgs.WrapMsg("empty time string")
	}
	ms, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.UnixMilli(ms), nil
}
