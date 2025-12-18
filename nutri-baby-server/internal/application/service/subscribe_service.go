package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/wxlbd/nutri-baby-server/internal/application/dto"
	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	errs "github.com/wxlbd/nutri-baby-server/pkg/errors"
	"go.uber.org/zap"
)

type SubscribeService struct {
	subscribeRepo         repository.SubscribeRepository
	subscriptionCacheRepo repository.SubscriptionCacheRepository
	userRepo              repository.UserRepository
	wechatService         *WechatService
	logger                *zap.Logger
}

func NewSubscribeService(
	subscribeRepo repository.SubscribeRepository,
	subscriptionCacheRepo repository.SubscriptionCacheRepository,
	userRepo repository.UserRepository,
	wechatService *WechatService,
	logger *zap.Logger,
) *SubscribeService {
	return &SubscribeService{
		subscribeRepo:         subscribeRepo,
		subscriptionCacheRepo: subscriptionCacheRepo,
		userRepo:              userRepo,
		wechatService:         wechatService,
		logger:                logger,
	}
}

// SaveSubscribeAuth 保存用户授权记录(一次性订阅消息机制回调处理)
//
// 微信授权完成后,通过回调通知应用授权结果
// 用户可选"记住我的选择",如勾选则将权限状态缓存30天
func (s *SubscribeService) SaveSubscribeAuth(ctx context.Context, openid string, records []dto.SubscribeAuthDTO) (*dto.SubscribeAuthResponse, error) {
	successCount := 0
	failedCount := 0

	for _, r := range records {
		// 判断用户是同意还是拒绝
		status := repository.StatusDeny
		if r.Status == "accept" {
			status = repository.StatusAllow
		}

		// 将权限状态保存到 Redis
		// 注: 微信的"总是保持以上选择"由微信端实现,我们只需记录状态即可(永久有效)
		err := s.subscriptionCacheRepo.SetSubscriptionStatus(
			ctx,
			openid,
			r.TemplateType,
			status,
		)
		if err != nil {
			s.logger.Error("Failed to cache subscription status",
				zap.String("openid", openid),
				zap.String("templateType", r.TemplateType),
				zap.Error(err),
			)
			failedCount++
		} else {
			successCount++
			s.logger.Info("Subscription status cached",
				zap.String("openid", openid),
				zap.String("templateType", r.TemplateType),
				zap.String("status", string(status)),
			)
		}
	}

	return &dto.SubscribeAuthResponse{
		SuccessCount: successCount,
		FailedCount:  failedCount,
	}, nil
}

// GetUserSubscriptions 获取用户订阅状态
//
// 返回用户在 Redis 缓存中存储的所有订阅权限记录
func (s *SubscribeService) GetUserSubscriptions(ctx context.Context, openid string) (*dto.SubscribeStatusResponse, error) {
	// 从缓存中获取用户的所有订阅权限状态
	subscriptionMap, err := s.subscriptionCacheRepo.GetAllSubscriptions(ctx, openid)
	if err != nil {
		s.logger.Error("Failed to get user subscriptions from cache",
			zap.String("openid", openid),
			zap.Error(err),
		)
		return nil, err
	}

	// 转换为响应格式
	subscriptions := make([]dto.SubscriptionItem, 0, len(subscriptionMap))
	for templateType, status := range subscriptionMap {
		item := dto.SubscriptionItem{
			TemplateType:  templateType,
			Status:        string(status),
			SubscribeTime: time.Now().Unix(), // 缓存中未记录精确时间,使用当前时间
		}
		subscriptions = append(subscriptions, item)
	}

	return &dto.SubscribeStatusResponse{
		Subscriptions: subscriptions,
	}, nil
}

// CheckAuthorizationStatus 检查用户对特定模板的授权状态
//
// 返回 true 当且仅当:
//  1. 缓存中有该记录且状态为 allow(允许)
//  2. 缓存未命中时,回复 false(需要重新授权询问)
func (s *SubscribeService) CheckAuthorizationStatus(ctx context.Context, openid, templateType string) (bool, error) {
	// 首先从缓存中查询用户的权限状态
	allowed, err := s.subscriptionCacheRepo.HasAllowedTemplate(ctx, openid, templateType)
	if err != nil {
		s.logger.Error("Failed to check subscription cache",
			zap.String("openid", openid),
			zap.String("templateType", templateType),
			zap.Error(err),
		)
		// 缓存查询失败,保守起见返回 false(需要重新授权)
		return false, nil
	}

	// 缓存命中且用户已授权
	if allowed {
		s.logger.Info("User has allowed subscription",
			zap.String("openid", openid),
			zap.String("templateType", templateType),
		)
		return true, nil
	}

	// 检查用户是否已明确拒绝
	denied, err := s.subscriptionCacheRepo.HasDeniedTemplate(ctx, openid, templateType)
	if err == nil && denied {
		s.logger.Info("User has denied subscription",
			zap.String("openid", openid),
			zap.String("templateType", templateType),
		)
		return false, nil
	}

	// 缓存未命中,需要向用户显示授权弹窗
	s.logger.Info("Subscription status not in cache, need to request authorization",
		zap.String("openid", openid),
		zap.String("templateType", templateType),
	)
	return false, nil
}

