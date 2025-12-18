package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/datatypes"
)

// AIAnalysisType AI分析类型
type AIAnalysisType string

const (
	AIAnalysisTypeFeeding  AIAnalysisType = "feeding"  // 喂养分析
	AIAnalysisTypeSleep    AIAnalysisType = "sleep"    // 睡眠分析
	AIAnalysisTypeGrowth   AIAnalysisType = "growth"   // 成长分析
	AIAnalysisTypeHealth   AIAnalysisType = "health"   // 健康分析
	AIAnalysisTypeBehavior AIAnalysisType = "behavior" // 行为分析
)

// AIAnalysisStatus 分析状态
type AIAnalysisStatus string

const (
	AIAnalysisStatusPending   AIAnalysisStatus = "pending"   // 待分析
	AIAnalysisStatusAnalyzing AIAnalysisStatus = "analyzing" // 分析中
	AIAnalysisStatusCompleted AIAnalysisStatus = "completed" // 已完成
	AIAnalysisStatusFailed    AIAnalysisStatus = "failed"    // 分析失败
)

// AIAnalysis AI分析结果实体
type AIAnalysis struct {
	ID           int64            `json:"id" gorm:"primaryKey;autoIncrement"`
	BabyID       int64            `json:"baby_id" gorm:"index;not null"`
	AnalysisType AIAnalysisType   `json:"analysis_type" gorm:"type:varchar(20);not null"`
	Status       AIAnalysisStatus `json:"status" gorm:"type:varchar(20);not null;default:pending"`
	StartDate    time.Time        `json:"start_date" gorm:"not null"`
	EndDate      time.Time        `json:"end_date" gorm:"not null"`
	InputData    string           `json:"input_data" gorm:"type:text"`              // 输入数据JSON
	Result       string           `json:"result" gorm:"type:text"`                  // 分析结果JSON
	Score        *float64         `json:"score,omitempty" gorm:"type:numeric(3,2)"` // 评分(0-100)
	Insights     []string         `json:"insights" gorm:"type:text"`                // 洞察建议
	Alerts       []string         `json:"alerts" gorm:"type:text"`                  // 异常警告
	CreatedAt    time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 表名
func (AIAnalysis) TableName() string {
	return "ai_analyses"
}

// AIAnalysisResult AI分析结果详细结构
type AIAnalysisResult struct {
	AnalysisID   int64                  `json:"analysis_id"`
	BabyID       int64                  `json:"baby_id"`
	AnalysisType AIAnalysisType         `json:"analysis_type"`
	Score        float64                `json:"score"`
	Insights     []AIInsight            `json:"insights"`
	Alerts       []AIAlert              `json:"alerts"`
	Patterns     []AIPattern            `json:"patterns"`
	Predictions  []AIPrediction         `json:"predictions"`
	Metadata     map[string]interface{} `json:"metadata"`
	UserFriendly *UserFriendlyResult    `json:"user_friendly,omitempty"` // 用户友好结果
}

// UserFriendlyResult 用户友好的分析结果
type UserFriendlyResult struct {
	OverallSummary   string                    `json:"overall_summary"`   // 总体评价
	ScoreExplanation string                    `json:"score_explanation"` // 评分说明
	KeyHighlights    []UserFriendlyHighlight   `json:"key_highlights"`    // 关键亮点
	ImprovementAreas []UserFriendlyImprovement `json:"improvement_areas"` // 改进建议
	NextStepActions  []UserFriendlyAction      `json:"next_step_actions"` // 下一步行动
	EncouragingWords string                    `json:"encouraging_words"` // 鼓励话语
}

// UserFriendlyHighlight 用户友好的亮点
type UserFriendlyHighlight struct {
	Title       string `json:"title"`       // 标题
	Description string `json:"description"` // 描述
	Icon        string `json:"icon"`        // 图标建议
}

// UserFriendlyImprovement 用户友好的改进建议
type UserFriendlyImprovement struct {
	Area       string `json:"area"`       // 改进领域
	Issue      string `json:"issue"`      // 问题描述
	Suggestion string `json:"suggestion"` // 具体建议
	Priority   string `json:"priority"`   // 优先级 (high/medium/low)
	Difficulty string `json:"difficulty"` // 实施难度 (easy/medium/hard)
}

// UserFriendlyAction 用户友好的行动建议
type UserFriendlyAction struct {
	Action   string `json:"action"`   // 行动内容
	Timeline string `json:"timeline"` // 时间安排
	Benefit  string `json:"benefit"`  // 预期收益
	HowTo    string `json:"how_to"`   // 具体做法
}

// AIInsight AI洞察建议
type AIInsight struct {
	Type        string `json:"type"`        // insight类型
	Title       string `json:"title"`       // 标题
	Description string `json:"description"` // 详细描述
	Priority    string `json:"priority"`    // 优先级(high/medium/low)
	Category    string `json:"category"`    // 分类
}

// AIAlert AI异常警告
type AIAlert struct {
	Level       string    `json:"level"`       // 警告级别(critical/warning/info)
	Type        string    `json:"type"`        // 警告类型
	Title       string    `json:"title"`       // 标题
	Description string    `json:"description"` // 描述
	Suggestion  string    `json:"suggestion"`  // 建议
	Timestamp   time.Time `json:"timestamp"`   // 时间戳
}

// AIPattern AI识别的模式
type AIPattern struct {
	PatternType string    `json:"pattern_type"` // 模式类型
	Description string    `json:"description"`  // 模式描述
	Confidence  float64   `json:"confidence"`   // 置信度(0-1)
	Frequency   string    `json:"frequency"`    // 出现频率
	TimeRange   TimeRange `json:"time_range"`   // 时间范围
}

// AIPrediction AI预测结果
type AIPrediction struct {
	PredictionType string  `json:"prediction_type"` // 预测类型
	Value          string  `json:"value"`           // 预测值
	Confidence     float64 `json:"confidence"`      // 置信度
	TimeFrame      string  `json:"time_frame"`      // 预测时间范围
	Reason         string  `json:"reason"`          // 预测理由
}

// TimeRange 时间范围
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// DailyTips 每日建议
type DailyTips struct {
	ID        int64                         `json:"id" gorm:"primaryKey;autoIncrement"`
	BabyID    int64                         `json:"baby_id" gorm:"index;not null"`
	Date      time.Time                     `json:"date" gorm:"type:date;index;not null"`
	Tips      datatypes.JSONSlice[DailyTip] `json:"tips" gorm:"type:jsonb"`
	ExpiredAt time.Time                     `json:"expired_at" gorm:"not null"`
	CreatedAt time.Time                     `json:"created_at" gorm:"autoCreateTime"`
}

// DailyTip 单个建议
type DailyTip struct {
	ID          string `json:"id"`
	Icon        string `json:"icon"`                 // 图标emoji
	Title       string `json:"title"`                // 标题
	Description string `json:"description"`          // 描述
	Type        string `json:"type"`                 // 类型(feeding/sleep/health)
	Priority    string `json:"priority"`             // 优先级
	ActionURL   string `json:"action_url,omitempty"` // 跳转链接
}

// TableName 表名
func (DailyTips) TableName() string {
	return "daily_tips"
}

// DailyTipList JSON序列化类型，支持数据库读写
type DailyTipList []DailyTip

// Value 实现 driver.Valuer，将建议列表转换为JSON
func (d DailyTipList) Value() (driver.Value, error) {
	if len(d) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(d)
}

// Scan 实现 sql.Scanner，从数据库JSON解析为结构体
func (d *DailyTipList) Scan(value interface{}) error {
	if value == nil {
		*d = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("DailyTipList: unsupported Scan type %T", value)
	}

	return json.Unmarshal(bytes, d)
}

// AIAnalysisRequest AI分析请求
type AIAnalysisRequest struct {
	BabyID       string                 `json:"baby_id" binding:"required"`
	AnalysisType AIAnalysisType         `json:"analysis_type" binding:"required"`
	StartDate    time.Time              `json:"start_date" binding:"required"`
	EndDate      time.Time              `json:"end_date" binding:"required"`
	Options      map[string]interface{} `json:"options,omitempty"`
}

// AIAnalysisResponse AI分析响应
type AIAnalysisResponse struct {
	AnalysisID string            `json:"analysis_id"`
	Status     AIAnalysisStatus  `json:"status"`
	Result     *AIAnalysisResult `json:"result,omitempty"`
	Message    string            `json:"message,omitempty"`
}
