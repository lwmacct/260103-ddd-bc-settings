// Package docs 提供 Swagger 文档。
//
// 运行 swag init -g cmd/server/main.go 生成完整文档
package docs

// @title           Settings Bounded Context API
// @version         1.0
// @description     Settings 模块管理系统配置定义和分类。
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
