package setting

import (
	"strings"
)

// SettingKey 配置键值对象。
//
// 配置键格式为 "category.name"，例如 "general.site_name"。
// 通过值对象封装确保 key 格式的一致性和验证。
type SettingKey struct {
	value string
}

// NewSettingKey 创建配置键，验证格式。
//
// 有效格式要求：
//   - 非空字符串
//   - 包含至少一个点号分隔符
//   - 点号前后都有内容
func NewSettingKey(key string) (SettingKey, error) {
	if key == "" {
		return SettingKey{}, ErrInvalidKeyFormat
	}
	idx := strings.Index(key, ".")
	if idx <= 0 || idx >= len(key)-1 {
		return SettingKey{}, ErrInvalidKeyFormat
	}
	return SettingKey{value: key}, nil
}

// MustSettingKey 创建配置键，格式无效时 panic。
//
// 仅用于常量初始化或确定有效的场景。
func MustSettingKey(key string) SettingKey {
	sk, err := NewSettingKey(key)
	if err != nil {
		panic("invalid setting key: " + key)
	}
	return sk
}

// String 返回配置键的字符串表示。
func (k SettingKey) String() string {
	return k.value
}

// Category 提取配置键的分类部分。
//
// 例如 "general.site_name" 返回 "general"。
func (k SettingKey) Category() string {
	idx := strings.Index(k.value, ".")
	if idx <= 0 {
		return ""
	}
	return k.value[:idx]
}

// Name 提取配置键的名称部分。
//
// 例如 "general.site_name" 返回 "site_name"。
func (k SettingKey) Name() string {
	idx := strings.Index(k.value, ".")
	if idx < 0 || idx >= len(k.value)-1 {
		return ""
	}
	return k.value[idx+1:]
}

// Equal 报告两个配置键是否相等。
func (k SettingKey) Equal(other SettingKey) bool {
	return k.value == other.value
}

// IsEmpty 报告配置键是否为空。
func (k SettingKey) IsEmpty() bool {
	return k.value == ""
}

// Category 配置分类值对象。
//
// 用于对配置项进行逻辑分组，便于管理界面展示和权限控制。
// 实际的分类常量定义在 entity_setting.go 中（如 CategoryGeneral）。
type Category struct {
	value string
}

// NewCategory 创建配置分类，验证有效性。
//
// 使用 entity_setting.go 中定义的分类常量：
//   - CategoryGeneral
//   - CategorySecurity
//   - CategoryNotification
//   - CategoryBackup
func NewCategory(cat string) (Category, error) {
	switch cat {
	case "general", "security", "notification", "backup":
		return Category{value: cat}, nil
	default:
		return Category{}, ErrCategoryNotFound
	}
}

// MustCategory 创建配置分类，无效时 panic。
//
// 仅用于常量初始化或确定有效的场景。
func MustCategory(cat string) Category {
	c, err := NewCategory(cat)
	if err != nil {
		panic("invalid category: " + cat)
	}
	return c
}

// String 返回分类的字符串表示。
func (c Category) String() string {
	return c.value
}

// IsValid 报告分类是否有效。
func (c Category) IsValid() bool {
	switch c.value {
	case "general", "security", "notification", "backup":
		return true
	default:
		return false
	}
}

// Equal 报告两个分类是否相等。
func (c Category) Equal(other Category) bool {
	return c.value == other.value
}

// IsEmpty 报告分类是否为空。
func (c Category) IsEmpty() bool {
	return c.value == ""
}

// AllCategoryStrings 返回所有有效分类的字符串表示。
func AllCategoryStrings() []string {
	return []string{"general", "security", "notification", "backup"}
}
