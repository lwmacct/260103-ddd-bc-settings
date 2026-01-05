// Package server 提供 HTTP 服务器命令。
//
//	@title           Settings Bounded Context API
//	@version         1.0
//	@description     基于 DDD + CQRS 架构的配置管理模块
//	@host            localhost:8080
//	@BasePath        /
//
//	@contact.name    API Support
//	@contact.url     https://github.com/lwmacct/260103-ddd-bc-settings
//
//	@license.name    MIT
//	@license.url     https://opensource.org/licenses/MIT
//
//	@securityDefinitions.apikey	BearerAuth
//	@in								header
//	@name							Authorization
//	@description					Bearer token authentication
package server

import (
	"github.com/urfave/cli/v3"
)

var (
	// 全局 flags 指针（由 main.go 设置）
	ConfigFile  *string
	EnvPrefix   *string
	FxLogEnable *bool
)

// Command 定义服务器命令。
var Command = &cli.Command{
	Name:    "server",
	Aliases: []string{"s"},
	Usage:   "启动 HTTP 服务器",
	Description: `启动 HTTP API 服务器。

示例:
  settings-server server              启动服务器
  settings-server server -c config.yaml  指定配置文件
  settings-server server --fx-log    启用 Fx 日志`,
	EnableShellCompletion: true,
	Action:                action,
}
