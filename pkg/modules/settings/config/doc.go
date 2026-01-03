// Package config 提供 Settings 模块的配置定义。
//
// # Settings 模块配置独立管理
//
// 本包定义 Settings 模块所需的所有配置项，实现模块配置自治：
//   - Redis：缓存键前缀配置
//
// 配置通过依赖注入由 internal/container 提供，Settings 模块不直接依赖 internal/config。
package config
