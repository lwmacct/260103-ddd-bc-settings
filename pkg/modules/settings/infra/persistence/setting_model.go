package persistence

import (
	"encoding/json"
	"time"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
	"gorm.io/datatypes"
)

// SettingModel 配置定义的 GORM 实体
//
// 索引设计：
//   - idx_settings_category_sort: 复合索引 (category_id, group, order, key) 覆盖分类查询和排序
//   - idx_settings_visible_at: 单列索引用于可见性过滤
//   - idx_settings_configurable_at: 单列索引用于可配置性过滤
//   - idx_settings_visible_configurable: 复合索引 (visible_at, configurable_at) 用于组合查询
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type SettingModel struct {
	ID             uint           `gorm:"primaryKey"`
	Key            string         `gorm:"uniqueIndex;size:100;not null"`
	DefaultValue   datatypes.JSON `gorm:"type:jsonb;not null;default:'null'"` // JSONB 原生值
	VisibleAt      string         `gorm:"size:20;not null;default:'user';index:idx_settings_visible_at;index:idx_settings_visible_configurable,priority:1"`
	ConfigurableAt string         `gorm:"size:20;not null;default:'user';index:idx_settings_configurable_at;index:idx_settings_visible_configurable,priority:2"`

	// 复合索引：覆盖 FindByCategoryID 的 WHERE + ORDER BY
	CategoryID uint   `gorm:"not null;index:idx_settings_category_sort,priority:1"`
	Group      string `gorm:"size:100;default:''"`                                   // 分组显示标签（直接存 label，空字符串表示无分组）
	Order      int    `gorm:"default:0;index:idx_settings_category_sort,priority:2"` // 排序（组间 + 组内）

	ValueType string `gorm:"size:20;default:'string'"`
	Label     string `gorm:"size:200"`

	// UI 配置
	InputType  string         `gorm:"column:input_type;size:32;not null;default:'text'"` // 控件类型
	Validation string         `gorm:"column:validation;type:text"`                       // JSON Logic 规则
	UIConfig   datatypes.JSON `gorm:"type:jsonb;default:'{}'"`                           // hint/options/depends_on

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName 指定配置定义表名
func (SettingModel) TableName() string {
	return "settings"
}

func newSettingModelFromEntity(entity *setting.Setting) *SettingModel {
	if entity == nil {
		return nil
	}

	// 将 any 类型的 DefaultValue 序列化为 JSON
	defaultValueJSON, _ := json.Marshal(entity.DefaultValue) //nolint:errchkjson // DefaultValue 是任意 JSONB 值

	return &SettingModel{
		ID:             entity.ID,
		Key:            entity.Key,
		DefaultValue:   datatypes.JSON(defaultValueJSON),
		VisibleAt:      entity.VisibleAt,
		ConfigurableAt: entity.ConfigurableAt,
		CategoryID:     entity.CategoryID,
		Group:          entity.Group,
		Order:          entity.Order,
		ValueType:      entity.ValueType,
		Label:          entity.Label,
		InputType:      entity.InputType,
		Validation:     entity.Validation,
		UIConfig:       datatypes.JSON(entity.UIConfig),
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}
}

// ToEntity 将 GORM Model 转换为 Domain Entity
func (m *SettingModel) ToEntity() *setting.Setting {
	if m == nil {
		return nil
	}

	// 将 JSON 反序列化为 any 类型
	var defaultValue any
	_ = json.Unmarshal(m.DefaultValue, &defaultValue)

	return &setting.Setting{
		ID:             m.ID,
		Key:            m.Key,
		DefaultValue:   defaultValue,
		VisibleAt:      m.VisibleAt,
		ConfigurableAt: m.ConfigurableAt,
		CategoryID:     m.CategoryID,
		Group:          m.Group,
		ValueType:      m.ValueType,
		Label:          m.Label,
		Order:          m.Order,
		InputType:      m.InputType,
		Validation:     m.Validation,
		UIConfig:       string(m.UIConfig),
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func mapSettingModelsToEntities(models []SettingModel) []*setting.Setting {
	if len(models) == 0 {
		return nil
	}

	defs := make([]*setting.Setting, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			defs = append(defs, entity)
		}
	}

	return defs
}
