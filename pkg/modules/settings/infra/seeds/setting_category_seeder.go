package seeds

import (
	"context"
	"log/slog"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/infra/persistence"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SettingCategorySeeder 配置分类种子数据
//
// Icon 字段使用 @mdi/font 图标库，格式为 "mdi-{icon-name}"
// 图标参考：https://pictogrammers.com/library/mdi/
type SettingCategorySeeder struct{}

// Seed 执行配置分类种子数据填充
func (s *SettingCategorySeeder) Seed(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	categories := []persistence.SettingCategoryModel{
		{
			Key:   "general",
			Label: "常规设置",
			Icon:  "mdi-cog",
			Order: 1,
			Scope: "user", // 用户可见
		},
		{
			Key:   "security",
			Label: "安全设置",
			Icon:  "mdi-shield-lock",
			Order: 2,
			Scope: "system", // 仅管理员
		},
		{
			Key:   "email",
			Label: "邮件服务",
			Icon:  "mdi-email-outline",
			Order: 3,
			Scope: "system", // 仅管理员
		},
		{
			Key:   "oauth",
			Label: "第三方登录",
			Icon:  "mdi-account-key",
			Order: 4,
			Scope: "system", // 仅管理员
		},
		{
			Key:   "notification",
			Label: "通知设置",
			Icon:  "mdi-bell",
			Order: 5,
			Scope: "user", // 用户可见
		},
		{
			Key:   "backup",
			Label: "备份设置",
			Icon:  "mdi-backup-restore",
			Order: 6,
			Scope: "system", // 仅管理员
		},
	}

	result := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"label", "icon", "sort_order", "scope",
		}), // 更新 UI 元数据和可见性，保持数据一致
	}).Create(&categories)
	if result.Error != nil {
		return result.Error
	}

	slog.Info("Seeded setting categories", "attempted", len(categories), "inserted", result.RowsAffected)
	return nil
}

// Name 返回种子器名称。
func (s *SettingCategorySeeder) Name() string {
	return "SettingCategorySeeder"
}
