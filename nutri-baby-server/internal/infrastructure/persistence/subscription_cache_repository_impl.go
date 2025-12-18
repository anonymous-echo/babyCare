package persistence

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
)

// subscriptionCacheRepositoryImpl 实现 repository.SubscriptionCacheRepository 接口
// 使用 Redis Hash 存储用户订阅权限状态
// Redis 结构: subscription:{openid} => {templateType: "allow"|"deny"}
// 注: 不设置过期时间,权限状态永久有效,直到用户在微信中主动拒绝授权时才更新
type subscriptionCacheRepositoryImpl struct {
	client *redis.Client
}

// NewSubscriptionCacheRepository 创建订阅权限缓存仓储实现
func NewSubscriptionCacheRepository(client *redis.Client) repository.SubscriptionCacheRepository {
	return &subscriptionCacheRepositoryImpl{
		client: client,
	}
}

// SetSubscriptionStatus 设置用户对特定模板的订阅权限状态
// 注: 权限状态永久有效,不设置过期时间
func (r *subscriptionCacheRepositoryImpl) SetSubscriptionStatus(
	ctx context.Context,
	openID string,
	templateType string,
	status repository.SubscriptionStatus,
) error {
	if openID == "" || templateType == "" {
		return fmt.Errorf("openID and templateType cannot be empty")
	}

	// 构建缓存键: subscription:{openid}
	hashKey := fmt.Sprintf("subscription:%s", openID)

	// 使用 Hash 结构存储: 每个模板对应一个权限状态(永久有效)
	return r.client.HSet(ctx, hashKey, templateType, string(status)).Err()
}

// GetSubscriptionStatus 获取用户对特定模板的订阅权限状态
func (r *subscriptionCacheRepositoryImpl) GetSubscriptionStatus(
	ctx context.Context,
	openID string,
	templateType string,
) (repository.SubscriptionStatus, bool, error) {
	if openID == "" || templateType == "" {
		return "", false, fmt.Errorf("openID and templateType cannot be empty")
	}

	// 构建缓存键
	hashKey := fmt.Sprintf("subscription:%s", openID)

	// 从 Hash 中获取该模板的权限状态
	status, err := r.client.HGet(ctx, hashKey, templateType).Result()
	if err != nil {
		// 键不存在或其他错误,说明缓存中没有该记录
		if err == redis.Nil {
			return "", false, nil
		}
		return "", false, err
	}

	return repository.SubscriptionStatus(status), true, nil
}

// HasAllowedTemplate 检查用户是否已授权特定模板
func (r *subscriptionCacheRepositoryImpl) HasAllowedTemplate(
	ctx context.Context,
	openID string,
	templateType string,
) (bool, error) {
	status, found, err := r.GetSubscriptionStatus(ctx, openID, templateType)
	if err != nil {
		return false, err
	}

	// 只有找到记录且状态为 allow 时,才返回 true
	return found && status == repository.StatusAllow, nil
}

// HasDeniedTemplate 检查用户是否已拒绝特定模板
func (r *subscriptionCacheRepositoryImpl) HasDeniedTemplate(
	ctx context.Context,
	openID string,
	templateType string,
) (bool, error) {
	status, found, err := r.GetSubscriptionStatus(ctx, openID, templateType)
	if err != nil {
		return false, err
	}

	// 只有找到记录且状态为 deny 时,才返回 true
	return found && status == repository.StatusDeny, nil
}

// GetAllSubscriptions 获取用户的所有订阅权限记录
func (r *subscriptionCacheRepositoryImpl) GetAllSubscriptions(
	ctx context.Context,
	openID string,
) (map[string]repository.SubscriptionStatus, error) {
	if openID == "" {
		return nil, fmt.Errorf("openID cannot be empty")
	}

	hashKey := fmt.Sprintf("subscription:%s", openID)

	// 获取整个 Hash
	subscriptions, err := r.client.HGetAll(ctx, hashKey).Result()
	if err != nil {
		return nil, err
	}

	// 转换为返回格式
	result := make(map[string]repository.SubscriptionStatus)
	for templateType, status := range subscriptions {
		result[templateType] = repository.SubscriptionStatus(status)
	}

	return result, nil
}

// ClearSubscriptions 清除用户的所有订阅权限缓存
func (r *subscriptionCacheRepositoryImpl) ClearSubscriptions(ctx context.Context, openID string) error {
	if openID == "" {
		return fmt.Errorf("openID cannot be empty")
	}

	hashKey := fmt.Sprintf("subscription:%s", openID)
	return r.client.Del(ctx, hashKey).Err()
}

// ClearSubscriptionStatus 清除用户对特定模板的订阅权限状态
func (r *subscriptionCacheRepositoryImpl) ClearSubscriptionStatus(
	ctx context.Context,
	openID string,
	templateType string,
) error {
	if openID == "" || templateType == "" {
		return fmt.Errorf("openID and templateType cannot be empty")
	}

	hashKey := fmt.Sprintf("subscription:%s", openID)
	return r.client.HDel(ctx, hashKey, templateType).Err()
}
