package wechat

import (
	"github.com/redis/go-redis/v9"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/cache"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/config"
)

// Client 微信客户端封装
type Client struct {
	wechat      *wechat.Wechat
	miniProgram *miniprogram.MiniProgram
}

// NewClient 创建微信客户端
func NewClient(cfg *config.Config, redisClient *redis.Client) *Client {
	// 使用应用层注入的 Redis 客户端,而不是创建新的连接
	redisCache := cache.NewRedisCache(redisClient)

	// 创建微信实例
	wc := wechat.NewWechat()

	// 配置小程序
	miniCfg := &miniConfig.Config{
		AppID:     cfg.Wechat.AppID,
		AppSecret: cfg.Wechat.AppSecret,
		Cache:     redisCache,
	}

	// 获取小程序实例
	mini := wc.GetMiniProgram(miniCfg)

	return &Client{
		wechat:      wc,
		miniProgram: mini,
	}
}

// GetMiniProgram 获取小程序实例
func (c *Client) GetMiniProgram() *miniprogram.MiniProgram {
	return c.miniProgram
}
