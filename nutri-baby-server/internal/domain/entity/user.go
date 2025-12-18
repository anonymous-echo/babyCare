package entity

import "gorm.io/plugin/soft_delete"

// User 用户实体
type User struct {
	ID            int64                 `gorm:"primaryKey;column:id" json:"id"`                                    // 雪花ID主键
	OpenID        string                `gorm:"column:openid;type:varchar(64);uniqueIndex;not null" json:"openid"` // 微信OpenID,唯一索引
	NickName      string                `gorm:"column:nick_name;type:varchar(64)" json:"nickName"`                 // 昵称
	AvatarURL     string                `gorm:"column:avatar_url;type:varchar(512)" json:"avatarUrl"`              // 头像URL
	DefaultBabyID int64                 `gorm:"column:default_baby_id;default:0" json:"defaultBabyId"`             // 默认宝宝ID (int64)
	LastLoginTime int64                 `gorm:"column:last_login_time" json:"lastLoginTime"`                       // 最后登录时间(毫秒时间戳)
	CreatedAt     int64                 `gorm:"column:created_at;autoCreateTime:milli;default:0" json:"createdAt"` // 创建时间(毫秒时间戳)
	UpdatedAt     int64                 `gorm:"column:updated_at;autoUpdateTime:milli;default:0" json:"updatedAt"` // 更新时间(毫秒时间戳)
	DeletedAt     soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;index;default:0" json:"-"`       // 软删除(毫秒时间戳)
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
