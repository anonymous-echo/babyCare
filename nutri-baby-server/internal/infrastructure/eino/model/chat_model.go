package model

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/cloudwego/eino/components/model"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/config"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/eino/chain"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/genai"
)

// NewToolCallingChatModel 创建支持工具调用的聊天模型实例
func NewToolCallingChatModel(cfg *config.Config, logger *zap.Logger) (model.ToolCallingChatModel, error) {
	// 使用配置文件中的AI配置
	aiConfig := cfg.AI

	// 默认使用支持工具调用的mock模式进行开发测试
	if aiConfig.Provider == "mock" || aiConfig.Provider == "" {
		logger.Info("使用支持工具调用的模拟AI模型进行开发测试")
		return chain.NewToolCallingMockChatModel(logger), nil
	}

	switch aiConfig.Provider {
	case "openai":
		return NewOpenAIToolCallingChatModel(aiConfig.OpenAI, logger)
	case "claude":
		return NewClaudeToolCallingChatModel(aiConfig.Claude, logger)
	case "gemini":
		return NewGeminiToolCallingChatModel(aiConfig.Gemini, logger)
	case "deepseek":
		cm, err := NewDeepSeekToolCallingChatModel(aiConfig.DeepSeek, logger)
		if err != nil || cm == nil {
			logger.Warn("DeepSeek 模型不可用，回退到支持工具调用的 Mock",
				zap.Error(err),
				zap.String("provider", aiConfig.Provider),
			)
			return chain.NewToolCallingMockChatModel(logger), nil
		}
		return cm, nil
	default:
		logger.Warn("未知的AI模型提供商，使用支持工具调用的模拟模型", zap.String("provider", aiConfig.Provider))
		return chain.NewToolCallingMockChatModel(logger), nil
	}
}

// NewOpenAIToolCallingChatModel 创建OpenAI工具调用聊天模型
func NewOpenAIToolCallingChatModel(config config.OpenAIConfig, logger *zap.Logger) (model.ToolCallingChatModel, error) {
	// 暂时返回错误，等待 eino-ext 更新支持 ToolCallingChatModel
	return nil, errors.New(errors.InternalError, "OpenAI ToolCallingChatModel 暂未实现")
}

// NewClaudeToolCallingChatModel 创建Claude工具调用聊天模型
func NewClaudeToolCallingChatModel(cfg config.ClaudeConfig, logger *zap.Logger) (model.ToolCallingChatModel, error) {
	// 暂时返回错误，等待 eino-ext 更新支持 ToolCallingChatModel
	return nil, errors.New(errors.InternalError, "Claude ToolCallingChatModel 暂未实现")
}

// NewDeepSeekToolCallingChatModel 创建DeepSeek工具调用聊天模型
func NewDeepSeekToolCallingChatModel(cfg config.DeepSeekConfig, logger *zap.Logger) (model.ToolCallingChatModel, error) {
	// 暂时返回错误，等待 eino-ext 更新支持 ToolCallingChatModel
	cm, err := deepseek.NewChatModel(context.TODO(), &deepseek.ChatModelConfig{
		Model:   cfg.Model,
		APIKey:  cfg.APIKey,
		BaseURL: cfg.BaseURL,
	})
	if err != nil {
		return nil, errors.Wrap(errors.InternalError, "创建 DeepSeek 模型失败", err)
	}
	return cm, nil
}

// 创建 Gemini 工具调用聊天模型
func NewGeminiToolCallingChatModel(cfg config.GeminiConfig, logger *zap.Logger) (model.ToolCallingChatModel, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:      cfg.APIKey,
		HTTPOptions: genai.HTTPOptions{BaseURL: cfg.BaseURL},
	})
	if err != nil {
		return nil, errors.Wrap(errors.InternalError, "创建 Gemini 客户端失败", err)
	}
	geminiChatModel, err := gemini.NewChatModel(context.TODO(), &gemini.Config{
		Model:  cfg.Model,
		Client: client,
	})
	if err != nil {
		return nil, errors.Wrap(errors.InternalError, "创建 Gemini 模型失败", err)
	}
	return geminiChatModel, nil
}

//
