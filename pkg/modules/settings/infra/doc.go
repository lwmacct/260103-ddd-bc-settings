// Package infra 提供 Settings 模块的基础设施实现。
//
// 本包包含技术实现的完整支持：
//   - persistence/：数据持久化（GORM Repository）
//   - cache/：缓存服务（RedisJSON 实现）
//   - seeds/：种子数据（系统初始化）
//
// 基础设施职责：
//   - 数据持久化（PostgreSQL + GORM）
//   - 缓存管理（Redis + RedisJSON）
//   - 数据初始化（Seeder 执行）
//   - 依赖注入（Fx 模块注册）
//
// # Architecture
//
// 技术选型：
//   - GORM：ORM 框架
//   - RedisJSON：缓存存储
//   - Fx：依赖注入
//
// # Thread Safety
//
// 所有 Repository 和 Service 都并发安全。
package infra
