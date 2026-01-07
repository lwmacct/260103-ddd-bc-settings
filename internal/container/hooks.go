package container

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"

	persistence "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/infra/persistence"
	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/infra/seeds"
	dbpkg "github.com/lwmacct/260103-ddd-shared/pkg/platform/db"
)

// GetAllModels 返回所有需要迁移的 Model。
func GetAllModels() []interface{} {
	return []interface{}{
		&persistence.SettingModel{},
		&persistence.SettingCategoryModel{},
	}
}

// RunMigration 执行数据库迁移（AutoMigrate）。
func RunMigration(lc fx.Lifecycle, db *gorm.DB) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("Running database migrations...")

			// 执行 AutoMigrate
			if err := db.AutoMigrate(GetAllModels()...); err != nil {
				return err
			}

			// 创建 Settings 索引
			if err := dbpkg.CreateIndexes(db, &persistence.SettingModel{}, []string{
				"idx_settings_category_sort",
				"idx_settings_visible_at",
				"idx_settings_configurable_at",
				"idx_settings_visible_configurable",
			}); err != nil {
				return err
			}

			slog.Info("Database migrations completed successfully")
			return nil
		},
	})
	return nil
}

// RunReset 重置数据库（删表+迁移+种子数据）。
func RunReset(lc fx.Lifecycle, db *gorm.DB, redis *redis.Client) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("Resetting db...")

			// 1. 清空 Redis 缓存
			if err := redis.FlushAll(ctx).Err(); err != nil {
				slog.Warn("Failed to flush Redis", "error", err)
			} else {
				slog.Info("Redis cache flushed")
			}

			// 2. 执行迁移
			if err := seeds.ExecuteSeeders(ctx, db); err != nil {
				return err
			}

			slog.Info("Database reset completed successfully")
			return nil
		},
	})
	return nil
}

// RunSeed 执行种子数据。
func RunSeed(lc fx.Lifecycle, db *gorm.DB) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("Running database seeders...")

			if err := seeds.ExecuteSeeders(ctx, db); err != nil {
				return err
			}

			slog.Info("Database seeding completed successfully")
			return nil
		},
	})
	return nil
}
