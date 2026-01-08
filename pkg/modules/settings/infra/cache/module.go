package cache

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"
	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/config"
	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// CacheResult 使用 fx.Out 同时提供两个接口。
type CacheResult struct {
	fx.Out

	// Application 层使用的完整缓存服务
	SettingsCacheService setting.SettingsCacheService

	// Infrastructure 层使用的配置变更通知接口（业务语义）
	SettingChangeNotifier settingdomain.SettingChangeNotifier
}

// CacheModule 提供 Settings 模块的缓存服务。
//
// 依赖 Settings 模块的 config.Config（由 internal/container 提供）。
var CacheModule = fx.Module("settings.cache",
	fx.Provide(
		NewSettingsCacheServiceAs,
	),
)

// NewSettingsCacheServiceAs 创建缓存服务并同时提供两个接口。
func NewSettingsCacheServiceAs(client *redis.Client, settingsCfg *config.Config) CacheResult {
	service := &settingsCacheService{
		client:    client,
		keyPrefix: settingsCfg.Redis.KeyPrefix,
	}

	return CacheResult{
		SettingsCacheService:  service,
		SettingChangeNotifier: service,
	}
}

// 编译时检查：确保实现了所有必需的接口
var (
	_ setting.SettingsCacheService        = (*settingsCacheService)(nil)
	_ settingdomain.SettingChangeNotifier = (*settingsCacheService)(nil)
)
