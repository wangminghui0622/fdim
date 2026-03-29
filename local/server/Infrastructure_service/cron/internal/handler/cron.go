package handler

import (
	"context"
	"fdim/Infrastructure_service/cron/internal/config"
	"fdim/Infrastructure_service/cron/internal/svc"
	"fmt"
	"time"

	"fdim/protocol/conversation"
	"fdim/protocol/msg"
	"fdim/protocol/third"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

// CronHandler Cron 任务处理?
type CronHandler struct {
	config             *config.Config
	cron               *cron.Cron
	msgClient          msg.MsgClient
	conversationClient conversation.ConversationClient
	thirdClient        third.ThirdClient
}

// NewCronHandler 创建 Cron 处理?
func NewCronHandler(ctx context.Context, cfg *config.Config) (*CronHandler, error) {
	svcCtx := svc.NewServiceContext(*cfg)

	handler := &CronHandler{
		config:             cfg,
		cron:               cron.New(cron.WithSeconds()),
		msgClient:          svcCtx.MsgClient,
		conversationClient: svcCtx.ConversationClient,
		thirdClient:        svcCtx.ThirdClient,
	}

	// 注册定时任务
	if err := handler.registerTasks(); err != nil {
		return nil, fmt.Errorf("failed to register tasks: %w", err)
	}

	return handler, nil
}

// registerTasks 注册定时任务
func (h *CronHandler) registerTasks() error {
	// 清理 S3 文件任务（每天凌?2 点执行）
	_, err := h.cron.AddFunc("0 0 2 * * *", func() {
		logx.Info("Running cleanup S3 files task...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		// 调用 ThirdClient 删除过期数据
		_, err := h.thirdClient.DeleteOutdatedData(ctx, &third.DeleteOutdatedDataReq{})
		if err != nil {
			logx.Errorf("Failed to delete outdated data: %v", err)
		} else {
			logx.Info("Cleanup S3 files task completed")
		}
	})
	if err != nil {
		return fmt.Errorf("failed to register cleanup S3 files task: %w", err)
	}

	// 删除过期消息任务（每天凌?3 点执行）
	_, err = h.cron.AddFunc("0 0 3 * * *", func() {
		logx.Info("Running delete expired messages task...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		// 调用 MsgClient 删除过期消息
		_, err := h.msgClient.DestructMsgs(ctx, &msg.DestructMsgsReq{})
		if err != nil {
			logx.Errorf("Failed to destruct messages: %v", err)
		} else {
			logx.Info("Delete expired messages task completed")
		}
	})
	if err != nil {
		return fmt.Errorf("failed to register delete expired messages task: %w", err)
	}

	// 清理用户消息任务（每天凌?4 点执行）
	_, err = h.cron.AddFunc("0 0 4 * * *", func() {
		logx.Info("Running cleanup user conversation messages task...")
		
		// 检?conversationClient 是否?nil
		if h.conversationClient == nil {
			logx.Error("ConversationClient is nil, skipping cleanup task")
			return
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		// 调用 ConversationClient 清理用户会话消息
		_, err := h.conversationClient.ClearUserConversationMsg(ctx, &conversation.ClearUserConversationMsgReq{})
		if err != nil {
			logx.Errorf("Failed to clear user conversation messages: %v", err)
		} else {
			logx.Info("Cleanup user conversation messages task completed")
		}
	})
	if err != nil {
		return fmt.Errorf("failed to register cleanup user conversation messages task: %w", err)
	}

	logx.Info("All cron tasks registered successfully")
	return nil
}

// Start 启动 Cron 服务
func (h *CronHandler) Start(ctx context.Context) error {
	logx.Info("Starting cron service...")
	h.cron.Start()

	// 等待上下文取?
	<-ctx.Done()
	logx.Info("Cron service context cancelled, stopping...")
	h.cron.Stop()
	return nil
}

// Stop 停止 Cron 服务
func (h *CronHandler) Stop(ctx context.Context) error {
	logx.Info("Stopping cron service...")
	h.cron.Stop()
	return nil
}

