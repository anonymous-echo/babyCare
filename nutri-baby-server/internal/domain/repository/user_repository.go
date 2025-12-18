package repository

import (
	"context"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *entity.User) error
	// FindByOpenID 根据OpenID查找用户
	FindByOpenID(ctx context.Context, openID string) (*entity.User, error)
	// FindByID 根据ID查找用户
	FindByID(ctx context.Context, userID int64) (*entity.User, error)
	// Update 更新用户
	Update(ctx context.Context, user *entity.User) error
	// UpdateLastLoginTime 更新最后登录时间
	UpdateLastLoginTime(ctx context.Context, openID string) error
	// UpdateDefaultBabyID 更新默认宝宝ID
	UpdateDefaultBabyID(ctx context.Context, openID string, babyID int64) error
}

// BabyRepository 宝宝仓储接口 (去家庭化架构)
type BabyRepository interface {
	// Create 创建宝宝
	Create(ctx context.Context, baby *entity.Baby) error
	// FindByID 根据ID查找宝宝
	FindByID(ctx context.Context, babyID int64) (*entity.Baby, error)
	// FindByUserID 查找用户可访问的宝宝列表(通过协作者关系)
	FindByUserID(ctx context.Context, userID int64) ([]*entity.Baby, error)

	// Update 更新宝宝信息
	Update(ctx context.Context, baby *entity.Baby) error
	// Delete 删除宝宝(软删除)
	Delete(ctx context.Context, babyID int64) error
	// FindByCreator 查找用户创建的宝宝列表
	FindByCreator(ctx context.Context, creatorID int64) ([]*entity.Baby, error)
	// FindAll 查找所有宝宝
	FindAll(ctx context.Context) ([]*entity.Baby, error)
	// FindActiveBabies 查找活跃宝宝（创建者或协作者在 activeSince 后登录过）
	FindActiveBabies(ctx context.Context, activeSince int64) ([]*entity.Baby, error)
}
