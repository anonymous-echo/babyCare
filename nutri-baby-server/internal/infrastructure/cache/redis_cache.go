package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache 实现 wechat SDK 的缓存接口
// 包装 go-redis 客户端以支持微信 SDK 的缓存操作
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache 创建 Redis 缓存实例
// 使用应用层注入的 redis.Client,确保使用主库而非副本
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
		ctx:    context.Background(),
	}
}

// Get 获取缓存值 - 实现 wechat SDK cache.Cache 接口
// 返回 interface{} 类型，键不存在或出错时返回 nil
func (rc *RedisCache) Get(key string) interface{} {
	val, err := rc.client.Get(rc.ctx, key).Result()
	if err != nil {
		// 键不存在或其他错误，返回 nil
		return nil
	}

	return val
}

// Set 设置缓存值 - 实现 wechat SDK cache.Cache 接口
func (rc *RedisCache) Set(key string, val interface{}, timeout time.Duration) error {
	// 将值转换为字符串（wechat SDK 期望字符串值）
	var value string
	switch v := val.(type) {
	case string:
		value = v
	case []byte:
		value = string(v)
	default:
		// 如果是其他类型，转换为字符串
		value = ""
	}

	return rc.client.Set(rc.ctx, key, value, timeout).Err()
}

// Delete 删除缓存值 - 实现 wechat SDK cache.Cache 接口
func (rc *RedisCache) Delete(key string) error {
	return rc.client.Del(rc.ctx, key).Err()
}

// IsExist 检查缓存值是否存在 - 实现 wechat SDK cache.Cache 接口
func (rc *RedisCache) IsExist(key string) bool {
	exists, err := rc.client.Exists(rc.ctx, key).Result()
	if err != nil {
		return false
	}

	return exists > 0
}
