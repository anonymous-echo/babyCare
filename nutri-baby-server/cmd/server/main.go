package main

// @title 宝宝喂养日志 API
// @version 1.0
// @description 全栈育儿记录系统 API 文档
// @host localhost:8080
// @BasePath /
// @schemes http https
// @x-logo {"url":"https://example.com/logo.png"}

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/config"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/logger"
	"github.com/wxlbd/nutri-baby-server/wire"
	"go.uber.org/zap"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "", "Configuration file path (e.g., config/config.yaml)")
	flag.Parse()

	// 确定配置文件路径：命令行 > 环境变量 > 默认值
	cfg, err := loadConfigWithFallback(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.Init(cfg.Log); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting Nutri Baby Server...")

	// 通过Wire初始化应用
	app, err := wire.InitApp(cfg)
	if err != nil {
		logger.Fatal("Failed to init app", zap.Error(err))
	}

	// 启动定时任务服务
	app.Scheduler.Start()
	defer app.Scheduler.Stop()

	// 启动HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      app.Router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// 在goroutine中启动服务器
	go func() {
		logger.Info("Server is running",
			zap.String("addr", server.Addr),
			zap.String("mode", cfg.Server.Mode))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

// loadConfigWithFallback 按优先级加载配置文件
// 优先级: 命令行flag > 环境变量 > 默认值
func loadConfigWithFallback(flagPath string) (*config.Config, error) {
	// 优先级1: 命令行参数
	if flagPath != "" {
		return config.Load(flagPath)
	}

	// 优先级2: 环境变量 NUTRI_BABY_CONFIG
	if envPath := os.Getenv("NUTRI_BABY_CONFIG"); envPath != "" {
		return config.Load(envPath)
	}

	// 优先级3: 默认值
	defaultPath := "config/config.yaml"

	// 检查配置文件是否存在
	if _, err := os.Stat(defaultPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(
				"configuration file not found\n"+
					"Please specify config path using one of these methods:\n"+
					"  1. Command line flag: -config /path/to/config.yaml\n"+
					"  2. Environment variable: export NUTRI_BABY_CONFIG=/path/to/config.yaml\n"+
					"  3. Default path: %s",
				defaultPath,
			)
		}
		return nil, fmt.Errorf("failed to stat config file: %w", err)
	}

	return config.Load(defaultPath)
}
