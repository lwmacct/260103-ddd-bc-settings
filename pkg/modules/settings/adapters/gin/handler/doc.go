// Package handler 实现 Settings 模块的 HTTP 处理器。
//
// 本包负责将 HTTP 请求绑定到 Application 层的 UseCase，
// 并将响应转换为统一的 JSON 格式。
//
// # Overview
//
// 本包包含以下处理器：
//   - [SettingHandler]: 系统配置和配置分类的 CRUD 操作
//
// # Thread Safety
//
// 所有 Handler 方法都是并发安全的，因为它们只调用无状态的 UseCase。
//
// # 依赖关系
//
// 本包依赖 Application 层的 UseCase（位于 `app/setting` 包），
// 并使用 `github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/response` 包
// 提供统一的响应格式。
package handler
