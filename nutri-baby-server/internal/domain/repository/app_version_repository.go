package repository

import (
	"context"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
)

// AppVersionRepository 应用版本仓储接口
type AppVersionRepository interface {
	// FindActive 获取当前活跃版本
	FindActive(ctx context.Context) (*entity.AppVersion, error)
	// FindByVersion 根据版本号查找版本信息
	FindByVersion(ctx context.Context, version string) (*entity.AppVersion, error)
	// Create 创建版本信息
	Create(ctx context.Context, appVersion *entity.AppVersion) error
	// Update 更新版本信息
	Update(ctx context.Context, appVersion *entity.AppVersion) error
	// SetActive 设置版本为活跃版本（将其他版本设为非活跃）
	SetActive(ctx context.Context, version string) error
}
