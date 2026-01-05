package db

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/260103-ddd-bc-settings/internal/config"
	"github.com/lwmacct/260103-ddd-bc-settings/internal/container"
	"github.com/urfave/cli/v3"
)

// actionMigrate 执行数据库迁移。
func actionMigrate(ctx context.Context, cmd *cli.Command) error {
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

// actionReset 重置数据库。
func actionReset(ctx context.Context, cmd *cli.Command) error {
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

// actionSeed 执行种子数据。
func actionSeed(ctx context.Context, cmd *cli.Command) error {
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

// loadConfig 加载配置。
func loadConfig(cmd *cli.Command) *config.Config {
	// 从父命令获取 env-prefix 值
	prefix := cmd.String("env-prefix")
	if prefix == "" {
		prefix = "APP_"
	}

	// 构建 cfgm 选项
	opts := []cfgm.Option{
		cfgm.WithEnvPrefix(prefix),
	}

	// 如果指定了配置文件，添加该路径
	path := cmd.String("config")
	if path != "" {
		opts = append(opts, cfgm.WithConfigPaths(path))
	}

	// 使用 cfgm.MustLoadCmd 加载配置
	return cfgm.MustLoadCmd(cmd, config.DefaultConfig(), "", opts...)
}

// nopLogger 空日志记录器，不输出任何 Fx 框架日志。
type nopLogger struct{}

func (nopLogger) LogEvent(fxevent.Event) {}
