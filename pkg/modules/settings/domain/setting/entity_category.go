package setting

import (
	"time"
)

// SettingCategory 配置分类实体。
//
// 存储 Category 的 UI 元数据，供前端直接渲染设置页面的 Tab。
// 与硬编码的常量（CategoryGeneral 等）不同，此实体的数据来自数据库，
// 允许运行时动态管理分类的显示属性。
//
// 设计说明：
//   - Key 字段与 [Setting.Category] 关联，用于分组查询
//   - Label/Icon/Order 供前端渲染 Tab 导航
//   - 前端无需硬编码任何 Category 信息
type SettingCategory struct {
	ID    uint   `json:"id"`    // 唯一标识
	Key   string `json:"key"`   // 分类键，唯一约束（general, security, notification, backup）
	Label string `json:"label"` // 显示名称（如 "常规设置"）
	Icon  string `json:"icon"`  // Tab 图标（mdi-xxx 格式）
	Order int    `json:"order"` // 排序权重（小的在前）

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// =============================================================================
// 验证方法
// =============================================================================

// Validate 验证实体完整性。
//
// 检查：
//   - Key 非空
//   - Label 非空
//   - Icon 非空
func (c *SettingCategory) Validate() error {
	if c.Key == "" {
		return ErrInvalidValue
	}
	if c.Label == "" {
		return ErrInvalidValue
	}
	if c.Icon == "" {
		return ErrInvalidValue
	}
	return nil
}

// IsValidKey 报告 Key 是否为已知的有效分类。
//
// 注意：此方法用于验证 Key 是否匹配预定义的分类常量。
// 如果需要支持动态分类，可以移除此验证或改为查询数据库。
func (c *SettingCategory) IsValidKey() bool {
	switch c.Key {
	case CategoryGeneral, CategorySecurity, CategoryNotification, CategoryBackup:
		return true
	default:
		return false
	}
}

// =============================================================================
// 查询方法
// =============================================================================

// MatchesCategory 报告是否与给定的 Setting.Category 匹配。
func (c *SettingCategory) MatchesCategory(category string) bool {
	return c.Key == category
}

// =============================================================================
// 状态变更方法
// =============================================================================

// UpdateLabel 更新显示名称。
func (c *SettingCategory) UpdateLabel(label string) {
	c.Label = label
}

// UpdateIcon 更新图标。
func (c *SettingCategory) UpdateIcon(icon string) {
	c.Icon = icon
}

// UpdateOrder 更新排序权重。
func (c *SettingCategory) UpdateOrder(order int) {
	c.Order = order
}
