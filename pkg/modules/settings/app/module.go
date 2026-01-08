// Package app 提供 Settings 模块的应用层聚合。
//
// 本包聚合所有子模块的 UseCase，并提供类型别名以简化依赖注入。
//
// # 架构模式
//
// 采用分布式模块架构：
//   - 子模块（setting/）独立管理依赖，包含 Module 变量和 UseCases 聚合
//   - 顶层（app/）提供类型别名和模块聚合，简化外部访问
//
// # 使用方式
//
//	import "github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/app"
//
//	// 引用 UseCase 类型（通过类型别名）
//	func NewHandler(useCases *app.SettingUseCases) *Handler { ... }
//
//	// 注册模块
//	fx.New(app.UseCaseModule)
package app

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/app/setting"
)

// SettingUseCases 是 setting.SettingUseCases 的类型别名。
//
// 提供便捷访问，避免外部代码直接导入子模块。
type SettingUseCases = setting.SettingUseCases

// UseCaseModule 提供 Settings 模块的 UseCase 层。
//
// 聚合所有子模块的 UseCase Module：
//   - setting.UseCaseModule: Setting 和 Category 的 UseCase Handlers
//
// 依赖：
//   - persistence.SettingRepositories（来自基础设施层）
//   - setting.SettingsCacheService（来自基础设施层）
var UseCaseModule = fx.Module("settings.usecase",
	setting.UseCaseModule,
)
