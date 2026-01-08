// Package app 提供 Settings 模块的应用层实现。
//
// 本包包含 Settings 业务逻辑的完整实现：
//   - setting/：配置管理用例（Command/Query Handlers）
//
// 应用层职责：
//   - 用例编排（Command/Query Handlers）
//   - DTO 定义（请求/响应数据结构）
//   - 业务流程控制（事务、缓存协调）
//   - 依赖注入（Fx 模块注册）
//
// # Architecture
//
// 采用分布式 Fx 模块架构：
//   - 每个子模块（如 setting）有独立的 module.go
//   - 顶层 app/module.go 提供类型别名和聚合
//
// # Thread Safety
//
// 所有 Handler 都是无状态的，并发安全。
package app
