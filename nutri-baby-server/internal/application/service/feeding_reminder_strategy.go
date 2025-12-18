package service

import (
	"fmt"
	"time"

	"errors"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/config"
)

// FeedingReminderStrategy 喂养提醒策略接口
type FeedingReminderStrategy interface {
	// GetTemplateType 获取微信订阅消息模板类型
	GetTemplateType() string

	// GetTemplateID 获取微信订阅消息模板ID
	GetTemplateID() string

	// BuildMessageData 构建消息数据
	BuildMessageData(record *entity.FeedingRecord, lastFeedingTime time.Time, hoursSinceLastFeeding float64) map[string]any

	// CanHandle 判断是否能处理该类型的喂养记录
	CanHandle(record *entity.FeedingRecord) bool
}

// BreastFeedingReminderStrategy 母乳喂养提醒策略
type BreastFeedingReminderStrategy struct {
	config *config.Config
}

func NewBreastFeedingReminderStrategy(cfg *config.Config) *BreastFeedingReminderStrategy {
	return &BreastFeedingReminderStrategy{
		config: cfg,
	}
}

func (s *BreastFeedingReminderStrategy) GetTemplateID() string {
	return s.config.Wechat.SubscribeTemplates["breast_feeding_reminder"]
}

func (s *BreastFeedingReminderStrategy) GetTemplateType() string {
	return "breast_feeding_reminder"
}

func (s *BreastFeedingReminderStrategy) BuildMessageData(record *entity.FeedingRecord, lastFeedingTime time.Time, hoursSinceLastFeeding float64) map[string]any {
	// 微信订阅消息模板字段: time1(上次时间), thing2(距离上次), character_string3(喂养量), phrase4(喂养类型), thing5(温馨提示)

	// 获取喂养侧
	side := "母乳"
	if sideVal, ok := record.Detail["breastSide"].(string); ok {
		switch sideVal {
		case "left":
			side = "左侧"
		case "right":
			side = "右侧"
		case "both":
			side = "两侧"
		}
	}

	return map[string]any{
		"time1":   lastFeedingTime.Format(time.DateTime),  // 上次时间
		"thing2":  formatTimeSince(hoursSinceLastFeeding), // 距离上次
		"phrase3": side,                                   // 喂养位置
		"thing4":  "该喂奶啦，注意观察宝宝的饥饿信号",                     // 温馨提示
	}
}

func (s *BreastFeedingReminderStrategy) CanHandle(record *entity.FeedingRecord) bool {
	return record.FeedingType == entity.FeedingTypeBreast
}

// BottleFeedingReminderStrategy 奶瓶喂养提醒策略
type BottleFeedingReminderStrategy struct {
	config *config.Config
}

func NewBottleFeedingReminderStrategy(cfg *config.Config) *BottleFeedingReminderStrategy {
	return &BottleFeedingReminderStrategy{
		config: cfg,
	}
}

func (s *BottleFeedingReminderStrategy) GetTemplateID() string {
	return s.config.Wechat.SubscribeTemplates["bottle_feeding_reminder"]
}

func (s *BottleFeedingReminderStrategy) GetTemplateType() string {
	return "bottle_feeding_reminder"
}

func (s *BottleFeedingReminderStrategy) BuildMessageData(record *entity.FeedingRecord, lastFeedingTime time.Time, hoursSinceLastFeeding float64) map[string]interface{} {
	// 微信订阅消息模板字段: time1(上次时间), thing2(距离上次), character_string3(喂养量), phrase4(喂养类型), thing5(温馨提示)

	// 获取奶瓶类型
	bottleType := "配方奶"
	if bottleTypeVal, ok := record.Detail["bottleType"].(string); ok {
		if bottleTypeVal == "breast-milk" {
			bottleType = "母乳"
		}
	}

	// 获取奶量
	amount := ""
	if amountVal, ok := record.Detail["amount"].(float64); ok {
		amount = fmt.Sprintf("%.0fml", amountVal)
	}

	return map[string]interface{}{
		"time1":             lastFeedingTime.Format(time.DateTime),  // 上次时间
		"thing2":            formatTimeSince(hoursSinceLastFeeding), // 距离上次
		"character_string3": amount,                                 // 喂养量
		"phrase4":           bottleType,                             // 喂养类型
		"thing5":            "该喂奶啦，记得准备好奶瓶哦",                        // 温馨提示
	}
}

