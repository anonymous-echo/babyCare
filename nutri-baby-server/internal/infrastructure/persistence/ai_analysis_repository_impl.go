package persistence

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// aiAnalysisRepositoryImpl AI分析结果仓储实现
type aiAnalysisRepositoryImpl struct {
	db *gorm.DB
}

// NewAIAnalysisRepository 创建AI分析结果仓储实例
func NewAIAnalysisRepository(db *gorm.DB) repository.AIAnalysisRepository {
	return &aiAnalysisRepositoryImpl{db: db}
}

// Create 创建分析记录
func (r *aiAnalysisRepositoryImpl) Create(ctx context.Context, analysis *entity.AIAnalysis) error {
	if err := r.db.WithContext(ctx).Create(analysis).Error; err != nil {
		return errors.Wrap(errors.DatabaseError, "创建AI分析记录失败", err)
	}
	return nil
}

// GetByID 根据ID获取分析结果
func (r *aiAnalysisRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.AIAnalysis, error) {
	var analysis entity.AIAnalysis
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&analysis).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errors.NotFound, "AI分析记录不存在")
		}
		return nil, errors.Wrap(errors.DatabaseError, "查询AI分析记录失败", err)
	}
	return &analysis, nil
}

// GetLatestByBabyIDAndType 根据宝宝ID和分析类型获取最新的分析结果
func (r *aiAnalysisRepositoryImpl) GetLatestByBabyIDAndType(ctx context.Context, babyID int64, analysisType entity.AIAnalysisType) (*entity.AIAnalysis, error) {
	var analysis entity.AIAnalysis
	err := r.db.WithContext(ctx).
		Where("baby_id = ? AND analysis_type = ? AND status = ?", babyID, analysisType, entity.AIAnalysisStatusCompleted).
		Order("created_at DESC").
		First(&analysis).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errors.NotFound, "未找到对应的AI分析记录")
		}
		return nil, errors.Wrap(errors.DatabaseError, "查询最新AI分析记录失败", err)
	}
	return &analysis, nil
}

// GetByDateRange 获取指定时间范围内的分析结果
func (r *aiAnalysisRepositoryImpl) GetByDateRange(ctx context.Context, babyID int64, analysisType entity.AIAnalysisType, startDate, endDate time.Time) ([]*entity.AIAnalysis, error) {
	var analyses []*entity.AIAnalysis
	query := r.db.WithContext(ctx).
		Where("baby_id = ? AND created_at BETWEEN ? AND ?", babyID, startDate, endDate)

	if analysisType != "" {
		query = query.Where("analysis_type = ?", analysisType)
	}

	err := query.Order("created_at DESC").Find(&analyses).Error
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "查询AI分析记录失败", err)
	}
	return analyses, nil
}

// GetPendingAnalyses 获取待分析或分析中的记录
func (r *aiAnalysisRepositoryImpl) GetPendingAnalyses(ctx context.Context, limit int) ([]*entity.AIAnalysis, error) {
	var analyses []*entity.AIAnalysis
	err := r.db.WithContext(ctx).
		Where("status IN (?, ?)", entity.AIAnalysisStatusPending, entity.AIAnalysisStatusAnalyzing).
		Order("created_at ASC").
		Limit(limit).
		Find(&analyses).Error
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "查询待分析记录失败", err)
	}
	return analyses, nil
}

// UpdateResult 更新分析结果
func (r *aiAnalysisRepositoryImpl) UpdateResult(ctx context.Context, id int64, result string, status entity.AIAnalysisStatus) error {
	updates := map[string]interface{}{
		"result": result,
		"status": status,
	}

	err := r.db.WithContext(ctx).Model(&entity.AIAnalysis{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return errors.Wrap(errors.DatabaseError, "更新AI分析结果失败", err)
	}
	return nil
}

// UpdateStatus 更新分析状态
func (r *aiAnalysisRepositoryImpl) UpdateStatus(ctx context.Context, id int64, status entity.AIAnalysisStatus) error {
	err := r.db.WithContext(ctx).Model(&entity.AIAnalysis{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		return errors.Wrap(errors.DatabaseError, "更新AI分析状态失败", err)
	}
	return nil
}

// Delete 删除分析记录
func (r *aiAnalysisRepositoryImpl) Delete(ctx context.Context, id int64) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.AIAnalysis{}).Error
	if err != nil {
		return errors.Wrap(errors.DatabaseError, "删除AI分析记录失败", err)
	}
	return nil
}

