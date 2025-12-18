package repository

import (
	"context"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
)

// BabyInvitationRepository 宝宝邀请仓储接口
type BabyInvitationRepository interface {
	// Create 创建邀请
	Create(ctx context.Context, invitation *entity.BabyInvitation) error

	// FindByToken 根据token查找邀请
	FindByToken(ctx context.Context, token string) (*entity.BabyInvitation, error)

	// FindByShortCode 根据短码查找邀请
	FindByShortCode(ctx context.Context, shortCode string) (*entity.BabyInvitation, error)

	// FindByBabyID 查找宝宝的所有邀请记录
	FindByBabyID(ctx context.Context, babyID int64) ([]*entity.BabyInvitation, error)

	// FindByBabyAndInviter 根据宝宝ID和邀请人查找未使用的邀请记录
	// 用于实现邀请码唯一性:同一用户对同一宝宝只有一个有效邀请
	FindByBabyAndInviter(ctx context.Context, babyID int64, inviterID int64) (*entity.BabyInvitation, error)

	// Delete 删除邀请(软删除)
	Delete(ctx context.Context, invitationID int64) error

	// CleanExpired 清理过期的邀请
	CleanExpired(ctx context.Context) error
}