// SendSubscribeMessage 发送订阅消息(一次性消息机制)
func (s *SubscribeService) SendSubscribeMessage(
	ctx context.Context,
	req *dto.SendMessageRequest,
) error {
	s.logger.Info("📤 [SendSubscribeMessage] START - 开始发送订阅消息",
		zap.String("openid", req.OpenID),
		zap.String("page", req.Page),
		zap.Any("data", req.Data),
	)

	s.logger.Info("✅ [SendSubscribeMessage] 授权可用,准备调用微信API")

	// 3. 调用微信API发送
	s.logger.Info("📞 [SendSubscribeMessage] STEP 3 - 调用微信API发送订阅消息",
		zap.String("openid", req.OpenID),
		//zap.String("templateID", record.TemplateID),
		zap.String("page", req.Page),
		zap.Any("data", req.Data),
	)

	err := s.wechatService.SendSubscribeMessage(
		req.OpenID,
		req.TemplateID,
		req.Data,
		req.Page,
		"formal",
	)
	if err != nil {
		s.logger.Error("Failed to send subscribe message", zap.Error(err))
	}

	// 4. 标记授权为已使用(无论发送成功或失败,授权都会被消耗)
	s.logger.Info("🔄 [SendSubscribeMessage] STEP 4 - 标记授权为已使用",
		zap.String("openid", req.OpenID),
		zap.String("templateID", req.TemplateID),
	)

	// 5. 记录发送日志
	s.logger.Info("📝 [SendSubscribeMessage] STEP 5 - 保存发送日志")

	// 获取用户信息以获取UserID
	user, err := s.userRepo.FindByOpenID(ctx, req.OpenID)
	if err != nil {
		s.logger.Error("Failed to find user",
			zap.String("openid", req.OpenID),
			zap.Error(err))
		return err
	}

	dataJSON, _ := json.Marshal(req.Data)
	now := time.Now().UnixMilli()
	log := &entity.MessageSendLog{
		UserID:           user.ID,
		TemplateID:       req.TemplateID,
		Data:             string(dataJSON),
		Page:             req.Page,
		MiniprogramState: "formal",
	}

	if err != nil {
		log.SendStatus = "failed"
		log.ErrMsg = err.Error()
		s.logger.Error("❌ [SendSubscribeMessage] 发送订阅消息失败",
			zap.String("openid", req.OpenID),
			zap.String("templateID", req.TemplateID),
			//zap.String("templateID", record.TemplateID),
			zap.Error(err),
		)
	} else {
		log.SendStatus = "success"
		log.SendTime = &now
		s.logger.Info("✅ [SendSubscribeMessage] 订阅消息发送成功",
			zap.String("openid", req.OpenID),
			zap.String("templateID", req.TemplateID),
			//zap.String("templateID", record.TemplateID),
		)
	}

	logErr := s.subscribeRepo.CreateSendLog(ctx, log)
	if logErr != nil {
		s.logger.Error("❌ [SendSubscribeMessage] 保存发送日志失败",
			zap.Error(logErr),
		)
	} else {
		s.logger.Info("✅ [SendSubscribeMessage] 发送日志已保存")
	}

	s.logger.Info("🏁 [SendSubscribeMessage] END - 订阅消息发送流程结束",
		zap.String("openid", req.OpenID),
		zap.String("templateID", req.TemplateID),
		zap.Bool("success", err == nil),
	)

	return err
}

// GetMessageLogs 获取消息发送日志
func (s *SubscribeService) GetMessageLogs(ctx context.Context, openid string, offset, limit int) (*dto.MessageLogsResponse, error) {
	// 获取用户信息以获取UserID
	user, err := s.userRepo.FindByOpenID(ctx, openid)
	if err != nil {
		s.logger.Error("Failed to find user",
			zap.String("openid", openid),
			zap.Error(err))
		return nil, errs.ErrInternal
	}

	logs, total, err := s.subscribeRepo.GetSendLogs(ctx, user.ID, offset, limit)
	if err != nil {
		s.logger.Error("Failed to get message logs",
			zap.String("openid", openid),
			zap.Error(err),
		)
		return nil, errs.ErrInternal
	}

	items := make([]dto.MessageLogItem, 0, len(logs))
	for _, log := range logs {
		item := dto.MessageLogItem{
			ID:         uint(log.ID),
			SendStatus: log.SendStatus,
			ErrMsg:     log.ErrMsg,
			CreatedAt:  log.CreatedAt / 1000, // 毫秒转秒
		}
		if log.SendTime != nil {
			item.SendTime = *log.SendTime / 1000 // 毫秒转秒
		}
		items = append(items, item)
	}

	return &dto.MessageLogsResponse{
		Logs:  items,
		Total: total,
	}, nil
}
