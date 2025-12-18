package chain

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

// ToolCallingMockChatModel 支持工具调用的模拟聊天模型
type ToolCallingMockChatModel struct {
	logger *zap.Logger
	tools  []*schema.ToolInfo
}

// NewToolCallingMockChatModel 创建支持工具调用的模拟聊天模型
func NewToolCallingMockChatModel(logger *zap.Logger) *ToolCallingMockChatModel {
	return &ToolCallingMockChatModel{
		logger: logger,
		tools:  []*schema.ToolInfo{},
	}
}

// Generate 生成响应
func (m *ToolCallingMockChatModel) Generate(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	m.logger.Debug("ToolCallingMockChatModel.Generate 被调用", zap.Int("message_count", len(messages)))

	// 检查最后一条消息
	lastMessage := messages[len(messages)-1]
	
	// 如果是用户消息且包含分析请求，模拟工具调用
	if lastMessage.Role == schema.User && m.shouldCallTool(lastMessage.Content) {
		return m.generateToolCallResponse(lastMessage.Content), nil
	}
	
	// 如果是工具调用结果，生成最终分析
	if m.hasToolResults(messages) {
		return m.generateFinalAnalysis(messages), nil
	}

	// 默认响应
	return &schema.Message{
		Role:    schema.Assistant,
		Content: m.generateMockResponse(lastMessage.Content),
	}, nil
}

// Stream 流式生成
func (m *ToolCallingMockChatModel) Stream(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	// 简化实现，直接返回完整消息
	message, err := m.Generate(ctx, messages, opts...)
	if err != nil {
		return nil, err
	}

	streamReader, streamWriter := schema.Pipe[*schema.Message](1)
	go func() {
		defer streamWriter.Close()
		streamWriter.Send(message, nil)
	}()

	return streamReader, nil
}

// WithTools 绑定工具
func (m *ToolCallingMockChatModel) WithTools(tools []*schema.ToolInfo) (model.ToolCallingChatModel, error) {
	newModel := &ToolCallingMockChatModel{
		logger: m.logger,
		tools:  tools,
	}
	m.logger.Debug("ToolCallingMockChatModel.WithTools 被调用", zap.Int("tool_count", len(tools)))
	return newModel, nil
}

// shouldCallTool 判断是否应该调用工具
func (m *ToolCallingMockChatModel) shouldCallTool(content string) bool {
	// 如果消息包含分析请求且有可用工具，则调用工具
	return len(m.tools) > 0 && (strings.Contains(content, "分析") || strings.Contains(content, "建议"))
}

// hasToolResults 检查消息历史中是否有工具调用结果
func (m *ToolCallingMockChatModel) hasToolResults(messages []*schema.Message) bool {
	for _, msg := range messages {
		if msg.Role == schema.Tool {
			return true
		}
	}
	return false
}

// generateToolCallResponse 生成工具调用响应
func (m *ToolCallingMockChatModel) generateToolCallResponse(content string) *schema.Message {
	var toolCalls []schema.ToolCall

	// 根据内容决定调用哪些工具
	if strings.Contains(content, "宝宝ID") {
		// 从消息中提取宝宝ID和日期范围
		babyID := m.extractBabyIDFromMessage(content)
		startDate, endDate := m.extractDateRangeFromMessage(content)
		
		// 调用获取宝宝信息工具
		toolCalls = append(toolCalls, schema.ToolCall{
			ID:   "call_baby_info",
			Type: "function",
			Function: schema.FunctionCall{
				Name:      "get_baby_info",
				Arguments: `{"baby_id": ` + babyID + `}`,
			},
		})

		// 根据分析类型调用相应的数据工具
		if strings.Contains(content, "喂养") {
			toolCalls = append(toolCalls, schema.ToolCall{
				ID:   "call_feeding_data",
				Type: "function",
				Function: schema.FunctionCall{
					Name:      "get_feeding_data",
					Arguments: `{"baby_id": ` + babyID + `, "start_date": "` + startDate + `", "end_date": "` + endDate + `", "limit": 100}`,
				},
			})
		}

		if strings.Contains(content, "睡眠") {
			toolCalls = append(toolCalls, schema.ToolCall{
				ID:   "call_sleep_data",
				Type: "function",
				Function: schema.FunctionCall{
					Name:      "get_sleep_data",
					Arguments: `{"baby_id": ` + babyID + `, "start_date": "` + startDate + `", "end_date": "` + endDate + `", "limit": 100}`,
				},
			})
		}

		if strings.Contains(content, "成长") {
			toolCalls = append(toolCalls, schema.ToolCall{
				ID:   "call_growth_data",
				Type: "function",
				Function: schema.FunctionCall{
					Name:      "get_growth_data",
					Arguments: `{"baby_id": ` + babyID + `, "start_date": "` + startDate + `", "end_date": "` + endDate + `", "limit": 100}`,
				},
			})
		}
	}

	return &schema.Message{
		Role:      schema.Assistant,
		Content:   "我需要获取相关数据来进行分析，让我调用一些工具来获取信息。",
		ToolCalls: toolCalls,
	}
}

