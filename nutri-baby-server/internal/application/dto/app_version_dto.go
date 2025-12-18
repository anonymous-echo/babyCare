package dto

// AppVersionDTO 应用版本 DTO
type AppVersionDTO struct {
	Version      string `json:"version"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	MinVersion   string `json:"minVersion,omitempty"`
	ForceUpdate  bool   `json:"forceUpdate"`
	ReleaseNotes string `json:"releaseNotes,omitempty"`
	BuildTime    int64  `json:"buildTime"` // 毫秒时间戳
}

// CreateAppVersionRequest 创建版本请求
type CreateAppVersionRequest struct {
	Version      string `json:"version" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	MinVersion   string `json:"minVersion"`
	IsActive     bool   `json:"isActive"`
	ForceUpdate  bool   `json:"forceUpdate"`
	ReleaseNotes string `json:"releaseNotes"`
}

// SetActiveVersionRequest 设置活跃版本请求
type SetActiveVersionRequest struct {
	Version string `json:"version" binding:"required"`
}
