package entity

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// Baby 宝宝实体 (去家庭化架构)
type Baby struct {
	ID          int64                 `gorm:"primaryKey;column:id" json:"id"`                                // 雪花ID主键
	Name        string                `gorm:"column:name;type:varchar(64)" json:"name"`                      // 姓名
	Nickname    string                `gorm:"column:nickname;type:varchar(64)" json:"nickname"`              // 昵称
	BirthDate   string                `gorm:"column:birth_date;type:varchar(10)" json:"birthDate"`           // 出生日期 YYYY-MM-DD
	Gender      string                `gorm:"column:gender;type:varchar(16)" json:"gender"`                  // 性别 male, female
	AvatarURL   string                `gorm:"column:avatar_url;type:varchar(512)" json:"avatarUrl"`          // 头像URL
	Height      float64               `gorm:"column:height;type:decimal(10,2)" json:"height"`                // 身高 cm
	Weight      float64               `gorm:"column:weight;type:decimal(10,2)" json:"weight"`                // 体重 kg
	UserID      int64                 `gorm:"column:user_id;index" json:"userId"`                            // 创建者用户ID (引用User.ID)
	FamilyGroup string                `gorm:"column:family_group;type:varchar(64);index" json:"familyGroup"` // 可选的家庭分组名称
	CreatedAt   int64                 `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"`       // 创建时间(毫秒时间戳)
	UpdatedAt   int64                 `gorm:"column:updated_at;autoUpdateTime:milli" json:"updatedAt"`       // 更新时间(毫秒时间戳)
	DeletedAt   soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;index;default:0" json:"-"`   // 软删除(毫秒时间戳)

	// 关联
	Collaborators []*BabyCollaborator `gorm:"foreignKey:BabyID;references:ID" json:"collaborators,omitempty"`
}

// TableName 指定表名
func (Baby) TableName() string {
	return "babies"
}

// BabyFamilyMember 宝宝亲友团成员实体 (原 BabyCollaborator)
type BabyCollaborator struct {
	ID           int64                 `gorm:"primaryKey;column:id" json:"id"`                                            // 雪花ID主键
	BabyID       int64                 `gorm:"column:baby_id;index;uniqueIndex:idx_baby_user" json:"babyId"`              // 宝宝ID (引用Baby.ID)
	UserID       int64                 `gorm:"column:user_id;index;uniqueIndex:idx_baby_user" json:"userId"`              // 用户ID (引用User.ID)
	Role         string                `gorm:"column:role;type:varchar(16)" json:"role"`                                  // 角色 admin, editor, viewer
	Relationship string                `gorm:"column:relationship;type:varchar(32)" json:"relationship"`                  // 与宝宝的关系: 爸爸, 妈妈, 爷爷, 奶奶, 外公, 外婆, 叔叔, 阿姨等
	AccessType   string                `gorm:"column:access_type;type:varchar(16);default:'permanent'" json:"accessType"` // 访问类型 permanent, temporary
	ExpiresAt    *int64                `gorm:"column:expires_at" json:"expiresAt"`                                        // 临时权限过期时间(毫秒时间戳)
	CreatedAt    int64                 `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"`                   // 创建时间(毫秒时间戳)
	UpdatedAt    int64                 `gorm:"column:updated_at;autoUpdateTime:milli" json:"updatedAt"`                   // 更新时间(毫秒时间戳)
	DeletedAt    soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;index;default:0" json:"-"`               // 软删除(毫秒时间戳)

	// 关联
	User *User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Baby *Baby `gorm:"foreignKey:BabyID;references:ID" json:"baby,omitempty"`
}

// TableName 指定表名
func (BabyCollaborator) TableName() string {
	return "baby_collaborators"
}

// IsExpired 检查临时权限是否过期
func (bc *BabyCollaborator) IsExpired() bool {
	if bc.AccessType != "temporary" || bc.ExpiresAt == nil {
		return false
	}
	return time.Now().UnixMilli() > *bc.ExpiresAt
}

// IsAdmin 检查是否为管理员
func (bc *BabyCollaborator) IsAdmin() bool {
	return bc.Role == "admin"
}

// CanEdit 检查是否有编辑权限
func (bc *BabyCollaborator) CanEdit() bool {
	return bc.Role == "admin" || bc.Role == "editor"
}
