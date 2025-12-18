package repository

import (
	"context"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
)

// BabyCollaboratorRepository 宝宝协作者仓储接口
type BabyCollaboratorRepository interface {
	// Create 创建协作者
	Create(ctx context.Context, collaborator *entity.BabyCollaborator) error

	// FindByBabyID 获取宝宝的所有协作者
	FindByBabyID(ctx context.Context, babyID int64) ([]*entity.BabyCollaborator, error)

	// FindByUserID 获取用户的所有协作宝宝
	FindByUserID(ctx context.Context, userID int64) ([]*entity.BabyCollaborator, error)

	// FindByBabyAndUser 查找特定宝宝的特定协作者
	FindByBabyAndUser(ctx context.Context, babyID int64, userID int64) (*entity.BabyCollaborator, error)

	// CheckPermission 检查用户对宝宝的访问权限
	// 返回协作者信息,如果没有权限返回 nil
	CheckPermission(ctx context.Context, babyID int64, userID int64) (*entity.BabyCollaborator, error)

	// Update 更新协作者信息
	Update(ctx context.Context, collaborator *entity.BabyCollaborator) error

	// Delete 移除协作者(软删除)
	Delete(ctx context.Context, babyID int64, userID int64) error

	// BatchCreate 批量创建协作者(用于复制协作者列表)
	BatchCreate(ctx context.Context, collaborators []*entity.BabyCollaborator) error

	// CleanExpired 清理过期的临时协作者
	CleanExpired(ctx context.Context) error

	// IsCollaborator 检查是否是协作者
	IsCollaborator(ctx context.Context, babyID int64, userID int64) (bool, error)

	// IsAdmin 检查是否是管理员
	IsAdmin(ctx context.Context, babyID int64, userID int64) (bool, error)

	// CanEdit 检查是否有编辑权限
	CanEdit(ctx context.Context, babyID int64, userID int64) (bool, error)
}
