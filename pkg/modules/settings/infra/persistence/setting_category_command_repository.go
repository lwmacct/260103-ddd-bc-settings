package persistence

import (
	"context"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
	"gorm.io/gorm"
)

// settingCategoryCommandRepository 配置分类写仓储实现。
type settingCategoryCommandRepository struct {
	db *gorm.DB
}

// NewSettingCategoryCommandRepository 创建配置分类写仓储。
func NewSettingCategoryCommandRepository(db *gorm.DB) setting.SettingCategoryCommandRepository {
	return &settingCategoryCommandRepository{db: db}
}

// Create 创建配置分类。
func (r *settingCategoryCommandRepository) Create(ctx context.Context, category *setting.SettingCategory) error {
	model := newSettingCategoryModelFromEntity(category)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// 回写生成的 ID
	category.ID = model.ID
	return nil
}

// Update 更新配置分类。
func (r *settingCategoryCommandRepository) Update(ctx context.Context, category *setting.SettingCategory) error {
	model := newSettingCategoryModelFromEntity(category)

	// 仅更新可修改字段，Key 不可修改
	return r.db.WithContext(ctx).
		Model(&SettingCategoryModel{}).
		Where("id = ?", category.ID).
		Updates(map[string]any{
			"label":      model.Label,
			"icon":       model.Icon,
			"sort_order": model.Order,
		}).Error
}

// Delete 删除配置分类。
func (r *settingCategoryCommandRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&SettingCategoryModel{}, id).Error
}
