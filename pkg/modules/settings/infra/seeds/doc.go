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
