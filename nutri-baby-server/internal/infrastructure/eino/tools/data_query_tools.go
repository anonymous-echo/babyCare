package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"go.uber.org/zap"
)

// DataQueryTools 数据查询工具集
type DataQueryTools struct {
	feedingRepo repository.FeedingRecordRepository
	sleepRepo   repository.SleepRecordRepository
	diaperRepo  repository.DiaperRecordRepository
	growthRepo  repository.GrowthRecordRepository
	vaccineRepo repository.BabyVaccineScheduleRepository
	babyRepo    repository.BabyRepository
	logger      *zap.Logger
}

// NewDataQueryTools 创建数据查询工具集
func NewDataQueryTools(
	feedingRepo repository.FeedingRecordRepository,
	sleepRepo repository.SleepRecordRepository,
	diaperRepo repository.DiaperRecordRepository,
	growthRepo repository.GrowthRecordRepository,
	vaccineRepo repository.BabyVaccineScheduleRepository,
	babyRepo repository.BabyRepository,
	logger *zap.Logger,
) *DataQueryTools {
	return &DataQueryTools{
		feedingRepo: feedingRepo,
		sleepRepo:   sleepRepo,
		diaperRepo:  diaperRepo,
		growthRepo:  growthRepo,
		vaccineRepo: vaccineRepo,
		babyRepo:    babyRepo,
		logger:      logger,
	}
}

// GetToolInfos 获取所有工具信息
func (t *DataQueryTools) GetToolInfos() []*schema.ToolInfo {
	return []*schema.ToolInfo{
		t.getFeedingDataToolInfo(),
		t.getSleepDataToolInfo(),
		t.getGrowthDataToolInfo(),
		t.getDiaperDataToolInfo(),
		t.getVaccineDataToolInfo(),
		t.getBabyInfoToolInfo(),
	}
}

// getFeedingDataToolInfo 获取喂养数据工具信息
func (t *DataQueryTools) getFeedingDataToolInfo() *schema.ToolInfo {
	return &schema.ToolInfo{
		Name: "get_feeding_data",
		Desc: "获取宝宝指定时间范围内的喂养记录数据，包括喂养类型、奶量、时长等信息",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"baby_id": {
				Type: "integer",
				Desc: "宝宝ID",
			},
			"start_date": {
				Type: "string",
				Desc: "开始日期，格式：YYYY-MM-DD",
			},
			"end_date": {
				Type: "string",
				Desc: "结束日期，格式：YYYY-MM-DD",
			},
			"limit": {
				Type: "integer",
				Desc: "返回记录数量限制，默认100",
			},
		}),
	}
}

// getSleepDataToolInfo 获取睡眠数据工具信息
func (t *DataQueryTools) getSleepDataToolInfo() *schema.ToolInfo {
	return &schema.ToolInfo{
		Name: "get_sleep_data",
		Desc: "获取宝宝指定时间范围内的睡眠记录数据，包括睡眠时长、质量等信息",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"baby_id": {
				Type: "integer",
				Desc: "宝宝ID",
			},
			"start_date": {
				Type: "string",
				Desc: "开始日期，格式：YYYY-MM-DD",
			},
			"end_date": {
				Type: "string",
				Desc: "结束日期，格式：YYYY-MM-DD",
			},
			"limit": {
				Type: "integer",
				Desc: "返回记录数量限制，默认100",
			},
		}),
	}
}

// getGrowthDataToolInfo 获取成长数据工具信息
func (t *DataQueryTools) getGrowthDataToolInfo() *schema.ToolInfo {
	return &schema.ToolInfo{
		Name: "get_growth_data",
		Desc: "获取宝宝指定时间范围内的成长记录数据，包括身高、体重、头围等信息",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"baby_id": {
				Type: "integer",
				Desc: "宝宝ID",
			},
			"start_date": {
				Type: "string",
				Desc: "开始日期，格式：YYYY-MM-DD",
			},
			"end_date": {
				Type: "string",
				Desc: "结束日期，格式：YYYY-MM-DD",
			},
			"limit": {
				Type: "integer",
				Desc: "返回记录数量限制，默认100",
			},
		}),
	}
}

