package logic

import (
	"context"
	"strconv"
	"sync"
	"time"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/pkg/errs"
	"fdim/protocol/constant"
	"fdim/protocol/msggateway"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type GetUsersOnlineStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUsersOnlineStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersOnlineStatusLogic {
	return &GetUsersOnlineStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsersOnlineStatusLogic) GetUsersOnlineStatus(req *types.GetUsersOnlineStatusReq) (resp *types.GetUsersOnlineStatusResp, err error) {
	if l.svcCtx.MsgGatewayClient == nil {
		return nil, errs.ErrInternalServer.WrapMsg("msgGateway service is not configured")
	}

	// 转换请求参数
	rpcReq := &msggateway.GetUsersOnlineStatusReq{
		UserIDs: req.UserIDs,
	}

	// 官方实现：open-im-server/internal/api/user.go GetUsersOnlineStatus
	allResults, err := l.queryAllGatewayInstances(rpcReq)
	if err != nil {
		logx.Errorf("Failed to query all gateway instances: %v", err)
		// 降级：使用单实例查询
		return l.querySingleInstance(rpcReq)
	}

	// 合并所有实例的结果
	return l.mergeResults(req.UserIDs, allResults), nil
}

func (l *GetUsersOnlineStatusLogic) queryAllGatewayInstances(req *msggateway.GetUsersOnlineStatusReq) ([]*msggateway.GetUsersOnlineStatusResp, error) {
	// 获取 etcd 配置
	etcdConfig := l.svcCtx.Config.MessageGatewayRpc.Etcd
	if len(etcdConfig.Hosts) == 0 {
		return nil, errs.ErrInternalServer.WrapMsg("etcd hosts not configured")
	}

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdConfig.Hosts,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, errs.ErrInternalServer.WrapMsg("failed to create etcd client: " + err.Error())
	}
	defer etcdClient.Close()

	// 获取服务实例列表
	ctx, cancel := context.WithTimeout(l.ctx, 3*time.Second)
	defer cancel()

	// 构建 etcd key
	serviceKey := etcdConfig.Key
	if serviceKey == "" {
		serviceKey = "msggateway.rpc"
	}

	getResp, err := etcdClient.Get(ctx, serviceKey, clientv3.WithPrefix())
	if err != nil {
		return nil, errs.ErrInternalServer.WrapMsg("failed to get service instances from etcd: " + err.Error())
	}

	if len(getResp.Kvs) == 0 {
		return nil, errs.ErrInternalServer.WrapMsg("no msggateway instances found in etcd")
	}

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []*msggateway.GetUsersOnlineStatusResp
	)

	// 限制并发数，避免过多连接
	semaphore := make(chan struct{}, 10)

	for _, kv := range getResp.Kvs {
		wg.Add(1)
		go func(endpoint string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			queryCtx, queryCancel := context.WithTimeout(l.ctx, 2*time.Second)
			defer queryCancel()

			conn, err := grpc.DialContext(queryCtx, endpoint,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
			)
			if err != nil {
				logx.Errorf("Failed to connect to msggateway instance %s: %v", endpoint, err)
				return
			}
			defer conn.Close()

			client := msggateway.NewMsgGatewayClient(conn)
			resp, err := client.GetUsersOnlineStatus(queryCtx, req)
			if err != nil {
				logx.Errorf("Failed to query msggateway instance %s: %v", endpoint, err)
				return
			}

			mu.Lock()
			results = append(results, resp)
			mu.Unlock()
		}(string(kv.Value))
	}

	wg.Wait()

	if len(results) == 0 {
		return nil, errs.ErrInternalServer.WrapMsg("all msggateway instances query failed")
	}

	return results, nil
}

func (l *GetUsersOnlineStatusLogic) querySingleInstance(req *msggateway.GetUsersOnlineStatusReq) (*types.GetUsersOnlineStatusResp, error) {
	rpcResp, err := l.svcCtx.MsgGatewayClient.GetUsersOnlineStatus(l.ctx, req)
	if err != nil {
		return nil, err
	}

	return l.mergeResults(req.UserIDs, []*msggateway.GetUsersOnlineStatusResp{rpcResp}), nil
}

func (l *GetUsersOnlineStatusLogic) mergeResults(userIDs []string, allResults []*msggateway.GetUsersOnlineStatusResp) *types.GetUsersOnlineStatusResp {
	userStatusMap := make(map[string]*types.UserOnlineStatus)

	for _, userID := range userIDs {
		userStatusMap[userID] = &types.UserOnlineStatus{
			UserID: userID,
			Status: constant.Offline,
		}
	}

	// 合并所有实例的结果
	for _, result := range allResults {
		for _, wsRes := range result.SuccessResult {
			userID := wsRes.UserID
			if status, exists := userStatusMap[userID]; exists {
				status.Status = constant.Online

				// 合并平台详情
				for _, detail := range wsRes.DetailPlatformStatus {
					status.DetailPlatformStatus = append(status.DetailPlatformStatus, types.UserOnlinePlatformStatus{
						Platform: strconv.Itoa(int(detail.PlatformID)),
						Status:   constant.Online, // DetailPlatformStatus 中的每个条目都表示在?
					})
				}
			}
		}
	}

	var successResult []types.UserOnlineStatus
	for _, userID := range userIDs {
		if status, exists := userStatusMap[userID]; exists {
			successResult = append(successResult, *status)
		}
	}

	return &types.GetUsersOnlineStatusResp{
		SuccessResult: successResult,
		FailedResult:  []interface{}{},
	}
}
