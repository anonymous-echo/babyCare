package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
)

// BatchDataTools 批量数据查询工具
type BatchDataTools struct {
	babyRepo    repository.BabyRepository
	feedingRepo repository.FeedingRecordRepository
	sleepRepo   repository.SleepRecordRepository
	growthRepo  repository.GrowthRecordRepository
	diaperRepo  repository.DiaperRecordRepository
	logger      *zap.Logger
}

// NewBatchDataTools 创建批量数据查询工具
func NewBatchDataTools(
	babyRepo repository.BabyRepository,
	feedingRepo repository.FeedingRecordRepository,
	sleepRepo repository.SleepRecordRepository,
	growthRepo repository.GrowthRecordRepository,
	diaperRepo repository.DiaperRecordRepository,
	logger *zap.Logger,
) *BatchDataTools {
	return &BatchDataTools{
		babyRepo:    babyRepo,
		feedingRepo: feedingRepo,
		sleepRepo:   sleepRepo,
		growthRepo:  growthRepo,
		diaperRepo:  diaperRepo,
		logger:      logger,
	}
}

// BatchDataRequest 批量数据请求
type BatchDataRequest struct {
	BabyID    int64     `json:"baby_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	DataTypes []string  `json:"data_types"` // ["baby_info", "feeding", "sleep", "growth", "diaper"]
}

// BatchDataResponse 批量数据响应
type BatchDataResponse struct {
	Type        string                  `json:"type"`
	BabyInfo    *entity.Baby            `json:"baby_info,omitempty"`
	FeedingData []entity.FeedingRecord  `json:"feeding_data,omitempty"`
	SleepData   []entity.SleepRecord    `json:"sleep_data,omitempty"`
	GrowthData  []entity.GrowthRecord   `json:"growth_data,omitempty"`
	DiaperData  []entity.DiaperRecord   `json:"diaper_data,omitempty"`
	Errors      map[string]string       `json:"errors,omitempty"`
}

// GetToolInfo 获取工具信息
func (t *BatchDataTools) GetToolInfo() *schema.ToolInfo {
	return &schema.ToolInfo{
		Name: "get_batch_data",
		Desc: "批量获取宝宝的多种数据，支持并行查询多个数据类型，提高效率",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"baby_id": {
				Type:     schema.Number,
				Desc:     "宝宝ID",
				Required: true,
			},
			"start_date": {
				Type:     schema.String,
				Desc:     "开始日期 (格式: 2006-01-02)",
				Required: true,
			},
			"end_date": {
				Type:     schema.String,
				Desc:     "结束日期 (格式: 2006-01-02)",
				Required: true,
			},
			"data_types": {
				Type:     schema.Array,
				Desc:     "需要获取的数据类型数组，可选值: baby_info, feeding, sleep, growth, diaper",
				Required: true,
			},
		}),
	}
}

// Execute 执行批量数据查询
func (t *BatchDataTools) Execute(ctx context.Context, params map[string]interface{}) (string, error) {
	// 解析参数
	babyID := int64(params["baby_id"].(float64))
	startDateStr := params["start_date"].(string)
	endDateStr := params["end_date"].(string)
	dataTypesRaw := params["data_types"].([]interface{})

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return "", fmt.Errorf("解析开始日期失败: %v", err)
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return "", fmt.Errorf("解析结束日期失败: %v", err)
	}

	dataTypes := make([]string, len(dataTypesRaw))
	for i, dt := range dataTypesRaw {
		dataTypes[i] = dt.(string)
	}

	// 并行获取数据
	response := &BatchDataResponse{
		Type:   "batch_data",
		Errors: make(map[string]string),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, dataType := range dataTypes {
		wg.Add(1)
		go func(dt string) {
			defer wg.Done()

			switch dt {
			case "baby_info":
				baby, err := t.babyRepo.FindByID(ctx, babyID)
				mu.Lock()
				if err != nil {
					response.Errors["baby_info"] = err.Error()
				} else {
					response.BabyInfo = baby
				}
				mu.Unlock()

			case "feeding":
				startTimeMs := startDate.UnixMilli()
				endTimeMs := endDate.UnixMilli()
				records, _, err := t.feedingRepo.FindByBabyID(ctx, babyID, startTimeMs, endTimeMs, 1, 1000)
				mu.Lock()
				if err != nil {
					response.Errors["feeding"] = err.Error()
				} else {
					// 转换为非指针切片
					feedingData := make([]entity.FeedingRecord, len(records))
					for i, r := range records {
						feedingData[i] = *r
					}
					response.FeedingData = feedingData
				}
				mu.Unlock()

			case "sleep":
				startTimeMs := startDate.UnixMilli()
				endTimeMs := endDate.UnixMilli()
				records, _, err := t.sleepRepo.FindByBabyID(ctx, babyID, startTimeMs, endTimeMs, 1, 1000)
				mu.Lock()
				if err != nil {
					response.Errors["sleep"] = err.Error()
				} else {
					sleepData := make([]entity.SleepRecord, len(records))
					for i, r := range records {
						sleepData[i] = *r
					}
					response.SleepData = sleepData
				}
				mu.Unlock()

			case "growth":
				startTimeMs := startDate.UnixMilli()
				endTimeMs := endDate.UnixMilli()
				records, _, err := t.growthRepo.FindByBabyID(ctx, babyID, startTimeMs, endTimeMs, 1, 1000)
				mu.Lock()
				if err != nil {
					response.Errors["growth"] = err.Error()
				} else {
					growthData := make([]entity.GrowthRecord, len(records))
					for i, r := range records {
						growthData[i] = *r
					}
					response.GrowthData = growthData
				}
				mu.Unlock()

			case "diaper":
				startTimeMs := startDate.UnixMilli()
				endTimeMs := endDate.UnixMilli()
				records, _, err := t.diaperRepo.FindByBabyID(ctx, babyID, startTimeMs, endTimeMs, 1, 1000)
				mu.Lock()
				if err != nil {
					response.Errors["diaper"] = err.Error()
				} else {
					diaperData := make([]entity.DiaperRecord, len(records))
					for i, r := range records {
						diaperData[i] = *r
					}
					response.DiaperData = diaperData
				}
				mu.Unlock()

			default:
				mu.Lock()
				response.Errors[dt] = fmt.Sprintf("未知的数据类型: %s", dt)
				mu.Unlock()
			}
		}(dataType)
	}

	wg.Wait()

	// 序列化响应
	result, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("序列化响应失败: %v", err)
	}

	t.logger.Debug("批量数据查询完成",
		zap.Int64("baby_id", babyID),
		zap.Strings("data_types", dataTypes),
		zap.Int("errors_count", len(response.Errors)),
	)

	return string(result), nil
}
