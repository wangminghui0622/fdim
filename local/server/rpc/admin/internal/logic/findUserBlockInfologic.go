package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserBlockInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserBlockInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserBlockInfoLogic {
	return &FindUserBlockInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserBlockInfoLogic) FindUserBlockInfo(req *admin.FindUserBlockInfoReq) (*admin.FindUserBlockInfoResp, error) {
	// 1. 验证参数
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("userIDs cannot be empty")
	}

	// 2. 查找用户封禁信息
	blocks, err := l.svcCtx.AdminDB.FindUserBlockInfo(l.ctx, req.UserIDs)
	if err != nil {
		l.Errorf("FindUserBlockInfo failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to find user block info")
	}

	// 3. 转换为响应格式
	results := make([]*admin.BlockInfo, 0, len(blocks))
	for _, block := range blocks {
		results = append(results, &admin.BlockInfo{
			UserID:     block.UserID,
			Reason:     block.Reason,
			OpUserID:   block.OperatorUserID,
			CreateTime: block.CreateTime.UnixMilli(),
		})
	}

	return &admin.FindUserBlockInfoResp{
		Blocks: results,
	}, nil
}
