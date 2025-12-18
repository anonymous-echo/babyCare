package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/eino/chain"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
	"go.uber.org/zap"
)

// 请求和响应类型定义

// CreateAnalysisRequest 创建分析请求
type CreateAnalysisRequest struct {
	BabyID       int64                 `json:"baby_id" binding:"required"`
	AnalysisType entity.AIAnalysisType `json:"analysis_type" binding:"required"`
	StartDate    CustomTime            `json:"start_date" binding:"required"`
	EndDate      CustomTime            `json:"end_date" binding:"required"`
}

// AnalysisResponse 分析响应
type AnalysisResponse struct {
	AnalysisID int64                    `json:"analysis_id"` // 修改为int64以匹配前端number类型
	Status     entity.AIAnalysisStatus  `json:"status"`
	Result     *entity.AIAnalysisResult `json:"result,omitempty"`
	CreatedAt  time.Time                `json:"created_at"`
}

// AnalysisStatusResponse 分析状态响应（用于轮询）
type AnalysisStatusResponse struct {
	AnalysisID string                  `json:"analysis_id"`
	Status     entity.AIAnalysisStatus `json:"status"`
	Progress   int                     `json:"progress"` // 进度百分比 0-100
	Message    string                  `json:"message"`  // 状态描述
	UpdatedAt  time.Time               `json:"updated_at"`
}

// DailyTipsResponse 每日建议响应
type DailyTipsResponse struct {
	Tips        []entity.DailyTip `json:"tips"`
	GeneratedAt time.Time         `json:"generated_at"`
	ExpiredAt   time.Time         `json:"expired_at"`
}

// BatchAnalysisRequest 批量分析请求
type BatchAnalysisRequest struct {
	BabyID        int64                   `json:"baby_id" binding:"required"`
	AnalysisTypes []entity.AIAnalysisType `json:"analysis_types" binding:"required"`
	StartDate     CustomTime              `json:"start_date" binding:"required"`
	EndDate       CustomTime              `json:"end_date" binding:"required"`
}

// BatchAnalysisResponse 批量分析响应
type BatchAnalysisResponse struct {
	TotalCount     int                `json:"total_count"`
	Analyses       []AnalysisResponse `json:"analyses"`
	CompletedCount int                `json:"completed_count"` // 新增：已完成数量
	FailedCount    int                `json:"failed_count"`    // 新增：失败数量
}

// AnalysisStatsResponse 分析统计响应
type AnalysisStatsResponse struct {
	TotalAnalyses int                             `json:"total_analyses"`
	ByType        map[entity.AIAnalysisType]int   `json:"by_type"`
	ByStatus      map[entity.AIAnalysisStatus]int `json:"by_status"`
	AvgScore      float64                         `json:"avg_score"`
	Trends        []AnalysisTrend                 `json:"trends"`
}

// AnalysisTrend 分析趋势
type AnalysisTrend struct {
	Date  time.Time             `json:"date"`
	Score float64               `json:"score"`
	Type  entity.AIAnalysisType `json:"type"`
}

// CustomTime 自定义时间类型
type CustomTime struct {
	time.Time
}

// UnmarshalJSON 自定义JSON反序列化
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1] // 移除引号
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

// MarshalJSON 自定义JSON序列化
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.Time.Format("2006-01-02"))
}

// AIAnalysisService AI分析服务接口
type AIAnalysisService interface {
	// 创建分析任务
	CreateAnalysis(ctx context.Context, req *CreateAnalysisRequest) (*AnalysisResponse, error)

	// 生成每日建议
	GenerateDailyTips(ctx context.Context, babyID string, date time.Time) (*DailyTipsResponse, error)

	// 处理待分析的任务
	ProcessPendingAnalyses(ctx context.Context) error

	// 获取分析结果
	GetAnalysisResult(ctx context.Context, analysisID string) (*AnalysisResponse, error)

	// 获取分析状态（用于轮询）
	GetAnalysisStatus(ctx context.Context, analysisID string) (*AnalysisStatusResponse, error)

	// 获取最新分析
	GetLatestAnalysis(ctx context.Context, babyID string, analysisType entity.AIAnalysisType) (*AnalysisResponse, error)

	// 批量分析
	BatchAnalyze(ctx context.Context, req *BatchAnalysisRequest) (*BatchAnalysisResponse, error)

	// 获取每日建议
	GetDailyTips(ctx context.Context, babyID string, date time.Time) (*DailyTipsResponse, error)

	// 获取分析统计
	GetAnalysisStats(ctx context.Context, babyID string, days int) (*AnalysisStatsResponse, error)
}

// aiAnalysisServiceImpl AI分析服务实现
type aiAnalysisServiceImpl struct {
	aiAnalysisRepo repository.AIAnalysisRepository
	dailyTipsRepo  repository.DailyTipsRepository
	babyRepo       repository.BabyRepository
	chainBuilder   *chain.AnalysisChainBuilder
	logger         *zap.Logger
}

