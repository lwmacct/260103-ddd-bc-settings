// Package settings 提供配置管理的完整模块实现。
//
// 本模块遵循 DDD 四层架构 + CQRS 模式：
//   - domain/setting：领域层（Setting 实体、Repository 接口、领域错误）
//   - infra：基础设施层（持久化、缓存、种子数据）
//   - app/setting：应用层（Command/Query Handlers、DTO）
//   - adapters/gin：适配器层（HTTP Handler、路由）
//
// 核心功能：
//   - 配置定义管理（CRUD + 批量操作）
//   - 配置分类管理（分组组织）
//   - 缓存失效策略（按 category 和 userID）
//   - 种子数据预置（系统初始化）
//
// 配置作用域：
//   - system：系统级配置（所有用户共享）
//   - user：用户级配置（可自定义覆盖）
//
// # Module Registration
//
// 本模块通过 fx.Module 注册所有子模块：
//
//	import "github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings"
//
//	fx.New(
//	    settings.Module(),  // Settings 完整模块
//	)
//
// # Thread Safety
//
// 所有导出的类型和函数都是并发安全的。
package settings
