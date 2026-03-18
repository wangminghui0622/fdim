package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/convert"
	"fdim/pkg/mcontext"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberListLogic {
	return &GetGroupMemberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupMemberListLogic) GetGroupMemberList(req *user.GetGroupMemberListReq) (*user.GetGroupMemberListResp, error) {
	resp := &user.GetGroupMemberListResp{}

	// Ȩ�޼�飺��Ҫ��Ⱥ��Ա�����Ա
	if !authverify.IsAdmin(l.ctx) {
		opUserID := mcontext.GetOpUserID(l.ctx)
		_, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, opUserID)
		if err != nil {
			return nil, fmt.Errorf("user not in group")
		}
	}

	// ��ȡ����Ⱥ��Ա
	members, err := l.svcCtx.GroupDB.FindGroupMemberAll(l.ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	// 转换为 Protocol Buffer 格式
	resp.Members = make([]*sdkws.GroupMemberFullInfo, 0, len(members))
	for _, m := range members {
		pbMember := convert.ModelGroupMemberDB2Pb(m)
		resp.Members = append(resp.Members, pbMember)
	}

	resp.Total = uint32(len(members))
	return resp, nil
}
