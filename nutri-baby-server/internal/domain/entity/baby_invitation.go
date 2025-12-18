package entity

import "gorm.io/plugin/soft_delete"

// BabyInvitation 宝宝邀请记录 (用于微信分享和二维码)
type BabyInvitation struct {
	ID           int64                 `gorm:"primaryKey;column:id" json:"id"`                                           // 雪花ID主键
	BabyID       int64                 `gorm:"column:baby_id;index;not null" json:"babyId"`                              // 宝宝ID (引用Baby.ID)
	UserID       int64                 `gorm:"column:user_id;not null" json:"userId"`                                    // 邀请人用户ID (引用User.ID)
	Token        string                `gorm:"column:token;type:varchar(64);uniqueIndex;not null" json:"token"`          // 临时token(用于验证)
	ShortCode    string                `gorm:"column:short_code;type:varchar(10);uniqueIndex;not null" json:"shortCode"` // 6位短码(用于小程序码scene参数)
	InviteType   string                `gorm:"column:invite_type;type:varchar(20);not null" json:"inviteType"`           // share=分享, qrcode=二维码
	Role         string                `gorm:"column:role;type:varchar(20);not null" json:"role"`                        // admin, editor, viewer
	Relationship string                `gorm:"column:relationship;type:varchar(32)" json:"relationship"`                 // 与宝宝的关系: 爸爸, 妈妈, 爷爷, 奶奶等
	AccessType   string                `gorm:"column:access_type;type:varchar(20);not null" json:"accessType"`           // permanent, temporary
	ExpiresAt    *int64                `gorm:"column:expires_at" json:"expiresAt"`                                       // 协作权限过期时间(毫秒)
	CreatedAt    int64                 `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"`                  // 创建时间(毫秒时间戳)
	DeletedAt    soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;index;default:0" json:"-"`              // 软删除(毫秒时间戳)
}

// TableName 指定表名
func (BabyInvitation) TableName() string {
	return "baby_invitations"
}
