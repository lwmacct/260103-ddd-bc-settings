package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/lwmacct/260103-ddd-bc-settings/internal/command/db"
	"github.com/lwmacct/260103-ddd-bc-settings/internal/command/server"
	"github.com/urfave/cli/v3"
)

var (
	// 全局 flags（持久化到子命令）
	configFile  string
	envPrefix   string
	fxLogEnable bool
)

func main() {
	app := &cli.Command{
		Name:    "settings-server",
		Version: "1.0.0",
		Usage:   "Settings Bounded Context Server - 配置管理模块",
		Description: `基于 DDD + CQRS 架构的配置管理模块。

示例:
  settings-server              启动 HTTP 服务器
  settings-server db migrate   执行数据库迁移
  settings-server db reset     重置数据库`,
		Commands: []*cli.Command{
			server.Command,
			db.Command,
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "配置文件路径",
				Destination: &configFile,
			},
			&cli.StringFlag{
				Name:        "env-prefix",
				Usage:       "环境变量前缀",
				Value:       "APP_",
				Destination: &envPrefix,
			},
			&cli.BoolFlag{
				Name:        "fx-log",
				Usage:       "启用 Fx 依赖注入日志",
				Destination: &fxLogEnable,
			},
		},
	}

	// 将全局 flags 传递给子命令
	server.ConfigFile = &configFile
	server.EnvPrefix = &envPrefix
	server.FxLogEnable = &fxLogEnable

	if err := app.Run(context.Background(), os.Args); err != nil {
		slog.Error("Application failed to run", "error", err)
		os.Exit(1)
	}
}