// getDiaperDataToolInfo 获取尿布数据工具信息
func (t *DataQueryTools) getDiaperDataToolInfo() *schema.ToolInfo {
	return &schema.ToolInfo{
		Name: "get_diaper_data",
		Desc: "获取宝宝指定时间范围内的尿布记录数据，包括排便类型、次数等信息",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"baby_id": {
				Type: "integer",
				Desc: "宝宝ID",
			},
			"start_date": {
				Type: "string",
				Desc: "开始日期，格式：YYYY-MM-DD",
			},
			"end_date": {
				Type: "string",
				Desc: "结束日期，格式：YYYY-MM-DD",
			},
			"limit": {
				Type: "integer",
				Desc: "返回记录数量限制，默认100",
			},
		}),
	}
}

// getVaccineDataToolInfo 获取疫苗数据工具信息
func (t *DataQueryTools) getVaccineDataToolInfo() *schema.ToolInfo {
	return &schema.ToolInfo{
		Name: "get_vaccine_data",
		Desc: "获取宝宝的疫苗接种记录和计划",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"baby_id": {
				Type: "integer",
				Desc: "宝宝ID",
			},
		}),
	}
}

// getBabyInfoToolInfo 获取宝宝信息工具信息
func (t *DataQueryTools) getBabyInfoToolInfo() *schema.ToolInfo {
	return &schema.ToolInfo{
		Name: "get_baby_info",
		Desc: "获取宝宝的基本信息，包括姓名、性别、出生日期、月龄等",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"baby_id": {
				Type: "integer",
				Desc: "宝宝ID",
			},
		}),
	}
}

// ExecuteTool 执行工具调用
func (t *DataQueryTools) ExecuteTool(ctx context.Context, toolName string, params map[string]interface{}) (string, error) {
	switch toolName {
	case "get_feeding_data":
		return t.getFeedingData(ctx, params)
	case "get_sleep_data":
		return t.getSleepData(ctx, params)
	case "get_growth_data":
		return t.getGrowthData(ctx, params)
	case "get_diaper_data":
		return t.getDiaperData(ctx, params)
	case "get_vaccine_data":
		return t.getVaccineData(ctx, params)
	case "get_baby_info":
		return t.getBabyInfo(ctx, params)
	default:
		return "", fmt.Errorf("未知的工具: %s", toolName)
	}
}

// getFeedingData 获取喂养数据
func (t *DataQueryTools) getFeedingData(ctx context.Context, params map[string]interface{}) (string, error) {
	babyID, startTime, endTime, limit, err := t.parseCommonParams(params)
	if err != nil {
		return "", err
	}

	records, _, err := t.feedingRepo.FindByBabyID(ctx, babyID, startTime, endTime, 1, limit)
	if err != nil {
		t.logger.Error("获取喂养数据失败", zap.Error(err))
		return "", fmt.Errorf("获取喂养数据失败: %v", err)
	}

	result := map[string]interface{}{
		"type":    "feeding_data",
		"count":   len(records),
		"records": records,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("序列化喂养数据失败: %v", err)
	}

	return string(data), nil
}

// getSleepData 获取睡眠数据
func (t *DataQueryTools) getSleepData(ctx context.Context, params map[string]interface{}) (string, error) {
	babyID, startTime, endTime, limit, err := t.parseCommonParams(params)
	if err != nil {
		return "", err
	}

	records, _, err := t.sleepRepo.FindByBabyID(ctx, babyID, startTime, endTime, 1, limit)
	if err != nil {
		t.logger.Error("获取睡眠数据失败", zap.Error(err))
		return "", fmt.Errorf("获取睡眠数据失败: %v", err)
	}

	result := map[string]interface{}{
		"type":    "sleep_data",
		"count":   len(records),
		"records": records,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("序列化睡眠数据失败: %v", err)
	}

	return string(data), nil
}

// getGrowthData 获取成长数据
func (t *DataQueryTools) getGrowthData(ctx context.Context, params map[string]interface{}) (string, error) {
	babyID, startTime, endTime, limit, err := t.parseCommonParams(params)
	if err != nil {
		return "", err
	}

	records, _, err := t.growthRepo.FindByBabyID(ctx, babyID, startTime, endTime, 1, limit)
	if err != nil {
		t.logger.Error("获取成长数据失败", zap.Error(err))
		return "", fmt.Errorf("获取成长数据失败: %v", err)
	}

	result := map[string]interface{}{
		"type":    "growth_data",
		"count":   len(records),
		"records": records,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("序列化成长数据失败: %v", err)
	}

	return string(data), nil
}