// generateFinalAnalysis 生成最终分析
func (m *ToolCallingMockChatModel) generateFinalAnalysis(messages []*schema.Message) *schema.Message {
	// 分析工具调用结果
	var hasFeeding, hasSleep, hasGrowth bool
	
	for _, msg := range messages {
		if msg.Role == schema.Tool {
			if strings.Contains(msg.Content, "feeding_data") {
				hasFeeding = true
			}
			if strings.Contains(msg.Content, "sleep_data") {
				hasSleep = true
			}
			if strings.Contains(msg.Content, "growth_data") {
				hasGrowth = true
			}
		}
	}

	// 根据获取到的数据类型生成相应的分析结果
	var analysisResult string
	
	if hasFeeding {
		analysisResult = `{
			"score": 85,
			"insights": [
				{
					"type": "feeding",
					"title": "喂养规律良好",
					"description": "基于获取的喂养数据分析，宝宝的喂养时间较为规律，建议继续保持",
					"priority": "medium",
					"category": "规律性"
				}
			],
			"alerts": [],
			"patterns": [
				{
					"pattern_type": "regular_feeding",
					"description": "每3-4小时喂养一次",
					"confidence": 0.9,
					"frequency": "daily"
				}
			],
			"predictions": []
		}`
	} else if hasSleep {
		analysisResult = `{
			"score": 78,
			"insights": [
				{
					"type": "sleep",
					"title": "睡眠质量良好",
					"description": "基于获取的睡眠数据分析，宝宝睡眠时长符合月龄标准",
					"priority": "high",
					"category": "睡眠质量"
				}
			],
			"alerts": [],
			"patterns": [],
			"predictions": []
		}`
	} else if hasGrowth {
		analysisResult = `{
			"score": 92,
			"insights": [
				{
					"type": "growth",
					"title": "生长发育正常",
					"description": "基于获取的成长数据分析，身高体重增长曲线正常",
					"priority": "high",
					"category": "发育评估"
				}
			],
			"alerts": [],
			"patterns": [],
			"predictions": []
		}`
	} else {
		// 生成每日建议
		analysisResult = `[
			{
				"id": "tip_1",
				"icon": "🍼",
				"title": "基于数据的喂养建议",
				"description": "根据获取的数据分析，建议在上午9-10点之间进行喂养",
				"type": "feeding",
				"priority": "high",
				"action_url": "/pages/record/feeding/index"
			},
			{
				"id": "tip_2",
				"icon": "😴",
				"title": "睡眠时间优化",
				"description": "基于睡眠数据分析，建议调整午睡时间",
				"type": "sleep",
				"priority": "medium",
				"action_url": "/pages/record/sleep/index"
			}
		]`
	}

	return &schema.Message{
		Role:    schema.Assistant,
		Content: analysisResult,
	}
}

// generateMockResponse 生成模拟响应（兜底）
func (m *ToolCallingMockChatModel) generateMockResponse(userInput string) string {
	return `{"score":80,"insights":[{"type":"general","title":"整体状况良好","description":"宝宝各项指标基本正常","priority":"medium","category":"综合评估"}],"alerts":[],"patterns":[],"predictions":[]}`
}

// extractBabyIDFromMessage 从消息中提取宝宝ID
func (m *ToolCallingMockChatModel) extractBabyIDFromMessage(content string) string {
	// 使用正则表达式提取宝宝ID
	re := regexp.MustCompile(`宝宝ID\s*(\d+)`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	
	// 如果没有找到，返回默认值
	m.logger.Warn("无法从消息中提取宝宝ID，使用默认值1", zap.String("content", content))
	return "1"
}

// extractDateRangeFromMessage 从消息中提取日期范围
func (m *ToolCallingMockChatModel) extractDateRangeFromMessage(content string) (startDate, endDate string) {
	// 提取日期范围
	dateRe := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})\s*至\s*(\d{4}-\d{2}-\d{2})`)
	matches := dateRe.FindStringSubmatch(content)
	if len(matches) > 2 {
		return matches[1], matches[2]
	}
	
	// 默认返回最近7天
	endDate = time.Now().Format("2006-01-02")
	startDate = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	return startDate, endDate
}

// 确保实现了 ToolCallingChatModel 接口
var _ model.ToolCallingChatModel = (*ToolCallingMockChatModel)(nil)
