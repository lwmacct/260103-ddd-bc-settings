package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"
	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

const (
	settingsCacheTTL         = 30 * time.Minute
	settingsKeyPrefix        = "settings:"
	settingsUserPrefix       = "user:"
	settingsAdminPrefix      = "admin:"
	settingsPublicPrefix     = "public:"
	settingsCategoryAll      = "_all"
	settingsCategoriesPrefix = "categories:"
	scanBatchSize            = 100
	deletePipelineSize       = 100
)

// settingsCacheService Settings 缓存服务的 Redis 实现。
//
// 缓存 Setting API 的最终响应：
//   - Key 格式：{prefix}settings:user:{userID}:{categoryKey}
//   - Key 格式：{prefix}settings:admin:{categoryKey}
//   - TTL：30 分钟
//   - RedisJSON 原生 JSON 类型存储
//   - 直接序列化 Application DTO（无独立缓存 DTO）
type settingsCacheService struct {
	client    *redis.Client
	keyPrefix string
}

// NewSettingsCacheService 创建 Settings 缓存服务。
func NewSettingsCacheService(client *redis.Client, keyPrefix string) setting.SettingsCacheService {
	return &settingsCacheService{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// =========================================================================
// 用户 Settings 操作
// =========================================================================

// GetUserSettings 获取用户 Settings 缓存。
func (s *settingsCacheService) GetUserSettings(ctx context.Context, userID uint, categoryKey string) ([]setting.SettingsItemDTO, error) {
	key := s.buildUserKey(userID, categoryKey)
	return s.get(ctx, key)
}

// SetUserSettings 设置用户 Settings 缓存。
func (s *settingsCacheService) SetUserSettings(ctx context.Context, userID uint, categoryKey string, settings []setting.SettingsItemDTO) error {
	key := s.buildUserKey(userID, categoryKey)
	return s.set(ctx, key, settings)
}

// DeleteUserSettings 删除用户的指定 category Settings 缓存。
func (s *settingsCacheService) DeleteUserSettings(ctx context.Context, userID uint, categoryKey string) error {
	key := s.buildUserKey(userID, categoryKey)
	return s.client.Del(ctx, key).Err()
}

// DeleteUserSettingsAll 删除用户的所有 Settings 缓存。
func (s *settingsCacheService) DeleteUserSettingsAll(ctx context.Context, userID uint) error {
	pattern := s.keyPrefix + settingsKeyPrefix + settingsUserPrefix + strconv.FormatUint(uint64(userID), 10) + ":*"
	return s.deleteByPattern(ctx, pattern)
}

// =========================================================================
// 管理员 Settings 操作
// =========================================================================

// GetAdminSettings 获取管理员 Settings 缓存。
func (s *settingsCacheService) GetAdminSettings(ctx context.Context, categoryKey string) ([]setting.SettingsItemDTO, error) {
	key := s.buildAdminKey(categoryKey)
	return s.get(ctx, key)
}

// SetAdminSettings 设置管理员 Settings 缓存。
func (s *settingsCacheService) SetAdminSettings(ctx context.Context, categoryKey string, settings []setting.SettingsItemDTO) error {
	key := s.buildAdminKey(categoryKey)
	return s.set(ctx, key, settings)
}

// DeleteAdminSettings 删除管理员的指定 category Settings 缓存。
func (s *settingsCacheService) DeleteAdminSettings(ctx context.Context, categoryKey string) error {
	key := s.buildAdminKey(categoryKey)
	return s.client.Del(ctx, key).Err()
}

// DeleteAdminSettingsAll 删除管理员的所有 Settings 缓存。
func (s *settingsCacheService) DeleteAdminSettingsAll(ctx context.Context) error {
	pattern := s.keyPrefix + settingsKeyPrefix + settingsAdminPrefix + "*"
	return s.deleteByPattern(ctx, pattern)
}

// =========================================================================
// 公开 Settings 操作
// =========================================================================

// GetPublicSettings 获取公开 Settings 缓存。
func (s *settingsCacheService) GetPublicSettings(ctx context.Context, categoryKey string) ([]setting.PublicSettingItemDTO, error) {
	key := s.buildPublicKey(categoryKey)
	return s.getPublic(ctx, key)
}

// SetPublicSettings 设置公开 Settings 缓存。
func (s *settingsCacheService) SetPublicSettings(ctx context.Context, categoryKey string, settings []setting.PublicSettingItemDTO) error {
	key := s.buildPublicKey(categoryKey)
	return s.setPublic(ctx, key, settings)
}

// DeletePublicSettings 删除公开的指定 category Settings 缓存。
func (s *settingsCacheService) DeletePublicSettings(ctx context.Context, categoryKey string) error {
	key := s.buildPublicKey(categoryKey)
	return s.client.Del(ctx, key).Err()
}

// DeletePublicSettingsAll 删除所有公开 Settings 缓存。
func (s *settingsCacheService) DeletePublicSettingsAll(ctx context.Context) error {
	pattern := s.keyPrefix + settingsKeyPrefix + settingsPublicPrefix + "*"
	return s.deleteByPattern(ctx, pattern)
}

// =========================================================================
// 批量失效操作
// =========================================================================

// DeleteByCategoryKey 删除所有用户和管理员的指定 category Settings 缓存。
func (s *settingsCacheService) DeleteByCategoryKey(ctx context.Context, categoryKey string) error {
	if categoryKey == "" {
		categoryKey = settingsCategoryAll
	}

	// 1. 删除管理员的指定 category 缓存
	adminKey := s.buildAdminKey(categoryKey)
	if err := s.client.Del(ctx, adminKey).Err(); err != nil {
		slog.Warn("failed to delete admin settings cache", "categoryKey", categoryKey, "error", err.Error())
	}

	// 2. 删除所有用户的指定 category 缓存
	pattern := s.keyPrefix + settingsKeyPrefix + settingsUserPrefix + "*:" + categoryKey
	if err := s.deleteByPattern(ctx, pattern); err != nil {
		return fmt.Errorf("failed to delete user settings caches: %w", err)
	}

	// 3. 同时删除所有用户的 _all 缓存（因为 _all 包含该 category）
	if categoryKey != settingsCategoryAll {
		allPattern := s.keyPrefix + settingsKeyPrefix + settingsUserPrefix + "*:" + settingsCategoryAll
		if err := s.deleteByPattern(ctx, allPattern); err != nil {
			slog.Warn("failed to delete user _all settings caches", "error", err.Error())
		}

		// 同时删除管理员的 _all 缓存
		adminAllKey := s.buildAdminKey(settingsCategoryAll)
		if err := s.client.Del(ctx, adminAllKey).Err(); err != nil {
			slog.Warn("failed to delete admin _all settings cache", "error", err.Error())
		}
	}

	return nil
}

// DeleteAll 删除所有 Settings 缓存。
func (s *settingsCacheService) DeleteAll(ctx context.Context) error {
	pattern := s.keyPrefix + settingsKeyPrefix + "*"
	return s.deleteByPattern(ctx, pattern)
}

// =========================================================================
// 分类列表缓存操作
// =========================================================================

// GetUserCategories 获取用户分类列表缓存。
func (s *settingsCacheService) GetUserCategories(ctx context.Context) ([]setting.CategoryMetaDTO, error) {
	key := s.buildCategoriesKey("user")
	return s.getCategories(ctx, key)
}

// SetUserCategories 设置用户分类列表缓存。
func (s *settingsCacheService) SetUserCategories(ctx context.Context, categories []setting.CategoryMetaDTO) error {
	key := s.buildCategoriesKey("user")
	return s.setCategories(ctx, key, categories)
}

// DeleteUserCategories 删除用户分类列表缓存。
func (s *settingsCacheService) DeleteUserCategories(ctx context.Context) error {
	key := s.buildCategoriesKey("user")
	return s.client.Del(ctx, key).Err()
}

// =========================================================================
// Category 实体缓存操作（供 Repository 装饰器使用）
// =========================================================================

// GetAllCategories 获取所有 SettingCategory 实体缓存。
// 直接序列化 Domain 实体（实体已有 json tags）。
func (s *settingsCacheService) GetAllCategories(ctx context.Context) ([]*settingdomain.SettingCategory, error) {
	key := s.buildCategoriesKey("all")

	data, err := s.client.JSONGet(ctx, key, "$").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // cache miss
		}
		return nil, fmt.Errorf("redis json get error: %w", err)
	}

	// JSON.GET $ 返回数组包装：[[actual_array]]
	var wrapper [][]*settingdomain.SettingCategory
	if err := json.Unmarshal([]byte(data), &wrapper); err != nil {
		_ = s.client.Del(ctx, key)
		slog.Warn("corrupted category cache, deleted", "key", key, "error", err.Error())
		return nil, nil // corrupted cache treated as miss
	}

	if len(wrapper) == 0 || len(wrapper[0]) == 0 {
		return nil, nil // empty wrapper treated as miss
	}

	return wrapper[0], nil
}