// getDiaperData 获取尿布数据
func (t *DataQueryTools) getDiaperData(ctx context.Context, params map[string]interface{}) (string, error) {
	babyID, startTime, endTime, limit, err := t.parseCommonParams(params)
	if err != nil {
		return "", err
	}

	records, _, err := t.diaperRepo.FindByBabyID(ctx, babyID, startTime, endTime, 1, limit)
	if err != nil {
		t.logger.Error("获取尿布数据失败", zap.Error(err))
		return "", fmt.Errorf("获取尿布数据失败: %v", err)
	}

	result := map[string]interface{}{
		"type":    "diaper_data",
		"count":   len(records),
		"records": records,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("序列化尿布数据失败: %v", err)
	}

	return string(data), nil
}

// getVaccineData 获取疫苗数据
func (t *DataQueryTools) getVaccineData(ctx context.Context, params map[string]interface{}) (string, error) {
	babyIDFloat, ok := params["baby_id"].(float64)
	if !ok {
		return "", fmt.Errorf("无效的宝宝ID")
	}
	babyID := int64(babyIDFloat)

	records, err := t.vaccineRepo.FindByBabyID(ctx, babyID, 1, 100)
	if err != nil {
		t.logger.Error("获取疫苗数据失败", zap.Error(err))
		return "", fmt.Errorf("获取疫苗数据失败: %v", err)
	}

	result := map[string]interface{}{
		"type":    "vaccine_data",
		"count":   len(records),
		"records": records,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("序列化疫苗数据失败: %v", err)
	}

	return string(data), nil
}

// getBabyInfo 获取宝宝信息
func (t *DataQueryTools) getBabyInfo(ctx context.Context, params map[string]interface{}) (string, error) {
	babyIDFloat, ok := params["baby_id"].(float64)
	if !ok {
		return "", fmt.Errorf("无效的宝宝ID")
	}
	babyID := int64(babyIDFloat)

	baby, err := t.babyRepo.FindByID(ctx, babyID)
	if err != nil {
		t.logger.Error("获取宝宝信息失败", zap.Error(err))
		return "", fmt.Errorf("获取宝宝信息失败: %v", err)
	}

	// 计算月龄
	birthDate, _ := time.Parse("2006-01-02", baby.BirthDate)
	now := time.Now()
	months := (now.Year()-birthDate.Year())*12 + int(now.Month()) - int(birthDate.Month())
	if now.Day() < birthDate.Day() {
		months--
	}
	if months < 0 {
		months = 0
	}

	result := map[string]interface{}{
		"type":       "baby_info",
		"baby":       baby,
		"age_months": months,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("序列化宝宝信息失败: %v", err)
	}

	return string(data), nil
}

// parseCommonParams 解析通用参数
func (t *DataQueryTools) parseCommonParams(params map[string]interface{}) (babyID int64, startTime, endTime int64, limit int, err error) {
	// 解析宝宝ID
	babyIDFloat, ok := params["baby_id"].(float64)
	if !ok {
		err = fmt.Errorf("无效的宝宝ID")
		return
	}
	babyID = int64(babyIDFloat)

	// 解析开始日期
	startDateStr, ok := params["start_date"].(string)
	if !ok {
		err = fmt.Errorf("无效的开始日期")
		return
	}
	startDate, parseErr := time.Parse("2006-01-02", startDateStr)
	if parseErr != nil {
		err = fmt.Errorf("开始日期格式错误: %v", parseErr)
		return
	}
	startTime = startDate.Unix() * 1000

	// 解析结束日期
	endDateStr, ok := params["end_date"].(string)
	if !ok {
		err = fmt.Errorf("无效的结束日期")
		return
	}
	endDate, parseErr := time.Parse("2006-01-02", endDateStr)
	if parseErr != nil {
		err = fmt.Errorf("结束日期格式错误: %v", parseErr)
		return
	}
	endTime = endDate.Unix() * 1000

	// 解析限制数量
	if limitFloat, exists := params["limit"].(float64); exists {
		limit = int(limitFloat)
	} else {
		limit = 100 // 默认值
	}

	return
}
