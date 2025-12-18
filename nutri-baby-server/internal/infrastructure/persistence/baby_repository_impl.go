package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// babyRepositoryImpl 宝宝仓储实现 (去家庭化架构)
type babyRepositoryImpl struct {
	db *gorm.DB
}

// NewBabyRepository 创建宝宝仓储
func NewBabyRepository(db *gorm.DB) repository.BabyRepository {
	return &babyRepositoryImpl{db: db}
}

// Create 创建宝宝
func (r *babyRepositoryImpl) Create(ctx context.Context, baby *entity.Baby) error {
	if err := r.db.WithContext(ctx).Create(baby).Error; err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to create baby", err)
	}
	return nil
}

// FindByID 根据ID查找宝宝
func (r *babyRepositoryImpl) FindByID(ctx context.Context, babyID int64) (*entity.Baby, error) {
	var baby entity.Baby
	err := r.db.WithContext(ctx).
		Where("id = ?", babyID).
		First(&baby).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New(errors.NotFound, "baby not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find baby", err)
	}

	return &baby, nil
}

// FindByUserID 查找用户可访问的宝宝列表(通过协作者关系)
func (r *babyRepositoryImpl) FindByUserID(ctx context.Context, userID int64) ([]*entity.Baby, error) {
	var babies []*entity.Baby

	// 通过 baby_collaborators 表关联查询
	err := r.db.WithContext(ctx).
		Joins("JOIN baby_collaborators ON baby_collaborators.baby_id = babies.id").
		Where("baby_collaborators.user_id = ?", userID).
		Order("baby_collaborators.created_at DESC").
		Find(&babies).Error

	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find babies by user id", err)
	}

	return babies, nil
}

// Update 更新宝宝信息
func (r *babyRepositoryImpl) Update(ctx context.Context, baby *entity.Baby) error {
	err := r.db.WithContext(ctx).
		Model(&entity.Baby{}).
		Where("id = ?", baby.ID).
		Updates(baby).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to update baby", err)
	}

	return nil
}

// Delete 删除宝宝(软删除)
func (r *babyRepositoryImpl) Delete(ctx context.Context, babyID int64) error {
	err := r.db.WithContext(ctx).
		Model(&entity.Baby{}).
		Where("id = ?", babyID).
		Delete(&entity.Baby{}).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to delete baby", err)
	}

	return nil
}

// FindByCreator 查找用户创建的宝宝列表
func (r *babyRepositoryImpl) FindByCreator(ctx context.Context, creatorID int64) ([]*entity.Baby, error) {
	var babies []*entity.Baby

	err := r.db.WithContext(ctx).
		Where("user_id = ?", creatorID).
		Order("created_at DESC").
		Find(&babies).Error

	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find babies by creator", err)
	}

	return babies, nil
}

// FindAll 查找所有宝宝
func (r *babyRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Baby, error) {
	var babies []*entity.Baby

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&babies).Error
	// gorm find方法不会返回 gorm.ErrRecordNotFound 错误
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find all babies", err)
	}

	return babies, nil
}

// FindActiveBabies 查找活跃宝宝（创建者或协作者在 activeSince 后登录过）
func (r *babyRepositoryImpl) FindActiveBabies(ctx context.Context, activeSince int64) ([]*entity.Baby, error) {
	var babies []*entity.Baby

	// 查找活跃用户ID
	var activeUserIDs []int64
	if err := r.db.WithContext(ctx).Model(&entity.User{}).
		Where("last_login_time > ?", activeSince).
		Pluck("id", &activeUserIDs).Error; err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find active users", err)
	}

	if len(activeUserIDs) == 0 {
		return []*entity.Baby{}, nil
	}

	// 查找这些用户创建的宝宝 或 作为协作者的宝宝
	// 使用 Distinct 去重
	err := r.db.WithContext(ctx).
		Distinct("babies.*").
		Model(&entity.Baby{}).
		Joins("LEFT JOIN baby_collaborators bc ON babies.id = bc.baby_id").
		Where("babies.user_id IN (?) OR bc.user_id IN (?)", activeUserIDs, activeUserIDs).
		Where("babies.deleted_at = 0").
		Find(&babies).Error

	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find active babies", err)
	}

	return babies, nil
}
