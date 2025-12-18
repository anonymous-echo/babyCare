package cache

import (
	"context"
	"sync"
	"time"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
)

// AnalysisDataCache 分析数据缓存
type AnalysisDataCache struct {
	mu           sync.RWMutex
	babyCache    map[int64]*babyCacheEntry
	cacheTTL     time.Duration
	maxCacheSize int
}

// babyCacheEntry 宝宝数据缓存条目
type babyCacheEntry struct {
	babyInfo      *entity.Baby
	feedingData   []entity.FeedingRecord
	sleepData     []entity.SleepRecord
	growthData    []entity.GrowthRecord
	diaperData    []entity.DiaperRecord
	expiresAt     time.Time
	lastAccessed  time.Time
}

// NewAnalysisDataCache 创建数据缓存
func NewAnalysisDataCache(ttl time.Duration, maxSize int) *AnalysisDataCache {
	cache := &AnalysisDataCache{
		babyCache:    make(map[int64]*babyCacheEntry),
		cacheTTL:     ttl,
		maxCacheSize: maxSize,
	}
	
	// 启动定期清理过期缓存的goroutine
	go cache.cleanupExpiredEntries()
	
	return cache
}

// GetBabyInfo 获取宝宝信息（带缓存）
func (c *AnalysisDataCache) GetBabyInfo(ctx context.Context, babyID int64, fetcher func(context.Context, int64) (*entity.Baby, error)) (*entity.Baby, error) {
	c.mu.RLock()
	entry, exists := c.babyCache[babyID]
	if exists && entry.babyInfo != nil && time.Now().Before(entry.expiresAt) {
		entry.lastAccessed = time.Now()
		c.mu.RUnlock()
		return entry.babyInfo, nil
	}
	c.mu.RUnlock()

	// 缓存未命中或已过期，重新获取
	baby, err := fetcher(ctx, babyID)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	
	if entry == nil {
		entry = &babyCacheEntry{
			expiresAt:    time.Now().Add(c.cacheTTL),
			lastAccessed: time.Now(),
		}
		c.babyCache[babyID] = entry
	}
	
	entry.babyInfo = baby
	entry.expiresAt = time.Now().Add(c.cacheTTL)
	
	return baby, nil
}

// GetFeedingData 获取喂养数据（带缓存）
func (c *AnalysisDataCache) GetFeedingData(ctx context.Context, babyID int64, startDate, endDate time.Time, fetcher func(context.Context, int64, time.Time, time.Time) ([]entity.FeedingRecord, error)) ([]entity.FeedingRecord, error) {
	c.mu.RLock()
	entry, exists := c.babyCache[babyID]
	if exists && entry.feedingData != nil && time.Now().Before(entry.expiresAt) {
		entry.lastAccessed = time.Now()
		c.mu.RUnlock()
		return entry.feedingData, nil
	}
	c.mu.RUnlock()

	// 缓存未命中或已过期，重新获取
	data, err := fetcher(ctx, babyID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	
	if entry == nil {
		entry = &babyCacheEntry{
			expiresAt:    time.Now().Add(c.cacheTTL),
			lastAccessed: time.Now(),
		}
		c.babyCache[babyID] = entry
	}
	
	entry.feedingData = data
	entry.expiresAt = time.Now().Add(c.cacheTTL)
	
	return data, nil
}

// GetSleepData 获取睡眠数据（带缓存）
func (c *AnalysisDataCache) GetSleepData(ctx context.Context, babyID int64, startDate, endDate time.Time, fetcher func(context.Context, int64, time.Time, time.Time) ([]entity.SleepRecord, error)) ([]entity.SleepRecord, error) {
	c.mu.RLock()
	entry, exists := c.babyCache[babyID]
	if exists && entry.sleepData != nil && time.Now().Before(entry.expiresAt) {
		entry.lastAccessed = time.Now()
		c.mu.RUnlock()
		return entry.sleepData, nil
	}
	c.mu.RUnlock()

	data, err := fetcher(ctx, babyID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	
	if entry == nil {
		entry = &babyCacheEntry{
			expiresAt:    time.Now().Add(c.cacheTTL),
			lastAccessed: time.Now(),
		}
		c.babyCache[babyID] = entry
	}
	
	entry.sleepData = data
	entry.expiresAt = time.Now().Add(c.cacheTTL)
	
	return data, nil
}

// GetGrowthData 获取成长数据（带缓存）
func (c *AnalysisDataCache) GetGrowthData(ctx context.Context, babyID int64, startDate, endDate time.Time, fetcher func(context.Context, int64, time.Time, time.Time) ([]entity.GrowthRecord, error)) ([]entity.GrowthRecord, error) {
	c.mu.RLock()
	entry, exists := c.babyCache[babyID]
	if exists && entry.growthData != nil && time.Now().Before(entry.expiresAt) {
		entry.lastAccessed = time.Now()
		c.mu.RUnlock()
		return entry.growthData, nil
	}
	c.mu.RUnlock()

	data, err := fetcher(ctx, babyID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	
	if entry == nil {
		entry = &babyCacheEntry{
			expiresAt:    time.Now().Add(c.cacheTTL),
			lastAccessed: time.Now(),
		}
		c.babyCache[babyID] = entry
	}
	
	entry.growthData = data
	entry.expiresAt = time.Now().Add(c.cacheTTL)
	
	return data, nil
}

// InvalidateCache 使缓存失效
func (c *AnalysisDataCache) InvalidateCache(babyID int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.babyCache, babyID)
}

// InvalidateAll 清空所有缓存
func (c *AnalysisDataCache) InvalidateAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.babyCache = make(map[int64]*babyCacheEntry)
}

// cleanupExpiredEntries 定期清理过期缓存
func (c *AnalysisDataCache) cleanupExpiredEntries() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for babyID, entry := range c.babyCache {
			if now.After(entry.expiresAt) {
				delete(c.babyCache, babyID)
			}
		}
		
		// 如果缓存超过最大大小，清理最久未访问的条目
		if len(c.babyCache) > c.maxCacheSize {
			c.evictLRU()
		}
		c.mu.Unlock()
	}
}

// evictLRU 清理最久未访问的缓存条目
func (c *AnalysisDataCache) evictLRU() {
	var oldestID int64
	var oldestTime time.Time
	first := true
	
	for babyID, entry := range c.babyCache {
		if first || entry.lastAccessed.Before(oldestTime) {
			oldestID = babyID
			oldestTime = entry.lastAccessed
			first = false
		}
	}
	
	if !first {
		delete(c.babyCache, oldestID)
	}
}

// GetCacheStats 获取缓存统计信息
func (c *AnalysisDataCache) GetCacheStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return map[string]interface{}{
		"total_entries": len(c.babyCache),
		"max_size":      c.maxCacheSize,
		"ttl_seconds":   c.cacheTTL.Seconds(),
	}
}
