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

type GetGroupAllMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupAllMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupAllMemberLogic {
	return &GetGroupAllMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupAllMemberLogic) GetGroupAllMember(req *user.GetGroupAllMemberReq) (*user.GetGroupAllMemberResp, error) {
	resp := &user.GetGroupAllMemberResp{}

	// ������֤
	if req.GroupID == "" {
		return nil, fmt.Errorf("groupID is empty")
	}

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

	return resp, nil
}
