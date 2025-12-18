package entity

import "time"

// AppVersion 应用版本实体
type AppVersion struct {
	ID           int64     `gorm:"primaryKey;column:id" json:"id"`
	Version      string    `gorm:"column:version;type:varchar(20);uniqueIndex;not null" json:"version"`                   // 版本号
	Name         string    `gorm:"column:name;type:varchar(100);not null;default:'宝宝喂养时刻'" json:"name"`              // 应用名称
	Description  string    `gorm:"column:description;type:text" json:"description"`                                       // 版本描述
	MinVersion   string    `gorm:"column:min_version;type:varchar(20)" json:"minVersion"`                                 // 最小支持版本
	IsActive     bool      `gorm:"column:is_active;type:boolean;not null;default:true;index" json:"isActive"`            // 是否为活跃版本
	ForceUpdate  bool      `gorm:"column:force_update;type:boolean;not null;default:false" json:"forceUpdate"`           // 是否强制更新
	ReleaseNotes string    `gorm:"column:release_notes;type:text" json:"releaseNotes"`                                    // 发布说明
	BuildTime    time.Time `gorm:"column:build_time;type:timestamp;default:CURRENT_TIMESTAMP" json:"buildTime"`          // 构建时间
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"createdAt"`          // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updatedAt"`          // 更新时间
}

// TableName 指定表名
func (AppVersion) TableName() string {
	return "app_versions"
}
