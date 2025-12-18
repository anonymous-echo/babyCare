package service

import (
	"context"

	"github.com/wxlbd/nutri-baby-server/internal/application/dto"
	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// AppVersionService 应用版本服务
type AppVersionService struct {
	appVersionRepo repository.AppVersionRepository
}

// NewAppVersionService 创建应用版本服务
func NewAppVersionService(appVersionRepo repository.AppVersionRepository) *AppVersionService {
	return &AppVersionService{
		appVersionRepo: appVersionRepo,
	}
}

// GetCurrentVersion 获取当前活跃版本信息
func (s *AppVersionService) GetCurrentVersion(ctx context.Context) (*dto.AppVersionDTO, error) {
	appVersion, err := s.appVersionRepo.FindActive(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.AppVersionDTO{
		Version:      appVersion.Version,
		Name:         appVersion.Name,
		Description:  appVersion.Description,
		MinVersion:   appVersion.MinVersion,
		ForceUpdate:  appVersion.ForceUpdate,
		ReleaseNotes: appVersion.ReleaseNotes,
		BuildTime:    appVersion.BuildTime.UnixMilli(),
	}, nil
}

// GetVersionByNumber 根据版本号获取版本信息
func (s *AppVersionService) GetVersionByNumber(ctx context.Context, version string) (*dto.AppVersionDTO, error) {
	appVersion, err := s.appVersionRepo.FindByVersion(ctx, version)
	if err != nil {
		return nil, err
	}

	return &dto.AppVersionDTO{
		Version:      appVersion.Version,
		Name:         appVersion.Name,
		Description:  appVersion.Description,
		MinVersion:   appVersion.MinVersion,
		ForceUpdate:  appVersion.ForceUpdate,
		ReleaseNotes: appVersion.ReleaseNotes,
		BuildTime:    appVersion.BuildTime.UnixMilli(),
	}, nil
}

// CreateVersion 创建新版本
func (s *AppVersionService) CreateVersion(ctx context.Context, req *dto.CreateAppVersionRequest) (*dto.AppVersionDTO, error) {
	// 验证版本号格式
	if req.Version == "" {
		return nil, errors.New(errors.ParamError, "version is required")
	}

	appVersion := &entity.AppVersion{
		Version:      req.Version,
		Name:         req.Name,
		Description:  req.Description,
		MinVersion:   req.MinVersion,
		IsActive:     req.IsActive,
		ForceUpdate:  req.ForceUpdate,
		ReleaseNotes: req.ReleaseNotes,
	}

	if err := s.appVersionRepo.Create(ctx, appVersion); err != nil {
		return nil, err
	}

	return &dto.AppVersionDTO{
		Version:      appVersion.Version,
		Name:         appVersion.Name,
		Description:  appVersion.Description,
		MinVersion:   appVersion.MinVersion,
		ForceUpdate:  appVersion.ForceUpdate,
		ReleaseNotes: appVersion.ReleaseNotes,
		BuildTime:    appVersion.BuildTime.UnixMilli(),
	}, nil
}

// SetActiveVersion 设置活跃版本
func (s *AppVersionService) SetActiveVersion(ctx context.Context, version string) error {
	// 验证版本是否存在
	_, err := s.appVersionRepo.FindByVersion(ctx, version)
	if err != nil {
		return err
	}

	// 设置为活跃版本
	return s.appVersionRepo.SetActive(ctx, version)
}
