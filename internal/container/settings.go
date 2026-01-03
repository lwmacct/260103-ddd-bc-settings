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
func newSettingsConfig(cfg *config.Config) *settingsconfig.Config {
	return &settingsconfig.Config{
		Redis: settingsconfig.Redis{
			KeyPrefix: cfg.Data.RedisKeyPrefix,
		},
	}
}
