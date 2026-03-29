package logic

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"fdim/pkg/database"
	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GenInvitationCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenInvitationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenInvitationCodeLogic {
	return &GenInvitationCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenInvitationCodeLogic) GenInvitationCode(req *admin.GenInvitationCodeReq) (*admin.GenInvitationCodeResp, error) {
	// 1. ֤
	if req.Len <= 0 {
		req.Len = 8 // Ĭϳ
	}
	if req.Num <= 0 {
		req.Num = 1 // Ĭ
	}
	if req.Chars == "" {
		req.Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ" // Ĭַ
	}

	// 2. 
	codes := make([]string, 0, req.Num)
	charsLen := big.NewInt(int64(len(req.Chars)))

	for i := int32(0); i < req.Num; i++ {
		var code string
		for j := int32(0); j < req.Len; j++ {
			n, err := rand.Int(rand.Reader, charsLen)
			if err != nil {
				return nil, errs.WrapMsg(err, "failed to generate random number")
			}
			code += string(req.Chars[n.Int64()])
		}
		codes = append(codes, code)
	}

	// 3. ǷѴڣظ
	exists, err := l.svcCtx.AdminDB.FindInvitationRegister(l.ctx, codes)
	if err != nil {
		l.Errorf("FindInvitationRegister failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to check invitation codes")
	}
	if len(exists) > 0 {
		// ظɣ򻯴ʵӦԣ
		return nil, fmt.Errorf("generated codes conflict with existing codes, please try again")
	}

	// 4. б
	invitations := make([]*database.InvitationRegister, len(codes))
	now := time.Now()
	for i, code := range codes {
		invitations[i] = &database.InvitationRegister{
			InvitationCode: code,
			UsedByUserID:   "",
			CreateTime:     now,
		}
	}

	if err := l.svcCtx.AdminDB.AddInvitationCode(l.ctx, invitations); err != nil {
		l.Errorf("AddInvitationCode failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to add generated invitation codes")
	}

	l.Infof("Generated invitation codes: count=%d, len=%d", req.Num, req.Len)
	return &admin.GenInvitationCodeResp{}, nil
}
