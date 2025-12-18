package repository

import (
	"context"
	"time"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
)

// AIAnalysisRepository AI分析结果仓储接口
type AIAnalysisRepository interface {
	// 创建分析记录
	Create(ctx context.Context, analysis *entity.AIAnalysis) error

	// 根据ID获取分析结果
	GetByID(ctx context.Context, id int64) (*entity.AIAnalysis, error)

	// 根据宝宝ID和分析类型获取最新的分析结果
	GetLatestByBabyIDAndType(ctx context.Context, babyID int64, analysisType entity.AIAnalysisType) (*entity.AIAnalysis, error)

	// 获取指定时间范围内的分析结果
	GetByDateRange(ctx context.Context, babyID int64, analysisType entity.AIAnalysisType, startDate, endDate time.Time) ([]*entity.AIAnalysis, error)

	// 获取待分析或分析中的记录
	GetPendingAnalyses(ctx context.Context, limit int) ([]*entity.AIAnalysis, error)

	// 更新分析结果
	UpdateResult(ctx context.Context, id int64, result string, status entity.AIAnalysisStatus) error

	// 更新分析状态
	UpdateStatus(ctx context.Context, id int64, status entity.AIAnalysisStatus) error

	// 删除分析记录
	Delete(ctx context.Context, id int64) error

	// 获取宝宝的分析历史统计
	GetAnalysisStats(ctx context.Context, babyID int64) (*AnalysisStats, error)
}

// DailyTipsRepository 每日建议仓储接口
type DailyTipsRepository interface {
	// 创建每日建议
	Create(ctx context.Context, tips *entity.DailyTips) error

	// 根据宝宝ID和日期获取每日建议
	GetByBabyIDAndDate(ctx context.Context, babyID int64, date time.Time) (*entity.DailyTips, error)

	// 获取宝宝的最新每日建议
	GetLatestByBabyID(ctx context.Context, babyID int64) (*entity.DailyTips, error)

	// 获取指定日期范围内的每日建议
	GetByDateRange(ctx context.Context, babyID int64, startDate, endDate time.Time) ([]*entity.DailyTips, error)

	// 更新每日建议
	Update(ctx context.Context, tips *entity.DailyTips) error

	// 删除过期的每日建议
	DeleteExpired(ctx context.Context, before time.Time) error
}

// AnalysisStats 分析统计信息
type AnalysisStats struct {
	TotalAnalyses      int64                    `json:"total_analyses"`
	CompletedAnalyses  int64                    `json:"completed_analyses"`
	FailedAnalyses     int64                    `json:"failed_analyses"`
	AverageScore       *float64                 `json:"average_score"`
	AnalysisTypeCounts map[string]int64         `json:"analysis_type_counts"`
	RecentAnalyses     []*entity.AIAnalysis     `json:"recent_analyses"`
}

// AIAnalysisParams 分析查询参数
type AIAnalysisParams struct {
	BabyID       int64
	AnalysisType *entity.AIAnalysisType
	StartDate    *time.Time
	EndDate      *time.Time
	Status       *entity.AIAnalysisStatus
	Limit        int
	Offset       int
}