// NewAIAnalysisService 创建AI分析服务实例
func NewAIAnalysisService(
	aiAnalysisRepo repository.AIAnalysisRepository,
	dailyTipsRepo repository.DailyTipsRepository,
	babyRepo repository.BabyRepository,
	chainBuilder *chain.AnalysisChainBuilder,
	logger *zap.Logger,
) AIAnalysisService {
	return &aiAnalysisServiceImpl{
		aiAnalysisRepo: aiAnalysisRepo,
		dailyTipsRepo:  dailyTipsRepo,
		babyRepo:       babyRepo,
		chainBuilder:   chainBuilder,
		logger:         logger,
	}
}

// CreateAnalysis 创建分析任务（异步模式）
func (s *aiAnalysisServiceImpl) CreateAnalysis(ctx context.Context, req *CreateAnalysisRequest) (*AnalysisResponse, error) {
	// 验证宝宝是否存在
	_, err := s.babyRepo.FindByID(ctx, req.BabyID)
	if err != nil {
		return nil, errors.Wrap(errors.NotFound, "获取宝宝信息失败", err)
	}

	// 验证时间范围
	if req.EndDate.Before(req.StartDate.Time) {
		return nil, errors.New(errors.ParamError, "结束日期不能早于开始日期")
	}

	// 创建分析记录
	analysis := &entity.AIAnalysis{
		BabyID:       req.BabyID,
		AnalysisType: req.AnalysisType,
		Status:       entity.AIAnalysisStatusPending,
		StartDate:    req.StartDate.Time,
		EndDate:      req.EndDate.Time,
	}

	if err := s.aiAnalysisRepo.Create(ctx, analysis); err != nil {
		return nil, errors.Wrap(errors.InternalError, "创建分析任务失败", err)
	}

	// 异步处理分析任务
	go func() {
		// 使用新的context，避免原context被取消
		bgCtx := context.Background()

		if err := s.processAnalysis(bgCtx, analysis.ID); err != nil {
			s.logger.Error("AI分析失败",
				zap.Int64("analysis_id", analysis.ID),
				zap.Error(err),
			)
			// 更新状态为失败
			s.aiAnalysisRepo.UpdateStatus(bgCtx, analysis.ID, entity.AIAnalysisStatusFailed)
		}
	}()

	// 立即返回任务ID和pending状态
	return &AnalysisResponse{
		AnalysisID: analysis.ID, // 直接使用int64类型
		Status:     entity.AIAnalysisStatusPending,
		CreatedAt:  analysis.CreatedAt,
	}, nil
}

// GenerateDailyTips 生成每日建议
func (s *aiAnalysisServiceImpl) GenerateDailyTips(ctx context.Context, babyID string, date time.Time) (*DailyTipsResponse, error) {
	// 转换ID类型
	id, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return nil, errors.Wrap(errors.ParamError, "无效的宝宝ID", err)
	}

	// 优先检查是否已存在当日建议，如果存在直接返回
	existingTips, err := s.dailyTipsRepo.GetByBabyIDAndDate(ctx, id, date)
	if err == nil && existingTips != nil {
		s.logger.Info("返回已存在的每日建议",
			zap.String("baby_id", babyID),
			zap.String("date", date.Format("2006-01-02")),
		)
		return &DailyTipsResponse{
			Tips:        existingTips.Tips,
			GeneratedAt: existingTips.CreatedAt,
			ExpiredAt:   existingTips.ExpiredAt,
		}, nil
	}

	// 记录查询错误但继续生成（可能是记录不存在）
	if err != nil {
		s.logger.Debug("查询已存在建议时出错，将生成新建议",
			zap.String("baby_id", babyID),
			zap.Error(err),
		)
	}

	// 获取宝宝信息
	baby, err := s.babyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(errors.NotFound, "获取宝宝信息失败", err)
	}

	s.logger.Info("开始生成新的每日建议",
		zap.String("baby_id", babyID),
		zap.String("date", date.Format("2006-01-02")),
	)

	// 使用分析链生成建议
	tips, err := s.chainBuilder.GenerateDailyTips(ctx, baby, date)
	if err != nil {
		s.logger.Error("生成每日建议失败",
			zap.String("baby_id", babyID),
			zap.Error(err),
		)
		return nil, errors.Wrap(errors.InternalError, "生成每日建议失败", err)
	}

	// 保存建议
	dailyTips := &entity.DailyTips{
		BabyID:    id,
		Date:      date,
		Tips:      tips,
		ExpiredAt: date.AddDate(0, 0, 1), // 24小时后过期
	}

	if err := s.dailyTipsRepo.Create(ctx, dailyTips); err != nil {
		return nil, errors.Wrap(errors.InternalError, "保存每日建议失败", err)
	}

	s.logger.Info("成功生成并保存每日建议",
		zap.String("baby_id", babyID),
		zap.Int("tips_count", len(tips)),
	)

	return &DailyTipsResponse{
		Tips:        tips,
		GeneratedAt: dailyTips.CreatedAt,
		ExpiredAt:   dailyTips.ExpiredAt,
	}, nil
}

