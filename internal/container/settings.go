package container

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-settings/internal/config"
	settingsconfig "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/config"
)

// SettingsConfigModule 提供 Settings 模块的配置。
//
// 从全局 internal/config.Config 提取 Settings 模块所需的配置项，
// 实现 Settings 模块的配置自治（不直接依赖 internal/config）。
var SettingsConfigModule = fx.Module("settings.config",
	fx.Provide(newSettingsConfig),
)

// newSettingsConfig 从全局配置构造 Settings 模块配置。
//
// 组合基础前缀和模块标识，形成完整的缓存键前缀：
//   - 基础前缀：cfg.Data.RedisKeyPrefix（如 "app:"）
//   - 模块标识："settings:"
//   - 组合结果："app:settings:"
func newSettingsConfig(cfg *config.Config) *settingsconfig.Config {
	return &settingsconfig.Config{
		RedisKeyPrefix: cfg.Data.RedisKeyPrefix + "settings:",
		API:            cfg.Settings.API,
	}
}
