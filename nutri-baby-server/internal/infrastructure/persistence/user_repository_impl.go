package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// userRepositoryImpl 用户仓储实现
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{db: db}
}

// Create 创建用户
func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to create user", err)
	}
	return nil
}

// FindByOpenID 根据OpenID查找用户
func (r *userRepositoryImpl) FindByOpenID(ctx context.Context, openID string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Where("openid = ?", openID).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrUserNotFound
	}
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find user", err)
	}

	return &user, nil
}

// FindByID 根据ID查找用户
func (r *userRepositoryImpl) FindByID(ctx context.Context, userID int64) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Where("id = ?", userID).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrUserNotFound
	}
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find user", err)
	}

	return &user, nil
}

// Update 更新用户
func (r *userRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("openid = ?", user.OpenID).
		Updates(user).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to update user", err)
	}

	return nil
}

// UpdateLastLoginTime 更新最后登录时间
func (r *userRepositoryImpl) UpdateLastLoginTime(ctx context.Context, openID string) error {
	err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("openid = ?", openID).
		Update("last_login_time", gorm.Expr("?", ctx.Value("current_time"))).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to update last login time", err)
	}

	return nil
}

// UpdateDefaultBabyID 更新默认宝宝ID
func (r *userRepositoryImpl) UpdateDefaultBabyID(ctx context.Context, openID string, babyID int64) error {
	err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("openid = ?", openID).
		Update("default_baby_id", babyID).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to update default baby id", err)
	}

	return nil
}