// ProcessPendingAnalyses 处理待分析的任务
func (s *aiAnalysisServiceImpl) ProcessPendingAnalyses(ctx context.Context) error {
	pendingAnalyses, err := s.aiAnalysisRepo.GetPendingAnalyses(ctx, 10)
	if err != nil {
		return errors.Wrap(errors.InternalError, "获取待分析任务失败", err)
	}

	for _, analysis := range pendingAnalyses {
		if err := s.processAnalysis(ctx, analysis.ID); err != nil {
			s.logger.Error("处理分析任务失败",
				zap.Int64("analysis_id", analysis.ID),
				zap.Error(err),
			)
			// 更新状态为失败
			s.aiAnalysisRepo.UpdateStatus(ctx, analysis.ID, entity.AIAnalysisStatusFailed)
			continue
		}
	}

	return nil
}

// GetAnalysisResult 获取分析结果
func (s *aiAnalysisServiceImpl) GetAnalysisResult(ctx context.Context, analysisID string) (*AnalysisResponse, error) {
	id, err := strconv.ParseInt(analysisID, 10, 64)
	if err != nil {
		return nil, errors.Wrap(errors.ParamError, "无效的分析ID", err)
	}

	analysis, err := s.aiAnalysisRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(errors.NotFound, "获取分析记录失败", err)
	}

	response := &AnalysisResponse{
		AnalysisID: id, // 使用int64类型的id
		Status:     analysis.Status,
		CreatedAt:  analysis.CreatedAt,
	}

	// 如果分析已完成，解析结果
	if analysis.Status == entity.AIAnalysisStatusCompleted && analysis.Result != "" {
		var result entity.AIAnalysisResult
		if err := json.Unmarshal([]byte(analysis.Result), &result); err != nil {
			s.logger.Error("解析分析结果失败",
				zap.Int64("analysis_id", id),
				zap.Error(err),
			)
			return response, nil
		}
		response.Result = &result
	}

	return response, nil
}

// GetAnalysisStatus 获取分析状态（用于轮询）
func (s *aiAnalysisServiceImpl) GetAnalysisStatus(ctx context.Context, analysisID string) (*AnalysisStatusResponse, error) {
	id, err := strconv.ParseInt(analysisID, 10, 64)
	if err != nil {
		return nil, errors.Wrap(errors.ParamError, "无效的分析ID", err)
	}

	analysis, err := s.aiAnalysisRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(errors.NotFound, "获取分析记录失败", err)
	}

	// 根据状态计算进度和消息
	var progress int
	var message string

	switch analysis.Status {
	case entity.AIAnalysisStatusPending:
		progress = 10
		message = "分析任务已创建，等待处理"
	case entity.AIAnalysisStatusAnalyzing:
		progress = 50
		message = "AI正在分析数据中..."
	case entity.AIAnalysisStatusCompleted:
		progress = 100
		message = "分析完成"
	case entity.AIAnalysisStatusFailed:
		progress = 0
		message = "分析失败，请重试"
	default:
		progress = 0
		message = "未知状态"
	}

	return &AnalysisStatusResponse{
		AnalysisID: analysisID,
		Status:     analysis.Status,
		Progress:   progress,
		Message:    message,
		UpdatedAt:  analysis.UpdatedAt,
	}, nil
}

// GetLatestAnalysis 获取最新分析
func (s *aiAnalysisServiceImpl) GetLatestAnalysis(ctx context.Context, babyID string, analysisType entity.AIAnalysisType) (*AnalysisResponse, error) {
	id, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return nil, errors.Wrap(errors.ParamError, "无效的宝宝ID", err)
	}

	analysis, err := s.aiAnalysisRepo.GetLatestByBabyIDAndType(ctx, id, analysisType)
	if err != nil {
		return nil, errors.Wrap(errors.NotFound, "获取最新分析失败", err)
	}

	return s.GetAnalysisResult(ctx, strconv.FormatInt(analysis.ID, 10))
}