// SetAllCategories 设置所有 SettingCategory 实体缓存。
// 直接序列化 Domain 实体（实体已有 json tags）。
func (s *settingsCacheService) SetAllCategories(ctx context.Context, categories []*settingdomain.SettingCategory) error {
	if len(categories) == 0 {
		return nil
	}

	key := s.buildCategoriesKey("all")

	// 写入缓存（直接序列化实体切片）
	pipe := s.client.Pipeline()
	pipe.JSONSet(ctx, key, "$", categories)
	pipe.Expire(ctx, key, settingsCacheTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set category cache: %w", err)
	}

	return nil
}

// DeleteAllCategories 删除 SettingCategory 实体缓存。
func (s *settingsCacheService) DeleteAllCategories(ctx context.Context) error {
	key := s.buildCategoriesKey("all")
	return s.client.Del(ctx, key).Err()
}

// =========================================================================
// 内部辅助方法
// =========================================================================

// buildUserKey 构建用户 Settings 缓存 key。
func (s *settingsCacheService) buildUserKey(userID uint, categoryKey string) string {
	if categoryKey == "" {
		categoryKey = settingsCategoryAll
	}
	return s.keyPrefix + settingsKeyPrefix + settingsUserPrefix + strconv.FormatUint(uint64(userID), 10) + ":" + categoryKey
}

