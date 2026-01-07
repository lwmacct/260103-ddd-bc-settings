// Package seeds 提供 Settings 模块的种子数据实现。
//
// # Overview
//
// 本包实现了 Settings 模块的种子数据填充：
//   - [SettingCategorySeeder]: 配置分类种子（通用、安全、外观等）
//   - [SettingSeeder]: 配置定义种子（site_name, maintenance_mode 等）
//   - [DefaultSeeders]: 返回按依赖顺序排列的 Seeder 列表
//
// 所有 Seeder 实现幂等执行（使用 UPSERT），可重复运行。
//
// # Usage
//
//	import "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/infra/seeds"
//
//	// 执行所有 Seeder
//	seeders := seeds.DefaultSeeders()
//	for _, seeder := range seeders {
//		log.Printf("Running seeder: %s", seeder.Name())
//		if err := seeder.Seed(ctx, db); err != nil {
//			return fmt.Errorf("seeder %s failed: %w", seeder.Name(), err)
//		}
//	}
//
// # Thread Safety
//
// Seeder 实现是无状态的，可并发调用 Name() 方法。
// Seed() 方法依赖数据库事务保证一致性，不应并发调用同一个 Seeder。
package seeds

import (
	"context"

	"gorm.io/gorm"
)

// Seeder 种子数据接口。
//
// 所有 Seeder 必须实现此接口，以便统一执行。
type Seeder interface {
	// Seed 执行种子数据填充。
	Seed(ctx context.Context, db *gorm.DB) error

	// Name 返回 Seeder 名称。
	Name() string
}
