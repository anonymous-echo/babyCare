package persistence

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// babyInvitationRepositoryImpl 宝宝邀请仓储实现
type babyInvitationRepositoryImpl struct {
	db *gorm.DB
}

// NewBabyInvitationRepository 创建宝宝邀请仓储
func NewBabyInvitationRepository(db *gorm.DB) repository.BabyInvitationRepository {
	return &babyInvitationRepositoryImpl{db: db}
}

// Create 创建邀请
func (r *babyInvitationRepositoryImpl) Create(ctx context.Context, invitation *entity.BabyInvitation) error {
	if err := r.db.WithContext(ctx).Create(invitation).Error; err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to create baby invitation", err)
	}
	return nil
}

// FindByToken 根据token查找邀请
func (r *babyInvitationRepositoryImpl) FindByToken(ctx context.Context, token string) (*entity.BabyInvitation, error) {
	var invitation entity.BabyInvitation
	err := r.db.WithContext(ctx).
		Where("token = ?", token).
		First(&invitation).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New(errors.NotFound, "邀请不存在或已失效")
	}
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find invitation by token", err)
	}

	return &invitation, nil
}

// FindByShortCode 根据短码查找邀请
func (r *babyInvitationRepositoryImpl) FindByShortCode(ctx context.Context, shortCode string) (*entity.BabyInvitation, error) {
	var invitation entity.BabyInvitation
	err := r.db.WithContext(ctx).
		Where("short_code = ?", shortCode).
		First(&invitation).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New(errors.NotFound, "邀请不存在或已失效")
	}
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "通过短码查找邀请记录失败", err)
	}

	return &invitation, nil
}

// FindByBabyID 查找宝宝的所有邀请记录
func (r *babyInvitationRepositoryImpl) FindByBabyID(ctx context.Context, babyID int64) ([]*entity.BabyInvitation, error) {
	var invitations []*entity.BabyInvitation

	err := r.db.WithContext(ctx).
		Where("baby_id = ?", babyID).
		Order("created_at DESC").
		Find(&invitations).Error

	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "通过宝宝ID查找邀请记录失败", err)
	}

	return invitations, nil
}

// FindByBabyAndInviter 根据宝宝ID和邀请人查找未使用的邀请记录
func (r *babyInvitationRepositoryImpl) FindByBabyAndInviter(ctx context.Context, babyID, inviterID int64) (*entity.BabyInvitation, error) {
	var invitation entity.BabyInvitation

	// 查询条件:
	// 1. 宝宝ID匹配
	// 2. 邀请人ID匹配
	// 3. 邀请未使用(used_by IS NULL)
	// 4. 未被删除(GORM soft_delete自动处理)
	err := r.db.WithContext(ctx).
		Where("baby_id = ? AND user_id = ?", babyID, inviterID).
		First(&invitation).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "通过邀请人和宝宝ID查找邀请记录失败", err)
	}

	return &invitation, nil
}

// CleanExpired 清理过期的邀请
func (r *babyInvitationRepositoryImpl) CleanExpired(ctx context.Context) error {
	now := time.Now().UnixMilli()

	err := r.db.WithContext(ctx).
		Model(&entity.BabyInvitation{}).
		Where("valid_until < ?", now).
		Delete(&entity.BabyInvitation{}).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "清理过期邀请失败", err)
	}

	return nil
}

// Delete 删除邀请(软删除)
func (r *babyInvitationRepositoryImpl) Delete(ctx context.Context, invitationID int64) error {
	err := r.db.WithContext(ctx).
		Model(&entity.BabyInvitation{}).
		Where("id = ?", invitationID).
		Delete(&entity.BabyInvitation{}).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "删除邀请失败", err)
	}

	return nil
}
