// Package adapters 提供 Settings 模块的外部适配器实现。
//
// 本包包含适配器层的所有实现：
//   - gin/：HTTP 协议适配器（Handler、Routes、Middleware）
//
// 适配器职责：
//   - 协议转换（HTTP → Application Command/Query）
//   - 请求验证（参数校验、格式转换）
//   - 响应格式化（统一的 JSON 响应）
//   - 路由定义（RESTful API 端点）
//
// # Architecture
//
// 适配器层依赖 Application 层的 UseCase，通过依赖注入组装。
//
// # Thread Safety
//
// 所有 Handler 都是并发安全的（只调用无状态 UseCase）。
package adapters
