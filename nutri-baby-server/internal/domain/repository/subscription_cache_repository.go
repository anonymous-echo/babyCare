package repository

import (
	"context"
)

// SubscriptionStatus 订阅权限状态
type SubscriptionStatus string

const (
	// StatusAllow 用户已授权该模板
	StatusAllow SubscriptionStatus = "accept"
	// StatusDeny 用户已拒绝该模板
	StatusDeny SubscriptionStatus = "reject"
)

// SubscriptionCacheRepository 用户订阅权限缓存仓储接口
//
// 定义用户对各消息模板的订阅权限状态存储和查询接口
// 实现: 基础设施层的 Redis Hash 存储
//
// 权限状态说明:
//   - allow: 用户已授权该模板,可以发送提醒
//   - deny: 用户已拒绝该模板,需要重新授权
//   - 未记录: 用户未做过选择,需要显示授权弹窗
type SubscriptionCacheRepository interface {
	// SetSubscriptionStatus 设置用户对特定模板的订阅权限状态
	//
	// 当用户在授权弹窗中做出选择后,微信会回调此接口,记录用户的授权状态
	// 注: "总是保持以上选择"由微信官方实现,后端无需处理
	//
	// 参数:
	//   - ctx: 上下文
	//   - openID: 用户OpenID
	//   - templateType: 模板类型(如 "breast_feeding_reminder")
	//   - status: 权限状态 (StatusAllow 或 StatusDeny)
	SetSubscriptionStatus(
		ctx context.Context,
		openID string,
		templateType string,
		status SubscriptionStatus,
	) error

	// GetSubscriptionStatus 获取用户对特定模板的订阅权限状态
	//
	// 返回值:
	//   - status: 权限状态 (StatusAllow/StatusDeny),如果未缓存返回空字符串
	//   - found: 是否在缓存中找到该记录(true=使用缓存值,false=需要再次询问用户)
	//   - err: 错误信息
	GetSubscriptionStatus(
		ctx context.Context,
		openID string,
		templateType string,
	) (SubscriptionStatus, bool, error)

	// HasAllowedTemplate 检查用户是否已授权特定模板
	HasAllowedTemplate(
		ctx context.Context,
		openID string,
		templateType string,
	) (bool, error)

	// HasDeniedTemplate 检查用户是否已拒绝特定模板
	HasDeniedTemplate(
		ctx context.Context,
		openID string,
		templateType string,
	) (bool, error)

	// GetAllSubscriptions 获取用户的所有订阅权限记录
	//
	// 返回用户缓存中所有模板的权限状态
	// 返回格式: map[templateType]status, 例如: {"breast_feeding_reminder": "allow", "bottle_feeding_reminder": "deny"}
	GetAllSubscriptions(
		ctx context.Context,
		openID string,
	) (map[string]SubscriptionStatus, error)

	// ClearSubscriptions 清除用户的所有订阅权限缓存
	// 用户可在设置中选择"重置订阅权限"时调用
	ClearSubscriptions(ctx context.Context, openID string) error

	// ClearSubscriptionStatus 清除用户对特定模板的订阅权限状态
	ClearSubscriptionStatus(
		ctx context.Context,
		openID string,
		templateType string,
	) error
}
