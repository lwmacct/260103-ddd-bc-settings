package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// settingCommandRepository 配置定义命令仓储的 GORM 实现
type settingCommandRepository struct {
	db *gorm.DB
}

// NewSettingCommandRepository 创建配置定义命令仓储实例
func NewSettingCommandRepository(db *gorm.DB) setting.CommandRepository {
	return &settingCommandRepository{db: db}
}

// Create 创建配置定义
func (r *settingCommandRepository) Create(ctx context.Context, s *setting.Setting) error {
	model := toSettingModel(s)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create setting definition: %w", err)
	}

	// 回写生成的 ID
	if entity := model.ToEntity(); entity != nil {
		*s = *entity
	}

	return nil
}

// Update 更新配置定义
func (r *settingCommandRepository) Update(ctx context.Context, s *setting.Setting) error {
	model := toSettingModel(s)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update setting definition: %w", err)
	}

	if entity := model.ToEntity(); entity != nil {
		*s = *entity
	}

	return nil
}

// Delete 删除配置定义
func (r *settingCommandRepository) Delete(ctx context.Context, key string) error {
	if err := r.db.WithContext(ctx).Where("key = ?", key).Delete(&SettingModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete setting definition: %w", err)
	}
	return nil
}

// BatchUpsert 批量插入或更新配置定义
func (r *settingCommandRepository) BatchUpsert(ctx context.Context, settings []*setting.Setting) error {
	if len(settings) == 0 {
		return nil
	}

	models := make([]*SettingModel, 0, len(settings))
	for _, s := range settings {
		if model := toSettingModel(s); model != nil {
			models = append(models, model)
		}
	}

	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"default_value", "category_id", "group", "value_type",
			"label", "ui_config", "order", "updated_at",
		}),
	}).Create(models).Error; err != nil {
		return fmt.Errorf("failed to batch upsert setting definitions: %w", err)
	}

	// 回写生成的 ID
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			*settings[i] = *entity
		}
	}

	return nil
}
