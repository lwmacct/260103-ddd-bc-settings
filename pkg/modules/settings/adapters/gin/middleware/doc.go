// Package middleware 提供 Settings 模块的 HTTP 中间件。
//
// # Overview
//
// Settings 模块当前没有自定义中间件，使用全局中间件即可满足需求。
//
// 如果将来需要添加 Settings 专属的中间件（如配置访问权限控制、
// 配置版本校验等），可以在此包中实现。
//
// # 依赖关系
//
// 本包依赖 Gin 框架（github.com/gin-gonic/gin）。
package middleware
