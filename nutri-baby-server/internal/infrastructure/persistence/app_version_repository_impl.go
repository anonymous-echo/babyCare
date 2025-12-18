package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// appVersionRepositoryImpl 应用版本仓储实现
type appVersionRepositoryImpl struct {
	db *gorm.DB
}

// NewAppVersionRepository 创建应用版本仓储
func NewAppVersionRepository(db *gorm.DB) repository.AppVersionRepository {
	return &appVersionRepositoryImpl{db: db}
}

// FindActive 获取当前活跃版本
func (r *appVersionRepositoryImpl) FindActive(ctx context.Context) (*entity.AppVersion, error) {
	var appVersion entity.AppVersion
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("created_at DESC").
		First(&appVersion).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New(errors.NotFound, "active app version not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find active app version", err)
	}

	return &appVersion, nil
}

// FindByVersion 根据版本号查找版本信息
func (r *appVersionRepositoryImpl) FindByVersion(ctx context.Context, version string) (*entity.AppVersion, error) {
	var appVersion entity.AppVersion
	err := r.db.WithContext(ctx).
		Where("version = ?", version).
		First(&appVersion).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New(errors.NotFound, "app version not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find app version", err)
	}

	return &appVersion, nil
}

// Create 创建版本信息
func (r *appVersionRepositoryImpl) Create(ctx context.Context, appVersion *entity.AppVersion) error {
	if err := r.db.WithContext(ctx).Create(appVersion).Error; err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to create app version", err)
	}
	return nil
}

// Update 更新版本信息
func (r *appVersionRepositoryImpl) Update(ctx context.Context, appVersion *entity.AppVersion) error {
	if err := r.db.WithContext(ctx).Save(appVersion).Error; err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to update app version", err)
	}
	return nil
}

// SetActive 设置版本为活跃版本（将其他版本设为非活跃）
func (r *appVersionRepositoryImpl) SetActive(ctx context.Context, version string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 将所有版本设为非活跃
		if err := tx.WithContext(ctx).
			Model(&entity.AppVersion{}).
			Update("is_active", false).Error; err != nil {
			return errors.Wrap(errors.DatabaseError, "failed to deactivate all versions", err)
		}

		// 2. 将指定版本设为活跃
		if err := tx.WithContext(ctx).
			Model(&entity.AppVersion{}).
			Where("version = ?", version).
			Update("is_active", true).Error; err != nil {
			return errors.Wrap(errors.DatabaseError, "failed to activate version", err)
		}

		return nil
	})
}
