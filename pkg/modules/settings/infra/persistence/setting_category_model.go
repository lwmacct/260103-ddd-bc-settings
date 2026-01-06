package persistence

import (
	"time"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// SettingCategoryModel 配置分类的 GORM 实体
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type SettingCategoryModel struct {
	ID    uint   `gorm:"primaryKey"`
	Key   string `gorm:"uniqueIndex;size:50;not null"`
	Label string `gorm:"size:200;not null"`
	Icon  string `gorm:"size:100;not null;default:'mdi-cog'"`
	Order int    `gorm:"column:sort_order;default:0;index"`
	Scope string `gorm:"size:20;not null;default:'system';index"` // 可见性级别

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName 指定配置分类表名
func (SettingCategoryModel) TableName() string {
	return "setting_categories"
}

// toCategoryModel 将 Domain Entity 转换为 GORM Model
func toCategoryModel(entity *setting.SettingCategory) *SettingCategoryModel {
	if entity == nil {
		return nil
	}

	return &SettingCategoryModel{
		ID:        entity.ID,
		Key:       entity.Key,
		Label:     entity.Label,
		Icon:      entity.Icon,
		Order:     entity.Order,
		Scope:     string(entity.Scope),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// ToEntity 将 GORM Model 转换为 Domain Entity
func (m *SettingCategoryModel) ToEntity() *setting.SettingCategory {
	if m == nil {
		return nil
	}

	return &setting.SettingCategory{
		ID:        m.ID,
		Key:       m.Key,
		Label:     m.Label,
		Icon:      m.Icon,
		Order:     m.Order,
		Scope:     setting.ScopeLevel(m.Scope),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// toCategoryEntities 将 GORM Model 切片转换为 Domain Entity 切片
func toCategoryEntities(models []SettingCategoryModel) []*setting.SettingCategory {
	if len(models) == 0 {
		return nil
	}

	entities := make([]*setting.SettingCategory, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			entities = append(entities, entity)
		}
	}

	return entities
}
