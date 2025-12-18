//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/wxlbd/nutri-baby-server/internal/application/service"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/config"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/eino/chain"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/eino/model"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/eino/tools"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/logger"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/persistence"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/wechat"
	"github.com/wxlbd/nutri-baby-server/internal/interface/http/handler"
	"github.com/wxlbd/nutri-baby-server/internal/interface/http/router"
)

// InitApp 初始化应用(Wire自动生成) (去家庭化架构)
func InitApp(cfg *config.Config) (*App, error) {
	wire.Build(
		// 基础设施层
		logger.NewLogger, // 日志系统
		persistence.NewDatabase,
		persistence.NewRedis, // Redis 客户端
		wechat.NewClient,     // 微信 SDK 客户端

		// Eino AI框架（工具调用架构）
		model.NewToolCallingChatModel, // 支持工具调用的AI模型客户端
		tools.NewDataQueryTools,       // 数据查询工具集
		tools.NewBatchDataTools,       // 批量数据查询工具
		chain.NewAnalysisChainBuilder, // AI分析链构建器

		// 仓储层
		persistence.NewSubscriptionCacheRepository, // 订阅权限缓存管理器
		persistence.NewUserRepository,
		persistence.NewBabyRepository,
		persistence.NewBabyCollaboratorRepository, // 去家庭化架构：宝宝协作者仓储
		persistence.NewBabyInvitationRepository,   // 去家庭化架构：宝宝邀请仓储
		persistence.NewFeedingRecordRepository,
		persistence.NewSleepRecordRepository,
		persistence.NewDiaperRecordRepository,
		persistence.NewGrowthRecordRepository,
		persistence.NewBabyVaccineScheduleRepository, // 新增：疫苗接种日程仓储
		persistence.NewVaccinePlanTemplateRepository, // 疫苗计划模板仓储
		persistence.NewSubscribeRepository,           // 订阅消息仓储
		persistence.NewAIAnalysisRepository,          // AI分析结果仓储
		persistence.NewDailyTipsRepository,           // 每日建议仓储
		persistence.NewAppVersionRepository,          // 应用版本仓储

		// 应用服务层
		service.NewWechatService,    // 微信服务
		service.NewSubscribeService, // 订阅消息服务
		service.NewAuthService,
		service.NewBabyService,
		service.NewFeedingRecordService,   // 喂养记录服务
		service.NewSleepRecordService,     // 睡眠记录服务
		service.NewDiaperRecordService,    // 尿布记录服务
		service.NewGrowthRecordService,    // 成长记录服务
		service.NewTimelineService,        // 时间线聚合服务
		service.NewVaccineScheduleService, // 新增：疫苗接种日程服务
		service.NewStatisticsService,      // 新增：统计服务
		service.NewDailyStatsService,      // 新增：按日统计服务
		service.NewSchedulerService,       // 定时任务服务
		service.NewUploadService,          // 文件上传服务
		service.NewAIAnalysisService,      // AI分析服务（工具调用架构）
		service.NewAppVersionService,      // 应用版本服务
		// service.NewSyncService, // TODO: WebSocket同步未实现，暂时注释

		// HTTP处理器
		handler.NewAuthHandler,
		handler.NewBabyHandler,
		handler.NewRecordHandler,
		handler.NewVaccineScheduleHandler, // 新增：疫苗接种日程处理器
		handler.NewStatisticsHandler,      // 新增：统计处理器
		handler.NewDailyStatsHandler,      // 新增：按日统计处理器
		handler.NewSubscribeHandler,       // 订阅消息处理器
		handler.NewAIAnalysisHandler,      // AI分析处理器（工具调用架构）
		handler.NewSyncHandler,
		handler.NewUploadHandler, // 文件上传处理器

		// 路由
		router.NewRouter,

		// 应用
		NewApp,
	)
	return &App{}, nil
}