func (s *BottleFeedingReminderStrategy) CanHandle(record *entity.FeedingRecord) bool {
	return record.FeedingType == entity.FeedingTypeBottle
}

// FoodFeedingReminderStrategy 辅食喂养提醒策略
type FoodFeedingReminderStrategy struct {
	config *config.Config
}

func NewFoodFeedingReminderStrategy(cfg *config.Config) *FoodFeedingReminderStrategy {
	return &FoodFeedingReminderStrategy{
		config: cfg,
	}
}

func (s *FoodFeedingReminderStrategy) GetTemplateType() string {
	return "food_feeding_reminder"
}

func (s *FoodFeedingReminderStrategy) GetTemplateID() string {
	return s.config.Wechat.SubscribeTemplates["food_feeding_reminder"]
}

func (s *FoodFeedingReminderStrategy) BuildMessageData(record *entity.FeedingRecord, lastFeedingTime time.Time, hoursSinceLastFeeding float64) map[string]interface{} {
	// 微信订阅消息模板字段: time1(上次时间), thing2(距离上次), character_string3(食物名称), phrase4(喂养类型), thing5(温馨提示)

	// 获取辅食名称
	foodName := "辅食"
	if foodNameVal, ok := record.Detail["foodName"].(string); ok && foodNameVal != "" {
		foodName = foodNameVal
	}

	return map[string]interface{}{
		"time1":             lastFeedingTime.Format(time.DateTime),  // 上次时间
		"thing2":            formatTimeSince(hoursSinceLastFeeding), // 距离上次
		"character_string3": foodName,                               // 食物名称
		"phrase4":           "辅食",                                   // 喂养类型
		"thing5":            "该给宝宝准备辅食啦，注意观察过敏反应",                   // 温馨提示
	}
}

func (s *FoodFeedingReminderStrategy) CanHandle(record *entity.FeedingRecord) bool {
	return record.FeedingType == entity.FeedingTypeFood
}

// FeedingReminderStrategyFactory 喂养提醒策略工厂
type FeedingReminderStrategyFactory struct {
	strategies map[string]FeedingReminderStrategy
}

// NewFeedingReminderStrategyFactory 创建策略工厂
func NewFeedingReminderStrategyFactory(cfg *config.Config) *FeedingReminderStrategyFactory {
	return &FeedingReminderStrategyFactory{
		strategies: map[string]FeedingReminderStrategy{
			entity.FeedingTypeBottle: NewBottleFeedingReminderStrategy(cfg),
			entity.FeedingTypeFood:   NewFoodFeedingReminderStrategy(cfg),
			entity.FeedingTypeBreast: NewBreastFeedingReminderStrategy(cfg),
		},
	}
}

// GetStrategy 根据喂养记录获取对应的策略
func (f *FeedingReminderStrategyFactory) GetStrategy(record *entity.FeedingRecord) (FeedingReminderStrategy, error) {
	strategy, ok := f.strategies[record.FeedingType]
	if ok {
		return strategy, nil
	}
	return nil, errors.New("unsupported feeding type")
}

// formatTimeSince 格式化距离上次的时间
func formatTimeSince(hours float64) string {
	if hours < 1 {
		minutes := int(hours * 60)
		return fmt.Sprintf("约%d分钟", minutes)
	}

	h := int(hours)
	if h == 1 {
		return "约1小时"
	}

	return fmt.Sprintf("约%d小时", h)
}
