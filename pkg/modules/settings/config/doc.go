// Package config 提供 Settings 模块的配置定义。
//
// # Settings 模块配置契约
//
// 本包定义 Settings 模块所需的业务配置结构，实现模块配置自治：
//   - RedisKeyPrefix：Redis 缓存键前缀（含模块标识）
//   - BasePath：Settings API 的基础路径
//
// # 配置通过依赖注入提供
//
// Settings 模块不直接依赖 internal/config，配置由容器层转换提供：
//
//	internal/config → internal/container → pkg/modules/settings/config
//
// # 默认值管理
//
// 所有默认值统一在 internal/config.DefaultConfig() 中定义：
//   - 生产环境：通过配置文件或环境变量提供
//   - 开发环境：使用 internal/config 中的默认值
//   - 测试环境：可在测试代码中直接构造 Config
//
// # 配置分层原则
//
// 本包定义的配置是模块级的业务配置：
//   - ✅ 包含：缓存键前缀、HTTP 路径等模块特定参数
//   - ❌ 不包含：Redis 连接、DB 连接等 Platform 层基础设施
//
// Platform 层基础设施（如 *redis.Client）由容器层直接注入。
package config
