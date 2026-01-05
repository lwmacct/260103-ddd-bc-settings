// Package main 提供完整的服务器入口，支持 CLI 命令。
//
// 命令：
//   - server      启动 HTTP 服务器（默认）
//   - db migrate  执行数据库迁移
//   - db reset    重置数据库（删表+重建+种子数据）
//   - db seed     执行种子数据
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	// Swagger docs - 空白导入触发 docs.go 的 init() 函数
	_ "github.com/lwmacct/260103-ddd-bc-settings/cmd/server/docs"

	// Swagger - 使用者完全控制文档生成
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/251219-go-pkg-logm/pkg/logm"
	"github.com/lwmacct/251219-go-pkg-logm/pkg/logm/formatter"
	"github.com/lwmacct/251219-go-pkg-logm/pkg/logm/writer"
	"github.com/lwmacct/260103-ddd-bc-settings/internal/config"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	// 启动器组装代码
	"github.com/lwmacct/260103-ddd-bc-settings/internal/container"

	// 业务模块 (Bounded Contexts)
	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings"
)

// Swagger 总体配置 - 使用者自定义
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

var (
	// 全局 flags（持久化到子命令）
	configFile  string
	envPrefix   string
	fxLogEnable bool
)

func main() {
	// 初始化日志 - 终端彩色 + 文件纯文本
	initLogger()

	app := &cli.Command{
		Name:    "server",
		Usage:   "Settings Bounded Context Server - 配置管理模块",
		Version: "1.0.0",
		Description: `基于 DDD + CQRS 架构的配置管理模块。

		示例:
		  server                  启动 HTTP 服务器
		  server db migrate       执行数据库迁移
		  server db reset         重置数据库
		  server -c config.yaml   指定配置文件
		  server --fx-log        启用 Fx 日志`,
		Commands: []*cli.Command{
			dbCommand(),
		},
		// Persistent flags 会自动传递给所有子命令
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
		EnableShellCompletion: true,
		Action:                startServer,
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}

// startServer 启动 HTTP 服务器。
func startServer(ctx context.Context, cmd *cli.Command) error {
	cfg := loadConfig(cmd)

	fxOptions := []fx.Option{
		fx.Supply(cfg),
		fx.StartTimeout(30 * time.Second),
		fx.StopTimeout(10 * time.Second),
		// Platform 层 (基础设施)
		container.InfraModule,
		// Settings 模块配置（从全局 config 提取）
		container.SettingsConfigModule,
		// 业务模块 (Bounded Contexts) - 完全自治
		settings.Module(),
		// HTTP 层 (跨模块handler + 路由)
		container.HTTPModule,
		// Swagger 端点注册 - 使用者决定
		fx.Invoke(func(r *gin.Engine) {
			r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}),
	}

	// CLI --fx-log 优先级高于配置文件
	if !cfg.Server.FxLogEnabled && !fxLogEnable {
		fxOptions = append(fxOptions, fx.WithLogger(func() fxevent.Logger {
			return nopLogger{}
		}))
	}

	fxApp := fx.New(fxOptions...)
	if err := fxApp.Err(); err != nil {
		return fmt.Errorf("create fx app: %w", err)
	}

	fxApp.Run()
	return nil
}

// dbCommand 数据库操作命令组。
// 使用 Persistent flags 使子命令继承父命令的配置选项。
func dbCommand() *cli.Command {
	cmd := &cli.Command{
		Name:  "db",
		Usage: "数据库操作",
		Description: `数据库迁移和种子数据管理。

		示例:
		  server db migrate       执行迁移
		  server db migrate -c dev.yaml  使用指定配置
		  server db reset         重置数据库
		  server db seed          执行种子数据`,
	}

	// 配置 flags，使用变量捕获以便子命令访问
	cmd.Flags = []cli.Flag{
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
	}

	// 添加子命令
	cmd.Commands = []*cli.Command{
		{
			Name:    "migrate",
			Aliases: []string{"m"},
			Usage:   "执行数据库迁移",
			Action:  migrateDatabase,
		},
		{
			Name:    "reset",
			Aliases: []string{"r"},
			Usage:   "重置数据库（删表+重建+种子数据）",
			Action:  resetDatabase,
		},
		{
			Name:    "seed",
			Aliases: []string{"s"},
			Usage:   "执行种子数据",
			Action:  seedDatabase,
		},
	}

	return cmd
}

// loadConfig 加载配置。
// 使用 cfgm.MustLoadCmd 从 CLI 命令和配置文件加载配置。
func loadConfig(cmd *cli.Command) *config.Config {
	// 从 CLI 获取 env-prefix 值
	prefix := envPrefix
	if prefix == "" {
		prefix = cmd.String("env-prefix")
		if prefix == "" {
			prefix = "APP_"
		}
	}

	// 构建 cfgm 选项
	opts := []cfgm.Option{
		cfgm.WithEnvPrefix(prefix),
	}

	// 如果指定了配置文件，添加该路径
	path := configFile
	if path == "" {
		path = cmd.String("config")
	}
	if path != "" {
		opts = append(opts, cfgm.WithConfigPaths(path))
	}

	// 使用 cfgm.MustLoadCmd 加载配置
	return cfgm.MustLoadCmd(cmd, config.DefaultConfig(), "", opts...)
}

// migrateDatabase 执行数据库迁移。
func migrateDatabase(ctx context.Context, cmd *cli.Command) error {
	cfg := loadConfig(cmd)

	appCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fxApp := fx.New(
		fx.Supply(cfg),
		fx.StartTimeout(5*time.Minute),
		fx.StopTimeout(10*time.Second),
		container.InfraModule,
		fx.Invoke(container.RunMigration),
		fx.WithLogger(func() fxevent.Logger { return nopLogger{} }),
	)

	if err := fxApp.Err(); err != nil {
		return fmt.Errorf("create fx app: %w", err)
	}

	// CLI 命令执行完成后立即退出，不等待信号
	if err := fxApp.Start(appCtx); err != nil {
		return err
	}
	return fxApp.Stop(appCtx)
}

// resetDatabase 重置数据库。
func resetDatabase(ctx context.Context, cmd *cli.Command) error {
	cfg := loadConfig(cmd)

	appCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fxApp := fx.New(
		fx.Supply(cfg),
		fx.StartTimeout(5*time.Minute),
		fx.StopTimeout(10*time.Second),
		container.InfraModule,
		fx.Invoke(container.RunReset),
		fx.WithLogger(func() fxevent.Logger { return nopLogger{} }),
	)

	if err := fxApp.Err(); err != nil {
		return fmt.Errorf("create fx app: %w", err)
	}

	// CLI 命令执行完成后立即退出，不等待信号
	if err := fxApp.Start(appCtx); err != nil {
		return err
	}
	return fxApp.Stop(appCtx)
}

// seedDatabase 执行种子数据。
func seedDatabase(ctx context.Context, cmd *cli.Command) error {
	cfg := loadConfig(cmd)

	appCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fxApp := fx.New(
		fx.Supply(cfg),
		fx.StartTimeout(5*time.Minute),
		fx.StopTimeout(10*time.Second),
		container.InfraModule,
		fx.Invoke(container.RunSeed),
		fx.WithLogger(func() fxevent.Logger { return nopLogger{} }),
	)

	if err := fxApp.Err(); err != nil {
		return fmt.Errorf("create fx app: %w", err)
	}

	// CLI 命令执行完成后立即退出，不等待信号
	if err := fxApp.Start(appCtx); err != nil {
		return err
	}
	return fxApp.Stop(appCtx)
}

// nopLogger 空日志记录器，不输出任何 Fx 框架日志。
type nopLogger struct{}

func (nopLogger) LogEvent(fxevent.Event) {}

// initLogger 初始化日志系统：终端彩色 + 文件纯文本。
func initLogger() {
	// 使用 logm 的 Handler（兼容 slog.Handler 接口）
	stdoutHandler := logm.NewHandler(&logm.HandlerConfig{
		LevelVar:   logm.GetLevelVar(),
		Formatter:  formatter.ColorText(),
		Writers:    []logm.Writer{writer.Stdout()},
		AddSource:  true,
		TimeFormat: "15:04:05.000",
	})

	// 文件 Handler
	fileHandler := logm.NewHandler(&logm.HandlerConfig{
		LevelVar:   logm.GetLevelVar(),
		Formatter:  formatter.Text(),
		Writers:    []logm.Writer{writer.File("/tmp/app.log")},
		AddSource:  true,
		TimeFormat: "2006-01-02 15:04:05.000",
	})

	// 使用自定义 multiHandler 组合两个 Handler
	logger := slog.New(&multiHandler{
		handlers: []slog.Handler{stdoutHandler, fileHandler},
	})
	slog.SetDefault(logger)
}

// multiHandler 将日志记录到多个 Handler。
type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		_ = h.Handle(ctx, r)
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}
