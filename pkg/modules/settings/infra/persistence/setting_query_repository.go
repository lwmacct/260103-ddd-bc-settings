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
	return toSettingEntities(models), nil
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
	return toSettingEntities(models), nil
}

// FindByScope 根据作用域查找配置定义列表
//
// Deprecated: 使用 FindByVisibleAt 或 FindByConfigurableAt 代替
//
// 为向后兼容保留，使用 VisibleAt 作为过滤条件
func (r *settingQueryRepository) FindByScope(ctx context.Context, scope string) ([]*setting.Setting, error) {
	var models []SettingModel
	err := r.db.WithContext(ctx).
		Where("visible_at = ?", scope).
		Order(`category_id ASC, "group" ASC, "order" ASC, key ASC`).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find setting definitions by scope: %w", err)
	}
	return toSettingEntities(models), nil
}

// FindVisibleToUser 查找普通用户可见的配置定义
// 可见条件：VisibleAt <= user（即 user 级别及以上的设置）
func (r *settingQueryRepository) FindVisibleToUser(ctx context.Context) ([]*setting.Setting, error) {
	var models []SettingModel
	// 用户可见：visible_at 为 system, org, team, user 的所有设置
	err := r.db.WithContext(ctx).
		Where("visible_at IN ?", []string{
			string(setting.ScopeLevelSystem),
			string(setting.ScopeLevelOrg),
			string(setting.ScopeLevelTeam),
			string(setting.ScopeLevelUser),
		}).
		Order(`category_id ASC, "group" ASC, "order" ASC, key ASC`).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find settings visible to user: %w", err)
	}
	return toSettingEntities(models), nil
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
	return toSettingEntities(models), nil
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

// =============================================================================
// 新增查询方法
// =============================================================================

// FindByVisibleAt 查询对指定级别可见的设置
// 返回条件：查询级别的层级 >= visible_at 的层级
func (r *settingQueryRepository) FindByVisibleAt(ctx context.Context, visibleAt setting.ScopeLevel) ([]*setting.Setting, error) {
	var models []SettingModel
	// 简化：查询 visible_at <= 查询级别的所有设置
	// PostgreSQL 的字符串比较可用于枚举值
	err := r.db.WithContext(ctx).
		Where("visible_at <= ?", visibleAt).
		Order(`category_id ASC, "group" ASC, "order" ASC, key ASC`).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find settings visible at %s: %w", visibleAt, err)
	}
	return toSettingEntities(models), nil
}

// FindByConfigurableAt 查询指定级别可配置的设置
// 返回条件：查询级别的层级 >= configurable_at 的层级
func (r *settingQueryRepository) FindByConfigurableAt(ctx context.Context, configurableAt setting.ScopeLevel) ([]*setting.Setting, error) {
	var models []SettingModel
	// 简化：查询 configurable_at <= 查询级别的所有设置
	err := r.db.WithContext(ctx).
		Where("configurable_at <= ?", configurableAt).
		Order(`category_id ASC, "group" ASC, "order" ASC, key ASC`).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find settings configurable at %s: %w", configurableAt, err)
	}
	return toSettingEntities(models), nil
}

// FindByVisibleAndConfigurable 查询同时满足可见性和可配置性的设置
func (r *settingQueryRepository) FindByVisibleAndConfigurable(
	ctx context.Context,
	visibleAt setting.ScopeLevel,
	configurableAt setting.ScopeLevel,
) ([]*setting.Setting, error) {
	var models []SettingModel
	err := r.db.WithContext(ctx).
		Where("visible_at <= ? AND configurable_at <= ?", visibleAt, configurableAt).
		Order(`category_id ASC, "group" ASC, "order" ASC, key ASC`).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find settings: %w", err)
	}
	return toSettingEntities(models), nil
}
