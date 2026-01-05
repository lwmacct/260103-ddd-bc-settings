package persistence

import (
	"context"
	"errors"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
	"gorm.io/gorm"
)

// settingCategoryQueryRepository 配置分类查询仓储实现
type settingCategoryQueryRepository struct {
	db *gorm.DB
}

// NewSettingCategoryQueryRepository 创建配置分类查询仓储
func NewSettingCategoryQueryRepository(db *gorm.DB) setting.SettingCategoryQueryRepository {
	return &settingCategoryQueryRepository{db: db}
}

// FindByID 根据 ID 查询分类
func (r *settingCategoryQueryRepository) FindByID(ctx context.Context, id uint) (*setting.SettingCategory, error) {
	var model SettingCategoryModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // nil, nil 表示未找到记录
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

// FindByKey 根据分类键查询
func (r *settingCategoryQueryRepository) FindByKey(ctx context.Context, key string) (*setting.SettingCategory, error) {
	var model SettingCategoryModel
	if err := r.db.WithContext(ctx).Where("key = ?", key).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // nil, nil 表示未找到记录
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

// FindAll 查询所有分类，按 Order 升序排列
func (r *settingCategoryQueryRepository) FindAll(ctx context.Context) ([]*setting.SettingCategory, error) {
	var models []SettingCategoryModel
	if err := r.db.WithContext(ctx).Order("sort_order ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	return toCategoryEntities(models), nil
}

// FindByIDs 根据 ID 列表批量查询分类，按 Order 升序排列
func (r *settingCategoryQueryRepository) FindByIDs(ctx context.Context, ids []uint) ([]*setting.SettingCategory, error) {
	if len(ids) == 0 {
		return []*setting.SettingCategory{}, nil
	}
	var models []SettingCategoryModel
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Order("sort_order ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	return toCategoryEntities(models), nil
}

// ExistsByKey 检查指定 Key 是否已存在
func (r *settingCategoryQueryRepository) ExistsByKey(ctx context.Context, key string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&SettingCategoryModel{}).Where("key = ?", key).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