// buildAdminKey 构建管理员 Settings 缓存 key。
func (s *settingsCacheService) buildAdminKey(categoryKey string) string {
	if categoryKey == "" {
		categoryKey = settingsCategoryAll
	}
	return s.keyPrefix + settingsKeyPrefix + settingsAdminPrefix + categoryKey
}

// buildPublicKey 构建公开 Settings 缓存 key。
func (s *settingsCacheService) buildPublicKey(categoryKey string) string {
	if categoryKey == "" {
		categoryKey = settingsCategoryAll
	}
	return s.keyPrefix + settingsKeyPrefix + settingsPublicPrefix + categoryKey
}

// buildCategoriesKey 构建分类列表缓存 key。
func (s *settingsCacheService) buildCategoriesKey(scope string) string {
	return s.keyPrefix + settingsKeyPrefix + settingsCategoriesPrefix + scope
}

// get 通用获取缓存方法（使用 RedisJSON）。
func (s *settingsCacheService) get(ctx context.Context, key string) ([]setting.SettingsItemDTO, error) {
	// 使用 JSON.GET 命令读取
	data, err := s.client.JSONGet(ctx, key, "$").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // cache miss
		}
		return nil, fmt.Errorf("redis json get error: %w", err)
	}

	// JSON.GET $ 返回数组包装：[actual_data]
	// 需要解包外层数组
	var wrapper [][]setting.SettingsItemDTO
	if err := json.Unmarshal([]byte(data), &wrapper); err != nil {
		// 缓存数据损坏，删除并返回未命中
		_ = s.client.Del(ctx, key)
		slog.Warn("corrupted settings cache, deleted", "key", key, "error", err.Error())
		return nil, nil // corrupted cache treated as miss
	}

	if len(wrapper) == 0 {
		return nil, nil // empty wrapper treated as miss
	}

	// 空切片也视为缓存未命中，避免返回无效数据
	if len(wrapper[0]) == 0 {
		_ = s.client.Del(ctx, key) // 删除无效的空缓存
		return nil, nil
	}

	return wrapper[0], nil
}

// set 通用设置缓存方法（使用 RedisJSON）。
func (s *settingsCacheService) set(ctx context.Context, key string, settings []setting.SettingsItemDTO) error {
	// 使用 Pipeline 执行 JSON.SET + EXPIRE
	pipe := s.client.Pipeline()
	pipe.JSONSet(ctx, key, "$", settings)
	pipe.Expire(ctx, key, settingsCacheTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set settings cache: %w", err)
	}

	return nil
}

