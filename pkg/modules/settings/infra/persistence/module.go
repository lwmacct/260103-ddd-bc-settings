package persistence

import (
	"go.uber.org/fx"
	"gorm.io/gorm"

	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// RepositoryModule 提供 Settings 模块的所有仓储实现。
//
// 装饰缓存层的仓储：
//   - Setting: 缓存查询 + 缓存命令（失效缓存）
//   - SettingCategory: 缓存查询 + 缓存命令（失效缓存）
var RepositoryModule = fx.Module("settings.repository",
	fx.Provide(
		// 带缓存装饰的仓储
		newSettingRepositoriesWithCache,
	),
)

// newSettingRepositoriesWithCache 创建带缓存装饰的仓储。
//
// 组合原始仓储和缓存服务，提供缓存失效策略：
//   - Query 操作：先查缓存，未命中则查数据库
//   - Command 操作：执行写操作后异步失效相关缓存
//
// 参数 cacheInvalidator 需要实现 settingdomain.CacheInvalidator 接口，
// 通常由 Infrastructure 层的 SettingsCacheService 实现。
func newSettingRepositoriesWithCache(
	db *gorm.DB,
	cacheInvalidator settingdomain.CacheInvalidator,
) SettingRepositories {
	rawRepos := NewSettingRepositories(db)

	// Setting 缓存装饰（只有 Command 需要失效缓存，Query 缓存在 Application 层）
	cachedCommand := NewCachedSettingCommandRepository(rawRepos.Command, cacheInvalidator)

	// SettingCategory Query 不做缓存装饰（Application 层有缓存）
	// SettingCategory Command 不需要缓存失效

	return SettingRepositories{
		Command:         cachedCommand,
		Query:           rawRepos.Query,
		CategoryCommand: NewSettingCategoryCommandRepository(db),
		CategoryQuery:   NewSettingCategoryQueryRepository(db),
	}
}
