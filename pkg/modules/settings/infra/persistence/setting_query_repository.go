package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
	"gorm.io/gorm"
)

// settingQueryRepository 配置定义查询仓储的 GORM 实现
type settingQueryRepository struct {
	db *gorm.DB
}

// NewSettingQueryRepository 创建配置定义查询仓储实例
func NewSettingQueryRepository(db *gorm.DB) setting.QueryRepository {
	return &settingQueryRepository{db: db}
}

// FindByKey 根据 Key 查找配置定义
func (r *settingQueryRepository) FindByKey(ctx context.Context, key string) (*setting.Setting, error) {
	var model SettingModel
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // returns nil for not found, valid pattern
		}
		return nil, fmt.Errorf("failed to find setting definition by key: %w", err)
	}
	return model.ToEntity(), nil
}

// FindByKeys 根据多个 Key 批量查找配置定义
func (r *settingQueryRepository) FindByKeys(ctx context.Context, keys []string) ([]*setting.Setting, error) {
	if len(keys) == 0 {
		return []*setting.Setting{}, nil
	}
	var models []SettingModel
	err := r.db.WithContext(ctx).Where("key IN ?", keys).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find setting definitions by keys: %w", err)
	}
	return mapSettingModelsToEntities(models), nil
}

// FindByCategoryID 根据分类 ID 查找配置定义列表
func (r *settingQueryRepository) FindByCategoryID(ctx context.Context, categoryID uint) ([]*setting.Setting, error) {
	var models []SettingModel
	err := r.db.WithContext(ctx).
		Where("category_id = ?", categoryID).
		Order(`"group" ASC, "order" ASC, key ASC`).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find setting definitions by category ID: %w", err)
	}
	return mapSettingModelsToEntities(models), nil
}

// FindByScope 根据作用域查找配置定义列表
func (r *settingQueryRepository) FindByScope(ctx context.Context, scope string) ([]*setting.Setting, error) {
	var models []SettingModel
	err := r.db.WithContext(ctx).
		Where("scope = ?", scope).
		Order(`category_id ASC, "group" ASC, "order" ASC, key ASC`).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find setting definitions by scope: %w", err)
	}
	return mapSettingModelsToEntities(models), nil
}

// FindVisibleToUser 查找普通用户可见的配置定义
// 包含: scope=user（用户设置）+ scope=system 且 public=true（公开系统设置）
func (r *settingQueryRepository) FindVisibleToUser(ctx context.Context) ([]*setting.Setting, error) {
	var models []SettingModel
	err := r.db.WithContext(ctx).
		Where("scope = ? OR (scope = ? AND public = ?)", setting.ScopeUser, setting.ScopeSystem, true).
		Order(`category_id ASC, "group" ASC, "order" ASC, key ASC`).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find settings visible to user: %w", err)
	}
	return mapSettingModelsToEntities(models), nil
}

// FindAll 查找所有配置定义
func (r *settingQueryRepository) FindAll(ctx context.Context) ([]*setting.Setting, error) {
	var models []SettingModel
	err := r.db.WithContext(ctx).
		Order(`category_id ASC, "group" ASC, "order" ASC, key ASC`).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all setting definitions: %w", err)
	}
	return mapSettingModelsToEntities(models), nil
}

// ExistsByKey 检查 Key 是否已存在
func (r *settingQueryRepository) ExistsByKey(ctx context.Context, key string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&SettingModel{}).Where("key = ?", key).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check setting definition existence: %w", err)
	}
	return count > 0, nil
}
