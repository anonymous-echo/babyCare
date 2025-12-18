package service

import (
	"context"
	"strconv"

	"go.uber.org/zap"

	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// BaseRecordService 基础记录服务，提供所有记录服务的共享逻辑
type BaseRecordService struct {
	babyRepo         repository.BabyRepository
	collaboratorRepo repository.BabyCollaboratorRepository
	userRepo         repository.UserRepository
	logger           *zap.Logger
}

// NewBaseRecordService 创建基础记录服务
func NewBaseRecordService(
	babyRepo repository.BabyRepository,
	collaboratorRepo repository.BabyCollaboratorRepository,
	userRepo repository.UserRepository,
	logger *zap.Logger,
) *BaseRecordService {
	return &BaseRecordService{
		babyRepo:         babyRepo,
		collaboratorRepo: collaboratorRepo,
		userRepo:         userRepo,
		logger:           logger,
	}
}

// CheckBabyAccess 检查用户对宝宝的访问权限 (去家庭化架构)
func (s *BaseRecordService) CheckBabyAccess(ctx context.Context, babyID, openID string) error {
	// 转换babyID from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return errors.New(errors.ParamError, "invalid baby id format")
	}

	// 获取用户信息以获取UserID
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return err
	}

	// 检查用户是否为宝宝的协作者
	isCollaborator, err := s.collaboratorRepo.IsCollaborator(ctx, babyIDInt64, user.ID)
	if err != nil {
		return err
	}

	if !isCollaborator {
		return errors.New(errors.PermissionDenied, "您没有权限访问该宝宝的记录")
	}

	return nil
}