// BatchAnalyze 批量分析（异步模式）
func (s *aiAnalysisServiceImpl) BatchAnalyze(ctx context.Context, req *BatchAnalysisRequest) (*BatchAnalysisResponse, error) {
	// 验证宝宝是否存在
	_, err := s.babyRepo.FindByID(ctx, req.BabyID)
	if err != nil {
		return nil, errors.Wrap(errors.NotFound, "获取宝宝信息失败", err)
	}

	// 验证时间范围
	if req.EndDate.Before(req.StartDate.Time) {
		return nil, errors.New(errors.ParamError, "结束日期不能早于开始日期")
	}

	var results []AnalysisResponse

	// 如果没有指定分析类型，使用所有类型
	analysisTypes := req.AnalysisTypes
	if len(analysisTypes) == 0 {
		analysisTypes = []entity.AIAnalysisType{
			entity.AIAnalysisTypeFeeding,
			entity.AIAnalysisTypeSleep,
			entity.AIAnalysisTypeGrowth,
			entity.AIAnalysisTypeHealth,
			entity.AIAnalysisTypeBehavior,
		}
	}

	for _, analysisType := range analysisTypes {
		// 创建分析记录
		analysis := &entity.AIAnalysis{
			BabyID:       req.BabyID,
			AnalysisType: analysisType,
			Status:       entity.AIAnalysisStatusPending,
			StartDate:    req.StartDate.Time,
			EndDate:      req.EndDate.Time,
		}

		if err := s.aiAnalysisRepo.Create(ctx, analysis); err != nil {
			s.logger.Error("创建批量分析任务失败",
				zap.Int64("baby_id", req.BabyID),
				zap.String("analysis_type", string(analysisType)),
				zap.Error(err),
			)
			continue
		}

		// 异步处理分析任务
		analysisID := analysis.ID
		go func(id int64, aType entity.AIAnalysisType) {
			bgCtx := context.Background()

			if err := s.processAnalysis(bgCtx, id); err != nil {
				s.logger.Error("批量AI分析失败",
					zap.Int64("analysis_id", id),
					zap.String("analysis_type", string(aType)),
					zap.Error(err),
				)
				s.aiAnalysisRepo.UpdateStatus(bgCtx, id, entity.AIAnalysisStatusFailed)
			}
		}(analysisID, analysisType)

		// 添加到结果列表
		results = append(results, AnalysisResponse{
			AnalysisID: analysis.ID, // 使用int64类型
			Status:     entity.AIAnalysisStatusPending,
			CreatedAt:  analysis.CreatedAt,
		})
	}

	return &BatchAnalysisResponse{
		TotalCount:     len(results),
		Analyses:       results,
		CompletedCount: 0, // 初始时都是pending状态
		FailedCount:    0,
	}, nil
}

// GetDailyTips 获取每日建议
func (s *aiAnalysisServiceImpl) GetDailyTips(ctx context.Context, babyID string, date time.Time) (*DailyTipsResponse, error) {
	return s.GenerateDailyTips(ctx, babyID, date)
}

// GetAnalysisStats 获取分析统计
func (s *aiAnalysisServiceImpl) GetAnalysisStats(ctx context.Context, babyID string, days int) (*AnalysisStatsResponse, error) {
	// 简化实现，返回基本统计信息
	stats := &AnalysisStatsResponse{
		TotalAnalyses: 0,
		ByType:        make(map[entity.AIAnalysisType]int),
		ByStatus:      make(map[entity.AIAnalysisStatus]int),
		AvgScore:      0.0,
		Trends:        []AnalysisTrend{},
	}

	// TODO: 实现完整的统计功能，需要扩展 repository 接口
	s.logger.Info("GetAnalysisStats 调用",
		zap.String("baby_id", babyID),
		zap.Int("days", days),
	)

	return stats, nil
}

// processAnalysis 处理单个分析任务
func (s *aiAnalysisServiceImpl) processAnalysis(ctx context.Context, analysisID int64) error {
	// 更新状态为分析中
	if err := s.aiAnalysisRepo.UpdateStatus(ctx, analysisID, entity.AIAnalysisStatusAnalyzing); err != nil {
		return errors.Wrap(errors.InternalError, "更新分析状态失败", err)
	}

	// 获取分析记录
	analysis, err := s.aiAnalysisRepo.GetByID(ctx, analysisID)
	if err != nil {
		return errors.Wrap(errors.InternalError, "获取分析记录失败", err)
	}

	// 使用分析链进行分析
	result, err := s.chainBuilder.Analyze(ctx, analysis)
	if err != nil {
		return errors.Wrap(errors.InternalError, "AI分析失败", err)
	}

	// 序列化结果
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return errors.Wrap(errors.ParamError, "序列化分析结果失败", err)
	}

	// 更新分析结果
	if err := s.aiAnalysisRepo.UpdateResult(ctx, analysisID, string(resultJSON), entity.AIAnalysisStatusCompleted); err != nil {
		return errors.Wrap(errors.InternalError, "更新分析结果失败", err)
	}

	s.logger.Info("AI分析任务完成",
		zap.Int64("analysis_id", analysisID),
		zap.String("analysis_type", string(analysis.AnalysisType)),
		zap.Float64("score", result.Score),
	)

	return nil
}
