package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/silenceper/wechat/v2/miniprogram/qrcode"
	"github.com/silenceper/wechat/v2/miniprogram/subscribe"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/config"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/wechat"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
	"go.uber.org/zap"
)

// WechatService 微信服务
type WechatService struct {
	wechatClient *wechat.Client
	config       *config.Config
	logger       *zap.Logger
}

// NewWechatService 创建微信服务实例
func NewWechatService(wechatClient *wechat.Client, cfg *config.Config, logger *zap.Logger) *WechatService {
	return &WechatService{
		wechatClient: wechatClient,
		config:       cfg,
		logger:       logger,
	}
}

// GenerateQRCode 生成小程序码
// scene: 场景值，最多32个字符，格式如 "c=ABC123"
// page: 小程序页面路径，如 "pages/baby/join/join"
// 返回: 小程序码图片的URL路径
func (s *WechatService) GenerateQRCode(ctx context.Context, scene, page string) (string, error) {
	s.logger.Info("🚀 [WechatService.GenerateQRCode] START - 开始生成小程序码",
		zap.String("scene", scene),
		zap.String("page", page),
	)

	// 验证参数
	if scene == "" {
		return "", errors.New(errors.ParamError, "scene参数不能为空")
	}
	if len(scene) > 32 {
		return "", errors.New(errors.ParamError, fmt.Sprintf("scene参数长度不能超过32个字符，当前长度: %d", len(scene)))
	}
	if page == "" {
		return "", errors.New(errors.ParamError, "page参数不能为空")
	}

	// 获取小程序二维码实例
	miniProgram := s.wechatClient.GetMiniProgram()
	qrcodeService := miniProgram.GetQRCode()

	// 构建二维码参数
	qrcodeParams := qrcode.QRCoder{
		Scene: scene,
		Page:  page,
		Width: 280, // 二维码宽度(像素)
	}

	s.logger.Info("📦 [WechatService.GenerateQRCode] 调用微信API生成小程序码",
		zap.Any("params", qrcodeParams),
	)

	// 调用微信API生成小程序码 (返回二进制流)
	imageBytes, err := qrcodeService.GetWXACodeUnlimit(qrcodeParams)
	if err != nil {
		s.logger.Error("❌ [WechatService.GenerateQRCode] 调用微信API失败",
			zap.Error(err),
			zap.String("scene", scene),
			zap.String("page", page),
		)
		return "", errors.Wrap(errors.InternalError, "生成小程序码失败", err)
	}

	s.logger.Info("✅ [WechatService.GenerateQRCode] 微信API调用成功，图片大小",
		zap.Int("bytes", len(imageBytes)),
	)

	// 生成文件名 (使用 scene 作为文件名的一部分)
	filename := fmt.Sprintf("qrcode_%s.png", scene)

	// 构建存储路径 (使用配置中的 upload.storage_path)
	storagePath := s.config.Upload.StoragePath
	if storagePath == "" {
		storagePath = "uploads/"
	}

	// 创建二维码子目录
	qrcodePath := filepath.Join(storagePath, "qrcodes")
	if err := os.MkdirAll(qrcodePath, 0755); err != nil {
		s.logger.Error("❌ [WechatService.GenerateQRCode] 创建目录失败",
			zap.Error(err),
			zap.String("path", qrcodePath),
		)
		return "", errors.Wrap(errors.InternalError, "创建存储目录失败", err)
	}

	// 完整文件路径
	filePath := filepath.Join(qrcodePath, filename)

	// 保存图片到文件
	if err := os.WriteFile(filePath, imageBytes, 0644); err != nil {
		s.logger.Error("❌ [WechatService.GenerateQRCode] 保存图片失败",
			zap.Error(err),
			zap.String("filePath", filePath),
		)
		return "", errors.Wrap(errors.InternalError, "保存小程序码图片失败", err)
	}

	s.logger.Info("✅ [WechatService.GenerateQRCode] 小程序码保存成功",
		zap.String("filePath", filePath),
	)

	// 拼接完整的URL访问地址
	relativePath := fmt.Sprintf("/uploads/qrcodes/%s", filename)

	// 从配置中获取服务器基础URL
	baseURL := s.config.Server.BaseURL
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost:%d", s.config.Server.Port)
	}

	// 返回完整的URL路径
	imageURL := fmt.Sprintf("%s%s", baseURL, relativePath)

	s.logger.Info("🎉 [WechatService.GenerateQRCode] 小程序码生成完成",
		zap.String("imageURL", imageURL),
	)

	return imageURL, nil
}

// SendSubscribeMessage 发送订阅消息
func (s *WechatService) SendSubscribeMessage(
	openid string,
	templateID string,
	data map[string]any,
	page string,
	miniprogramState string,
) error {
	s.logger.Info("🚀 [WechatService.SendSubscribeMessage] START - 开始发送微信订阅消息",
		zap.String("openid", openid),
		zap.String("templateID", templateID),
		zap.String("page", page),
		zap.String("miniprogramState", miniprogramState),
		zap.Any("data", data),
	)

	// 获取小程序订阅消息实例
	miniProgram := s.wechatClient.GetMiniProgram()
	subscribeService := miniProgram.GetSubscribe()

	// 格式化数据为 SDK 要求的格式
	formattedData := make(map[string]*subscribe.DataItem)
	for k, v := range data {
		formattedData[k] = &subscribe.DataItem{
			Value: v,
		}
	}

	// 构造消息
	msg := &subscribe.Message{
		ToUser:           openid,
		TemplateID:       templateID,
		Page:             page,
		Data:             formattedData,
		MiniprogramState: miniprogramState,
		Lang:             "zh_CN",
	}

	s.logger.Info("📦 [WechatService.SendSubscribeMessage] 发送消息",
		zap.Any("message", msg),
	)

	// 发送订阅消息
	err := subscribeService.Send(msg)
	if err != nil {
		s.logger.Error("❌ [WechatService.SendSubscribeMessage] 发送失败",
			zap.Error(err),
			zap.String("openid", openid),
			zap.String("templateID", templateID),
		)
		return err
	}

	s.logger.Info("✅ [WechatService.SendSubscribeMessage] 订阅消息发送成功",
		zap.String("openid", openid),
		zap.String("templateID", templateID),
	)

	return nil
}
