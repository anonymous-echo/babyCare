package dto

// CreateBabyRequest 创建宝宝请求 (去家庭化架构)
type CreateBabyRequest struct {
	Name                  string `json:"name" binding:"required"`
	Nickname              string `json:"nickname"`
	Gender                string `json:"gender" binding:"required,oneof=male female"`
	BirthDate             string `json:"birthDate" binding:"required"` // YYYY-MM-DD
	AvatarURL             string `json:"avatarUrl"`
	CopyCollaboratorsFrom string `json:"copyCollaboratorsFrom"` // 可选:复制协作者的源宝宝ID
}

// UpdateBabyRequest 更新宝宝请求
type UpdateBabyRequest struct {
	Name      string `json:"name"`
	Nickname  string `json:"nickname"`
	Gender    string `json:"gender" binding:"omitempty,oneof=male female"`
	BirthDate string `json:"birthDate"` // YYYY-MM-DD
	AvatarURL string `json:"avatarUrl"`
	Height    int    `json:"height"` // cm
	Weight    int    `json:"weight"` // g
}

// BabyDTO 宝宝DTO (去家庭化架构)
type BabyDTO struct {
	BabyID     string `json:"babyId"`
	Name       string `json:"name"`
	Nickname   string `json:"nickname"`
	Gender     string `json:"gender"`
	BirthDate  string `json:"birthDate"`
	AvatarURL  string `json:"avatarUrl"`
	CreatorID  string `json:"creatorId"` // 创建者 openid
	Height     int    `json:"height"`
	Weight     int    `json:"weight"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

// FamilyMemberDTO 亲友团成员DTO (原 CollaboratorDTO)
type CollaboratorDTO struct {
	OpenID       string `json:"openid"`
	NickName     string `json:"nickName"`
	AvatarURL    string `json:"avatarUrl"`
	Role         string `json:"role"`         // admin, editor, viewer
	Relationship string `json:"relationship"` // 与宝宝的关系: 爸爸, 妈妈, 爷爷, 奶奶, 外公, 外婆等
	AccessType   string `json:"accessType"`   // permanent, temporary
	ExpiresAt    *int64 `json:"expiresAt"`    // 临时权限过期时间
	JoinTime     int64  `json:"joinTime"`
}

// InviteFamilyMemberRequest 邀请亲友团成员请求 (微信分享/二维码)
type InviteCollaboratorRequest struct {
	InviteType   string `json:"inviteType" binding:"required,oneof=share qrcode"` // share=微信分享, qrcode=二维码
	Role         string `json:"role" binding:"required,oneof=admin editor viewer"`
	Relationship string `json:"relationship"` // 与宝宝的关系: 爸爸, 妈妈, 爷爷, 奶奶, 外公, 外婆等
	AccessType   string `json:"accessType" binding:"required,oneof=permanent temporary"`
	ExpiresAt    *int64 `json:"expiresAt"` // 仅当 accessType=temporary 时需要
}

// BabyInvitationDTO 宝宝邀请信息DTO
type BabyInvitationDTO struct {
	BabyID       string        `json:"babyId"`       // 宝宝ID
	Name         string        `json:"name"`         // 宝宝名称
	InviterName  string        `json:"inviterName"`  // 邀请人名称
	Role         string        `json:"role"`         // 角色
	ShortCode    string        `json:"shortCode"`    // 6位短码(用于小程序码scene参数)
	QRCodeParams *QRCodeParams `json:"qrcodeParams"` // 二维码参数
	ExpiresAt    *int64        `json:"expiresAt"`    // 协作者权限过期时间(临时权限)
}

// ShareParams 微信小程序分享参数
type ShareParams struct {
	Title    string `json:"title"`    // 分享标题: "邀请你一起记录{宝宝名}的成长"
	Path     string `json:"path"`     // 小程序路径: pages/baby/join/join?babyId=xxx&token=xxx
	ImageURL string `json:"imageUrl"` // 分享图片(宝宝头像或默认图)
}

// QRCodeParams 二维码参数
type QRCodeParams struct {
	Scene     string `json:"scene"`     // 二维码场景值: babyId=xxx&token=xxx
	Page      string `json:"page"`      // 小程序页面路径
	QRCodeURL string `json:"qrcodeUrl"` // 二维码图片URL(前端生成或后端生成)
}

// JoinBabyRequest 加入宝宝协作请求
type JoinBabyRequest struct {
	BabyID string `json:"babyId" binding:"required"` // 宝宝ID
	Token  string `json:"token" binding:"required"`  // 临时token(验证邀请有效性)
}

// UpdateFamilyMemberRequest 更新亲友团成员请求
type UpdateFamilyMemberRequest struct {
	Role         string `json:"role" binding:"omitempty,oneof=admin editor viewer"` // 角色
	Relationship string `json:"relationship"`                                       // 与宝宝的关系
}

// InvitationDetailDTO 邀请详情DTO (用于通过短码查询)
type InvitationDetailDTO struct {
	BabyID      string `json:"babyId"`      // 宝宝ID
	BabyName    string `json:"babyName"`    // 宝宝名称
	BabyAvatar  string `json:"babyAvatar"`  // 宝宝头像
	InviterName string `json:"inviterName"` // 邀请人名称
	Role        string `json:"role"`        // 角色
	AccessType  string `json:"accessType"`  // 访问类型
	ExpiresAt   *int64 `json:"expiresAt"`   // 权限过期时间(临时权限)
	Token       string `json:"token"`       // Token(用于加入)
}
