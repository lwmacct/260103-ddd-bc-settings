// Package db 提供数据库管理命令。
package db

import (
	"github.com/urfave/cli/v3"
)

// Command 定义数据库管理命令。
var Command = &cli.Command{
	Name:  "db",
	Usage: "数据库操作",
	Description: `数据库迁移和种子数据管理。

示例:
  settings-server db migrate              执行迁移
  settings-server db migrate -c dev.yaml  使用指定配置
  settings-server db reset                重置数据库
  settings-server db seed                 执行种子数据`,
	Commands: []*cli.Command{
		{
			Name:    "migrate",
			Aliases: []string{"m"},
			Usage:   "执行数据库迁移",
			Action:  actionMigrate,
		},
		{
			Name:    "reset",
			Aliases: []string{"r"},
			Usage:   "重置数据库（删表+重建+种子数据）",
			Action:  actionReset,
		},
		{
			Name:    "seed",
			Aliases: []string{"s"},
			Usage:   "执行种子数据",
			Action:  actionSeed,
		},
	},
}
