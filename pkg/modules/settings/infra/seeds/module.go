package seeds

import (
	"context"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// SeedsModule 提供 Settings 模块的 Seeds 执行模块。
//
// 在应用启动时（OnStart 钩子）执行种子数据。
var SeedsModule = fx.Module("settings.seeds",
	fx.Invoke(func(db *gorm.DB) {
		// TODO: 在 OnStart 钩子中执行种子数据
		_ = db
	}),
)

// ExecuteSeeders 执行 Settings 模块的种子数据。
func ExecuteSeeders(ctx context.Context, db *gorm.DB) error {
	seeders := DefaultSeeders()
	for _, seeder := range seeders {
		if err := seeder.Seed(ctx, db); err != nil {
			return err
		}
	}
	return nil
}
