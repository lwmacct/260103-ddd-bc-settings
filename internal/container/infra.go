package container

import (
	"context"
	"log/slog"
	"net/url"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"

	dbpkg "github.com/lwmacct/260103-ddd-shared/pkg/platform/db"

	"github.com/lwmacct/260103-ddd-settings-bc/internal/config"
)

// InfraModule 提供基础设施模块（DB、Redis）。
var InfraModule = fx.Module("infra",
	fx.Provide(
		newDatabase,
		newRedisClient,
	),
)

func newDatabase(lc fx.Lifecycle, cfg *config.Config) (*gorm.DB, error) {
	ctx := context.Background()

	// 使用 DefaultConfig 创建数据库配置
	dbConfig := dbpkg.DefaultConfig(cfg.Data.PgsqlURL)

	db, err := dbpkg.NewConnection(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	// 数据库连接生命周期
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Ping 检查
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			return sqlDB.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			sqlDB, _ := db.DB()
			return sqlDB.Close()
		},
	})

	slog.Info("Database connected", "driver", "postgres")

	return db, nil
}

func newRedisClient(lc fx.Lifecycle, cfg *config.Config) (*redis.Client, error) {
	// 从 RedisURL 解析连接信息
	redisURL := cfg.Data.RedisURL

	// 解析 Redis URL (格式: redis://[:password@]host:port[/db])
	addr := "localhost:6379"
	var password string
	var db int = 0

	if redisURL != "" {
		u, err := url.Parse(redisURL)
		if err == nil {
			if u.Host != "" {
				addr = u.Host
			}
			if u.User != nil {
				password, _ = u.User.Password()
			}
			if len(u.Path) > 1 {
				dbStr := strings.TrimPrefix(u.Path, "/")
				dbInt, _ := strconv.Atoi(dbStr)
				db = dbInt
			}
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: 10,
	})

	// Redis 连接生命周期
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Ping(ctx).Err()
		},
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	slog.Info("Redis connected", "addr", addr, "db", db)

	return client, nil
}
