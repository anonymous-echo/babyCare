package persistence

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/config"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/logger"
)

// NewDatabase 创建数据库连接
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	// GORM配置
	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
		NowFunc: func() time.Time {
			// 返回 UTC 时间，PostgreSQL 会根据 DSN 中的 timezone 参数自动转换
			return time.Now().UTC()
		},
		// 禁用外键约束检查，避免迁移顺序问题
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 获取底层连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	// 自动迁移
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	logger.Info("Database connected successfully")

	return db, nil
}

// autoMigrate 自动迁移数据表 (去家庭化架构)
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.Baby{},
		&entity.BabyCollaborator{}, // 去家庭化架构：宝宝协作者
		&entity.BabyInvitation{},   // 去家庭化架构：宝宝邀请(微信分享/二维码)
		&entity.FeedingRecord{},
		&entity.SleepRecord{},
		&entity.DiaperRecord{},
		&entity.GrowthRecord{},
		&entity.VaccinePlanTemplate{},
		&entity.BabyVaccineSchedule{}, // 新表：合并计划、记录和提醒
		&entity.SubscribeRecord{},     // 订阅消息：用户订阅记录
		&entity.MessageSendLog{},      // 订阅消息：消息发送日志
		&entity.MessageSendQueue{},    // 订阅消息：消息发送队列
	)
}
