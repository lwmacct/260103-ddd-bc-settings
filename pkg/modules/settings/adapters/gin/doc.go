// Package gin 实现 Settings 模块的 HTTP 适配器。
//
// 本包包含基于 Gin 框架的完整 HTTP 实现：
//   - handler/：HTTP 请求处理器
//   - routes/：RESTful API 路由定义
//   - middleware/：自定义中间件
//   - manualtest/：集成测试工具
//
// # API Endpoints
//
// Settings API：
//   - GET    /api/admin/settings           - 配置列表（层级）
//   - GET    /api/admin/settings/{key}     - 配置详情
//   - POST   /api/admin/settings           - 创建配置
//   - PUT    /api/admin/settings/{key}     - 更新配置
//   - DELETE /api/admin/settings/{key}     - 删除配置
//   - POST   /api/admin/settings/batch     - 批量更新
//
// Categories API：
//   - GET    /api/admin/settings/categories       - 分类列表
//   - GET    /api/admin/settings/categories/{id}  - 分类详情
//   - POST   /api/admin/settings/categories       - 创建分类
//   - PUT    /api/admin/settings/categories/{id}  - 更新分类
//   - DELETE /api/admin/settings/categories/{id}  - 删除分类
//
// # Thread Safety
//
// 所有 Handler 和 Middleware 都是并发安全的。
package gin
