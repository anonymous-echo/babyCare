package wire

import (
	"github.com/gin-gonic/gin"
	"github.com/wxlbd/nutri-baby-server/internal/application/service"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/config"
	"github.com/wxlbd/nutri-baby-server/internal/interface/http/handler"
)

// App 应用程序
type App struct {
	Config            *config.Config
	Router            *gin.Engine
	Scheduler         *service.SchedulerService
	AIAnalysisService service.AIAnalysisService
	AIAnalysisHandler *handler.AIAnalysisHandler
}

// NewApp 创建应用实例
func NewApp(
	cfg *config.Config,
	router *gin.Engine,
	scheduler *service.SchedulerService,
	aiAnalysisService service.AIAnalysisService,
	aiAnalysisHandler *handler.AIAnalysisHandler,
) *App {
	return &App{
		Config:            cfg,
		Router:            router,
		Scheduler:         scheduler,
		AIAnalysisService: aiAnalysisService,
		AIAnalysisHandler: aiAnalysisHandler,
	}
}
