package server

import (
	"context"

	"fdim/protocol/third"
	"fdim/rpc/third/internal/logic"
	"fdim/rpc/third/internal/svc"
)

type ThirdServer struct {
	third.UnimplementedThirdServer
	svcCtx *svc.ServiceContext
}

func NewThirdServer(svcCtx *svc.ServiceContext) *ThirdServer {
	return &ThirdServer{
		svcCtx: svcCtx,
	}
}

// PartLimit 获取分片限制
func (s *ThirdServer) PartLimit(ctx context.Context, req *third.PartLimitReq) (*third.PartLimitResp, error) {
	l := logic.NewPartLimitLogic(ctx, s.svcCtx)
	return l.PartLimit(req)
}

// PartSize 获取分片大小
func (s *ThirdServer) PartSize(ctx context.Context, req *third.PartSizeReq) (*third.PartSizeResp, error) {
	l := logic.NewPartSizeLogic(ctx, s.svcCtx)
	return l.PartSize(req)
}

// InitiateMultipartUpload 初始化分片上?
func (s *ThirdServer) InitiateMultipartUpload(ctx context.Context, req *third.InitiateMultipartUploadReq) (*third.InitiateMultipartUploadResp, error) {
	l := logic.NewInitiateMultipartUploadLogic(ctx, s.svcCtx)
	return l.InitiateMultipartUpload(req)
}

// AuthSign 授权签名
func (s *ThirdServer) AuthSign(ctx context.Context, req *third.AuthSignReq) (*third.AuthSignResp, error) {
	l := logic.NewAuthSignLogic(ctx, s.svcCtx)
	return l.AuthSign(req)
}

// CompleteMultipartUpload 完成分片上传
func (s *ThirdServer) CompleteMultipartUpload(ctx context.Context, req *third.CompleteMultipartUploadReq) (*third.CompleteMultipartUploadResp, error) {
	l := logic.NewCompleteMultipartUploadLogic(ctx, s.svcCtx)
	return l.CompleteMultipartUpload(req)
}

// AccessURL 获取访问URL
func (s *ThirdServer) AccessURL(ctx context.Context, req *third.AccessURLReq) (*third.AccessURLResp, error) {
	l := logic.NewAccessURLLogic(ctx, s.svcCtx)
	return l.AccessURL(req)
}

// InitiateFormData 初始化表单数据上?
func (s *ThirdServer) InitiateFormData(ctx context.Context, req *third.InitiateFormDataReq) (*third.InitiateFormDataResp, error) {
	l := logic.NewInitiateFormDataLogic(ctx, s.svcCtx)
	return l.InitiateFormData(req)
}

// CompleteFormData 完成表单数据上传
func (s *ThirdServer) CompleteFormData(ctx context.Context, req *third.CompleteFormDataReq) (*third.CompleteFormDataResp, error) {
	l := logic.NewCompleteFormDataLogic(ctx, s.svcCtx)
	return l.CompleteFormData(req)
}

// DeleteOutdatedData 删除过期数据
func (s *ThirdServer) DeleteOutdatedData(ctx context.Context, req *third.DeleteOutdatedDataReq) (*third.DeleteOutdatedDataResp, error) {
	l := logic.NewDeleteOutdatedDataLogic(ctx, s.svcCtx)
	return l.DeleteOutdatedData(req)
}

// FcmUpdateToken 更新FCM Token
func (s *ThirdServer) FcmUpdateToken(ctx context.Context, req *third.FcmUpdateTokenReq) (*third.FcmUpdateTokenResp, error) {
	l := logic.NewFcmUpdateTokenLogic(ctx, s.svcCtx)
	return l.FcmUpdateToken(req)
}

// SetAppBadge 设置应用角标
func (s *ThirdServer) SetAppBadge(ctx context.Context, req *third.SetAppBadgeReq) (*third.SetAppBadgeResp, error) {
	l := logic.NewSetAppBadgeLogic(ctx, s.svcCtx)
	return l.SetAppBadge(req)
}

// UploadLogs 上传日志
func (s *ThirdServer) UploadLogs(ctx context.Context, req *third.UploadLogsReq) (*third.UploadLogsResp, error) {
	l := logic.NewUploadLogsLogic(ctx, s.svcCtx)
	return l.UploadLogs(req)
}

// DeleteLogs 删除日志
func (s *ThirdServer) DeleteLogs(ctx context.Context, req *third.DeleteLogsReq) (*third.DeleteLogsResp, error) {
	l := logic.NewDeleteLogsLogic(ctx, s.svcCtx)
	return l.DeleteLogs(req)
}

// SearchLogs 搜索日志
func (s *ThirdServer) SearchLogs(ctx context.Context, req *third.SearchLogsReq) (*third.SearchLogsResp, error) {
	l := logic.NewSearchLogsLogic(ctx, s.svcCtx)
	return l.SearchLogs(req)
}
