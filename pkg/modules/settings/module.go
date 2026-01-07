// Package settings 提供配置管理的完整模块实现。
//
// 本模块遵循 DDD 四层架构 + CQRS 模式：
//   - domain/setting：领域层（实体、Repository 接口）
//   - infra/persistence：基础设施层（GORM Repository 实现）
//   - infra/cache：缓存服务（Redis 实现）
//   - app/setting：应用层（UseCase Handlers、DTO）
//   - adapters/gin：适配器层（HTTP Handler、路由）
//
// 模块自治：所有业务相关代码内聚在本模块内，可独立编译和测试。
package settings

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/adapters/gin/handler"
	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app"
	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/infra/cache"
	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/infra/persistence"
)

// Module 提供 Settings Bounded Context 的完整 Fx 模块。
//
// 依赖顺序（严格按此顺序注册）：
//   1. cache.CacheModule - 缓存服务
//   2. persistence.RepositoryModule - 数据持久化（依赖缓存）
//   3. app.UseCaseModule - 用例处理器（依赖仓储）
//   4. handler.HandlerModule - HTTP 处理器（依赖用例）
//
// 使用示例：
//
//	fx.New(
//	    fx.Supply(cfg),
//	    container.InfraModule,      // DB, Redis
//	    settings.Module(),           // Settings 模块
//	    container.HTTPModule,        // HTTP 路由
//	)
func Module() fx.Option {
	return fx.Module("settings",
		// 基础设施层
		cache.CacheModule,
		persistence.RepositoryModule,

		// 应用层
		app.UseCaseModule,

		// 适配器层
		handler.HandlerModule,
	)
}
