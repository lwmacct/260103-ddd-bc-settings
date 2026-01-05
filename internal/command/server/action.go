package server

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/251219-go-pkg-logm/pkg/logm"
	"github.com/lwmacct/251219-go-pkg-logm/pkg/logm/formatter"
	"github.com/lwmacct/251219-go-pkg-logm/pkg/logm/writer"
	"github.com/lwmacct/260103-ddd-bc-settings/internal/config"
	"github.com/lwmacct/260103-ddd-bc-settings/internal/container"
	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings"
	"github.com/urfave/cli/v3"

	// Swagger docs - 空白导入触发 docs.go 的 init() 函数
	_ "github.com/lwmacct/260103-ddd-bc-settings/internal/command/server/docs"
)

// action 启动 HTTP 服务器。
func action(ctx context.Context, cmd *cli.Command) error {
	initLogger()
	cfg := LoadConfig(cmd)

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
	fxLogEnabled := FxLogEnable != nil && *FxLogEnable
	if !cfg.Server.FxLogEnabled && !fxLogEnabled {
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

// LoadConfig 加载配置。
func LoadConfig(cmd *cli.Command) *config.Config {
	// 从 CLI 获取 env-prefix 值
	var prefix string
	if EnvPrefix != nil && *EnvPrefix != "" {
		prefix = *EnvPrefix
	} else {
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
	var path string
	if ConfigFile != nil && *ConfigFile != "" {
		path = *ConfigFile
	} else {
		path = cmd.String("config")
	}
	if path != "" {
		opts = append(opts, cfgm.WithConfigPaths(path))
	}

	// 使用 cfgm.MustLoadCmd 加载配置
	return cfgm.MustLoadCmd(cmd, config.DefaultConfig(), "", opts...)
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
