package persistence

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// babyCollaboratorRepositoryImpl 宝宝协作者仓储实现
type babyCollaboratorRepositoryImpl struct {
	db *gorm.DB
}

// NewBabyCollaboratorRepository 创建宝宝协作者仓储
func NewBabyCollaboratorRepository(db *gorm.DB) repository.BabyCollaboratorRepository {
	return &babyCollaboratorRepositoryImpl{db: db}
}

// Create 创建协作者
func (r *babyCollaboratorRepositoryImpl) Create(ctx context.Context, collaborator *entity.BabyCollaborator) error {
	if err := r.db.WithContext(ctx).Create(collaborator).Error; err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to create baby collaborator", err)
	}
	return nil
}

// FindByBabyID 获取宝宝的所有协作者
func (r *babyCollaboratorRepositoryImpl) FindByBabyID(ctx context.Context, babyID int64) ([]*entity.BabyCollaborator, error) {
	var collaborators []*entity.BabyCollaborator
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("baby_id = ?", babyID).
		Order("created_at ASC").
		Find(&collaborators).Error

	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find collaborators by baby id", err)
	}

	return collaborators, nil
}

// FindByUserID 获取用户的所有协作宝宝
func (r *babyCollaboratorRepositoryImpl) FindByUserID(ctx context.Context, userID int64) ([]*entity.BabyCollaborator, error) {
	var collaborators []*entity.BabyCollaborator
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&collaborators).Error

	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find collaborators by user id", err)
	}

	return collaborators, nil
}

// FindByBabyAndUser 查找特定宝宝的特定协作者
func (r *babyCollaboratorRepositoryImpl) FindByBabyAndUser(ctx context.Context, babyID, userID int64) (*entity.BabyCollaborator, error) {
	var collaborator entity.BabyCollaborator
	err := r.db.WithContext(ctx).
		Where("baby_id = ? AND user_id = ?", babyID, userID).
		First(&collaborator).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // 不是协作者,返回 nil 而不是错误
	}
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find collaborator", err)
	}

	return &collaborator, nil
}

// CheckPermission 检查用户对宝宝的访问权限
func (r *babyCollaboratorRepositoryImpl) CheckPermission(ctx context.Context, babyID, userID int64) (*entity.BabyCollaborator, error) {
	collaborator, err := r.FindByBabyAndUser(ctx, babyID, userID)
	if err != nil {
		return nil, err
	}

	// 检查是否过期
	if collaborator != nil && collaborator.IsExpired() {
		return nil, nil // 权限已过期,返回 nil
	}

	return collaborator, nil
}

// Update 更新协作者信息
func (r *babyCollaboratorRepositoryImpl) Update(ctx context.Context, collaborator *entity.BabyCollaborator) error {
	err := r.db.WithContext(ctx).
		Model(&entity.BabyCollaborator{}).
		Where("baby_id = ? AND user_id = ?", collaborator.BabyID, collaborator.UserID).
		Updates(collaborator).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to update collaborator", err)
	}

	return nil
}

// Delete 移除协作者(软删除)
func (r *babyCollaboratorRepositoryImpl) Delete(ctx context.Context, babyID, userID int64) error {
	err := r.db.WithContext(ctx).
		Model(&entity.BabyCollaborator{}).
		Where("baby_id = ? AND user_id = ?", babyID, userID).
		Delete(&entity.BabyCollaborator{}).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to delete collaborator", err)
	}

	return nil
}

// BatchCreate 批量创建协作者(用于复制协作者列表)
func (r *babyCollaboratorRepositoryImpl) BatchCreate(ctx context.Context, collaborators []*entity.BabyCollaborator) error {
	if len(collaborators) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).CreateInBatches(collaborators, 100).Error; err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to batch create collaborators", err)
	}

	return nil
}

// CleanExpired 清理过期的临时协作者
func (r *babyCollaboratorRepositoryImpl) CleanExpired(ctx context.Context) error {
	now := time.Now().UnixMilli()

	err := r.db.WithContext(ctx).
		Model(&entity.BabyCollaborator{}).
		Where("access_type = ? AND expires_at IS NOT NULL AND expires_at < ?", "temporary", now).
		Delete(&entity.BabyCollaborator{}).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to clean expired collaborators", err)
	}

	return nil
}

// IsCollaborator 检查是否是协作者
func (r *babyCollaboratorRepositoryImpl) IsCollaborator(ctx context.Context, babyID, userID int64) (bool, error) {
	collaborator, err := r.CheckPermission(ctx, babyID, userID)
	if err != nil {
		return false, err
	}
	return collaborator != nil, nil
}

// IsAdmin 检查是否是管理员
func (r *babyCollaboratorRepositoryImpl) IsAdmin(ctx context.Context, babyID, userID int64) (bool, error) {
	collaborator, err := r.CheckPermission(ctx, babyID, userID)
	if err != nil {
		return false, err
	}
	if collaborator == nil {
		return false, nil
	}
	return collaborator.IsAdmin(), nil
}

// CanEdit 检查是否有编辑权限
func (r *babyCollaboratorRepositoryImpl) CanEdit(ctx context.Context, babyID, userID int64) (bool, error) {
	collaborator, err := r.CheckPermission(ctx, babyID, userID)
	if err != nil {
		return false, err
	}
	if collaborator == nil {
		return false, nil
	}
	return collaborator.CanEdit(), nil
}
