package setting

import (
	"context"

	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// SettingsCacheService Settings 缓存服务接口。
//
// 缓存 Settings API 的最终响应：
//   - Key 格式：{prefix}settings:user:{userID}:{categoryKey}
//   - Key 格式：{prefix}settings:admin:{categoryKey}
//   - TTL：30 分钟
//   - RedisJSON 原生 JSON 类型存储
//   - 直接序列化 Application DTO（无独立缓存 DTO）
type SettingsCacheService interface {
	// =========================================================================
	// 用户 Settings 操作
	// =========================================================================

	// GetUserSettings 获取用户 Settings 缓存。
	GetUserSettings(ctx context.Context, userID uint, categoryKey string) ([]SettingsCategoryDTO, error)

	// SetUserSettings 设置用户 Settings 缓存。
	SetUserSettings(ctx context.Context, userID uint, categoryKey string, settings []SettingsCategoryDTO) error

	// DeleteUserSettings 删除用户的指定 category Settings 缓存。
	DeleteUserSettings(ctx context.Context, userID uint, categoryKey string) error

	// DeleteUserSettingsAll 删除用户的所有 Settings 缓存。
	DeleteUserSettingsAll(ctx context.Context, userID uint) error

	// =========================================================================
	// 管理员 Settings 操作
	// =========================================================================

	// GetAdminSettings 获取管理员 Settings 缓存。
	GetAdminSettings(ctx context.Context, categoryKey string) ([]SettingsCategoryDTO, error)

	// SetAdminSettings 设置管理员 Settings 缓存。
	SetAdminSettings(ctx context.Context, categoryKey string, settings []SettingsCategoryDTO) error

	// DeleteAdminSettings 删除管理员的指定 category Settings 缓存。
	DeleteAdminSettings(ctx context.Context, categoryKey string) error

	// DeleteAdminSettingsAll 删除管理员的所有 Settings 缓存。
	DeleteAdminSettingsAll(ctx context.Context) error

	// =========================================================================
	// 批量失效操作
	// =========================================================================

	// DeleteByCategoryKey 删除所有用户和管理员的指定 category Settings 缓存。
	DeleteByCategoryKey(ctx context.Context, categoryKey string) error

	// DeleteAll 删除所有 Settings 缓存。
	DeleteAll(ctx context.Context) error

	// =========================================================================
	// Category 实体缓存操作（供 Repository 装饰器使用）
	// =========================================================================

	// GetAllCategories 获取所有 SettingCategory 实体缓存。
	GetAllCategories(ctx context.Context) ([]*settingdomain.SettingCategory, error)

	// SetAllCategories 设置所有 SettingCategory 实体缓存。
	SetAllCategories(ctx context.Context, categories []*settingdomain.SettingCategory) error

	// DeleteAllCategories 删除 SettingCategory 实体缓存。
	DeleteAllCategories(ctx context.Context) error
}