// GetAnalysisStats 获取宝宝的分析历史统计
func (r *aiAnalysisRepositoryImpl) GetAnalysisStats(ctx context.Context, babyID int64) (*repository.AnalysisStats, error) {
	var stats repository.AnalysisStats

	// 总分析数
	var totalCount int64
	if err := r.db.WithContext(ctx).Model(&entity.AIAnalysis{}).Where("baby_id = ?", babyID).Count(&totalCount).Error; err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "统计总分析数失败", err)
	}
	stats.TotalAnalyses = totalCount

	// 完成数
	var completedCount int64
	if err := r.db.WithContext(ctx).Model(&entity.AIAnalysis{}).
		Where("baby_id = ? AND status = ?", babyID, entity.AIAnalysisStatusCompleted).
		Count(&completedCount).Error; err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "统计完成分析数失败", err)
	}
	stats.CompletedAnalyses = completedCount

	// 失败数
	stats.FailedAnalyses = totalCount - completedCount

	// 平均分
	var avgScore *float64
	if err := r.db.WithContext(ctx).
		Model(&entity.AIAnalysis{}).
		Where("baby_id = ? AND status = ? AND score IS NOT NULL", babyID, entity.AIAnalysisStatusCompleted).
		Select("AVG(score)").
		Scan(&avgScore).Error; err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "统计平均分失败", err)
	}
	stats.AverageScore = avgScore

	// 各类型统计
	typeCounts := make(map[string]int64)
	var typeStats []struct {
		AnalysisType string
		Count        int64
	}

	if err := r.db.WithContext(ctx).
		Model(&entity.AIAnalysis{}).
		Where("baby_id = ?", babyID).
		Select("analysis_type, COUNT(*) as count").
		Group("analysis_type").
		Scan(&typeStats).Error; err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "统计各类型分析数失败", err)
	}

	for _, stat := range typeStats {
		typeCounts[stat.AnalysisType] = stat.Count
	}
	stats.AnalysisTypeCounts = typeCounts

	// 最近分析记录
	var recentAnalyses []*entity.AIAnalysis
	if err := r.db.WithContext(ctx).
		Where("baby_id = ?", babyID).
		Order("created_at DESC").
		Limit(5).
		Find(&recentAnalyses).Error; err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "查询最近分析记录失败", err)
	}
	stats.RecentAnalyses = recentAnalyses

	return &stats, nil
}

// dailyTipsRepositoryImpl 每日建议仓储实现
type dailyTipsRepositoryImpl struct {
	db *gorm.DB
}

// NewDailyTipsRepository 创建每日建议仓储实例
func NewDailyTipsRepository(db *gorm.DB) repository.DailyTipsRepository {
	return &dailyTipsRepositoryImpl{db: db}
}

// Create 创建每日建议
func (r *dailyTipsRepositoryImpl) Create(ctx context.Context, tips *entity.DailyTips) error {
	if err := r.db.WithContext(ctx).Create(tips).Error; err != nil {
		return errors.Wrap(errors.DatabaseError, "创建每日建议失败", err)
	}
	return nil
}

// GetByBabyIDAndDate 根据宝宝ID和日期获取每日建议
func (r *dailyTipsRepositoryImpl) GetByBabyIDAndDate(ctx context.Context, babyID int64, date time.Time) (*entity.DailyTips, error) {
	var tips entity.DailyTips
	err := r.db.WithContext(ctx).
		Where("baby_id = ? AND date = ?", babyID, date.Format("2006-01-02")).
		First(&tips).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errors.NotFound, "未找到对应的每日建议")
		}
		return nil, errors.Wrap(errors.DatabaseError, "查询每日建议失败", err)
	}
	return &tips, nil
}

// GetLatestByBabyID 获取宝宝的最新每日建议
func (r *dailyTipsRepositoryImpl) GetLatestByBabyID(ctx context.Context, babyID int64) (*entity.DailyTips, error) {
	var tips entity.DailyTips
	err := r.db.WithContext(ctx).
		Where("baby_id = ? AND expired_at > ?", babyID, time.Now()).
		Order("date DESC").
		First(&tips).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errors.NotFound, "未找到有效的每日建议")
		}
		return nil, errors.Wrap(errors.DatabaseError, "查询最新每日建议失败", err)
	}
	return &tips, nil
}

// GetByDateRange 获取指定日期范围内的每日建议
func (r *dailyTipsRepositoryImpl) GetByDateRange(ctx context.Context, babyID int64, startDate, endDate time.Time) ([]*entity.DailyTips, error) {
	var tips []*entity.DailyTips
	err := r.db.WithContext(ctx).
		Where("baby_id = ? AND date BETWEEN ? AND ?", babyID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("date DESC").
		Find(&tips).Error
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "查询每日建议失败", err)
	}
	return tips, nil
}

// Update 更新每日建议
func (r *dailyTipsRepositoryImpl) Update(ctx context.Context, tips *entity.DailyTips) error {
	tipsJSON, err := json.Marshal(tips.Tips)
	if err != nil {
		return errors.Wrap(errors.ParamError, "序列化建议数据失败", err)
	}

	updates := map[string]interface{}{
		"tips": tipsJSON,
	}

	err = r.db.WithContext(ctx).Model(&entity.DailyTips{}).Where("id = ?", tips.ID).Updates(updates).Error
	if err != nil {
		return errors.Wrap(errors.DatabaseError, "更新每日建议失败", err)
	}
	return nil
}

// DeleteExpired 删除过期的每日建议
func (r *dailyTipsRepositoryImpl) DeleteExpired(ctx context.Context, before time.Time) error {
	err := r.db.WithContext(ctx).
		Where("expired_at < ?", before).
		Delete(&entity.DailyTips{}).Error
	if err != nil {
		return errors.Wrap(errors.DatabaseError, "删除过期每日建议失败", err)
	}
	return nil
}