// getCategories 获取分类列表缓存（使用 RedisJSON）。
func (s *settingsCacheService) getCategories(ctx context.Context, key string) ([]setting.CategoryMetaDTO, error) {
	data, err := s.client.JSONGet(ctx, key, "$").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // cache miss
		}
		return nil, fmt.Errorf("redis json get error: %w", err)
	}

	// JSON.GET $ 返回数组包装：[actual_data]
	var wrapper [][]setting.CategoryMetaDTO
	if err := json.Unmarshal([]byte(data), &wrapper); err != nil {
		_ = s.client.Del(ctx, key)
		slog.Warn("corrupted categories cache, deleted", "key", key, "error", err.Error())
		return nil, nil
	}

	if len(wrapper) == 0 {
		return nil, nil
	}

	return wrapper[0], nil
}

// setCategories 设置分类列表缓存（使用 RedisJSON）。
func (s *settingsCacheService) setCategories(ctx context.Context, key string, categories []setting.CategoryMetaDTO) error {
	pipe := s.client.Pipeline()
	pipe.JSONSet(ctx, key, "$", categories)
	pipe.Expire(ctx, key, settingsCacheTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set categories cache: %w", err)
	}

	return nil
}

// getPublic 获取公开 Settings 缓存（使用 RedisJSON）。
func (s *settingsCacheService) getPublic(ctx context.Context, key string) ([]setting.PublicSettingItemDTO, error) {
	data, err := s.client.JSONGet(ctx, key, "$").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // cache miss
		}
		return nil, fmt.Errorf("redis json get error: %w", err)
	}

	// JSON.GET $ 返回数组包装：[actual_data]
	var wrapper [][]setting.PublicSettingItemDTO
	if err := json.Unmarshal([]byte(data), &wrapper); err != nil {
		_ = s.client.Del(ctx, key)
		slog.Warn("corrupted public settings cache, deleted", "key", key, "error", err.Error())
		return nil, nil
	}

	if len(wrapper) == 0 || len(wrapper[0]) == 0 {
		return nil, nil
	}

	return wrapper[0], nil
}

// setPublic 设置公开 Settings 缓存（使用 RedisJSON）。
func (s *settingsCacheService) setPublic(ctx context.Context, key string, settings []setting.PublicSettingItemDTO) error {
	pipe := s.client.Pipeline()
	pipe.JSONSet(ctx, key, "$", settings)
	pipe.Expire(ctx, key, settingsCacheTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set public settings cache: %w", err)
	}

	return nil
}

// deleteByPattern 使用 SCAN 批量删除匹配模式的 key。
func (s *settingsCacheService) deleteByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var allKeys []string

	for {
		keys, nextCursor, err := s.client.Scan(ctx, cursor, pattern, scanBatchSize).Result()
		if err != nil {
			return fmt.Errorf("failed to scan keys: %w", err)
		}

		allKeys = append(allKeys, keys...)
		cursor = nextCursor

		if cursor == 0 {
			break
		}
	}

	if len(allKeys) == 0 {
		return nil
	}

	// 使用 Pipeline 批量删除
	pipe := s.client.Pipeline()
	for i := 0; i < len(allKeys); i += deletePipelineSize {
		end := min(i+deletePipelineSize, len(allKeys))
		pipe.Del(ctx, allKeys[i:end]...)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute delete pipeline: %w", err)
	}

	slog.Debug("deleted settings caches", "pattern", pattern, "count", len(allKeys))
	return nil
}

// NotifyAllChanged 通知所有配置已变更（Domain 业务语义接口实现）。
func (s *settingsCacheService) NotifyAllChanged(ctx context.Context) error {
	return s.DeleteAll(ctx)
}

// NotifyCategoryChanged 通知指定分类的配置已变更（Domain 业务语义接口实现）。
func (s *settingsCacheService) NotifyCategoryChanged(ctx context.Context, categoryKey string) error {
	return s.DeleteByCategoryKey(ctx, categoryKey)
}

// 编译时检查：确保 settingsCacheService 实现了必要的接口
var (
	_ setting.SettingsCacheService        = (*settingsCacheService)(nil)
	_ settingdomain.SettingChangeNotifier = (*settingsCacheService)(nil)
)
