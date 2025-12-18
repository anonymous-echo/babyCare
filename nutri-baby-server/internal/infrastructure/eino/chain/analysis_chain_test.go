package chain

import (
	"context"
	"testing"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/eino/tools"
	"go.uber.org/zap"
)

// MockChatModel 模拟ChatModel
type MockChatModel struct {
	mock.Mock
}

func (m *MockChatModel) Generate(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	args := m.Called(ctx, messages, opts)
	return args.Get(0).(*schema.Message), args.Error(1)
}

func (m *MockChatModel) Stream(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	args := m.Called(ctx, messages, opts)
	return args.Get(0).(*schema.StreamReader[*schema.Message]), args.Error(1)
}

func (m *MockChatModel) BindTools(tools []*schema.ToolInfo) error {
	args := m.Called(tools)
	return args.Error(0)
}

// WithTools 模拟 WithTools 方法
func (m *MockChatModel) WithTools(toolInfos []*schema.ToolInfo) (model.ToolCallingChatModel, error) {
	// 在这个简单的 mock 中，我们直接返回自己，因为我们实现了 ToolCallingChatModel 接口
	// 实际使用中可能需要返回一个新的 mock 对象
	return m, nil
}

func TestAnalysisChainBuilder_Analyze(t *testing.T) {
	// 准备测试数据
	mockModel := new(MockChatModel)
	logger := zap.NewNop()

	// 创建一个空的 DataQueryTools，因为我们在测试中不会真正执行工具
	// 注意：这里需要确保 tools 包的构造函数可用，或者 mock 它
	// 由于 tools 包比较复杂，我们这里假设 DataQueryTools 可以被实例化
	// 如果 DataQueryTools 依赖很多，可能需要 mock 它的接口
	// 这里为了简化，我们假设 NewAnalysisChainBuilder 接受 nil tools 也能运行到 Generate 调用

	// 创建 builder
	builder := &AnalysisChainBuilder{
		chatModel:      mockModel,
		dataTools:      &tools.DataQueryTools{}, // 假设这里可以是 nil 或者空结构体
		batchDataTools: &tools.BatchDataTools{},
		logger:         logger,
		enableParallel: true,
	}

	// 模拟 AI 响应，包含 user_friendly 字段
	aiResponseContent := `{
		"score": 85.5,
		"insights": [
			{
				"type": "feeding",
				"title": "喂养规律",
				"description": "宝宝喂养很有规律",
				"priority": "high",
				"category": "positive"
			}
		],
		"alerts": [],
		"patterns": [],
		"predictions": [],
		"user_friendly": {
			"overall_summary": "宝宝表现很棒",
			"score_explanation": "85分是很高的分数",
			"key_highlights": [
				{
					"title": "胃口好",
					"description": "吃得香",
					"icon": "😋"
				}
			],
			"improvement_areas": [],
			"next_step_actions": [],
			"encouraging_words": "继续保持"
		}
	}`

	// 设置 mock 期望
	mockModel.On("Generate", mock.Anything, mock.Anything, mock.Anything).Return(&schema.Message{
		Role:    schema.Assistant,
		Content: aiResponseContent,
	}, nil)

	// 执行测试
	analysis := &entity.AIAnalysis{
		BabyID:       1,
		AnalysisType: entity.AIAnalysisTypeFeeding,
		StartDate:    time.Now().Add(-24 * time.Hour),
		EndDate:      time.Now(),
	}

	result, err := builder.Analyze(context.Background(), analysis)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 85.5, result.Score)
	assert.Len(t, result.Insights, 1)
	assert.Equal(t, "喂养规律", result.Insights[0].Title)

	// 验证 user_friendly 字段是否正确解析
	assert.NotNil(t, result.UserFriendly)
	assert.Equal(t, "宝宝表现很棒", result.UserFriendly.OverallSummary)
	assert.Equal(t, "85分是很高的分数", result.UserFriendly.ScoreExplanation)
	assert.Len(t, result.UserFriendly.KeyHighlights, 1)
	assert.Equal(t, "胃口好", result.UserFriendly.KeyHighlights[0].Title)
}
