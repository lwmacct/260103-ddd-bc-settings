package setting

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ScopeLevel 定义作用域层级类型
type ScopeLevel string

const (
	ScopeLevelSystem ScopeLevel = "system"
	ScopeLevelOrg    ScopeLevel = "org"
	ScopeLevelTeam   ScopeLevel = "team"
	ScopeLevelUser   ScopeLevel = "user"
	ScopeLevelPublic ScopeLevel = "public"
)

// scopeLevelOrder 定义层级顺序（权限从高到低）
// 数值越小，权限越高；数值越大，可见范围越广
var scopeLevelOrder = map[ScopeLevel]int{
	ScopeLevelSystem: 0, // 最高权限（仅管理员）
	ScopeLevelOrg:    1,
	ScopeLevelTeam:   2,
	ScopeLevelUser:   3, // 所有登录用户
	ScopeLevelPublic: 4, // 最广可见（包括未登录）
}

// compareScopeLevel 比较两个作用域的层级
// 返回值: -1 表示 a < b, 0 表示 a == b, 1 表示 a > b
func compareScopeLevel(a, b ScopeLevel) int {
	orderA, okA := scopeLevelOrder[a]
	orderB, okB := scopeLevelOrder[b]
	if !okA || !okB {
		return 0
	}
	if orderA < orderB {
		return -1
	} else if orderA > orderB {
		return 1
	}
	return 0
}

// Setting 配置定义实体。
// 存储配置项的 Schema 和默认值，支持分类、类型标注和 UI 元数据。
//
// DefaultValue 字段（JSONB）直接存储原生 JSON 值：
//   - 字符串: "My Site"
//   - 数值: 30
//   - 布尔值: true
//   - JSON 对象/数组: {"key": "value"} 或 [1, 2, 3]
//
// VisibleAt 字段决定配置的可见权限级别（权限从高到低）：
//   - "system": 仅管理员可见（敏感配置：SMTP、OAuth）
//   - "org": 组织管理者及以上可见
//   - "team": 团队管理者及以上可见
//   - "user": 所有登录用户可见
//   - "public": 任何人可见（包括未登录）
//
// ConfigurableAt 字段决定配置的最大可配置级别（到此级别为止可配置）：
//   - "system": 只有系统管理员可配置
//   - "org": Org 及以下级别可配置
//   - "team": Team 及以下级别可配置
//   - "user": 所有人可配置
//
// 场景示例：
//   - VisibleAt="system", ConfigurableAt="system": 敏感配置，仅管理员可见可改
//   - VisibleAt="user", ConfigurableAt="system": 所有用户可见，仅管理员可改
//   - VisibleAt="public", ConfigurableAt="system": 公开可见（如站点名称），仅管理员可改
//   - VisibleAt="user", ConfigurableAt="user": 用户可见可改
//
// InputType 决定前端控件类型和后端自动校验规则（email/url/password 等）。
// Validation 存储自定义 JSON Logic 规则，用于业务级增强校验。
// UIConfig 存储前端展示配置：hint（提示）、options（下拉选项）、depends_on（依赖关系）。
type Setting struct {
	ID             uint   // 唯一标识
	Key            string // 配置键，唯一约束
	DefaultValue   any    // 默认值（JSONB 原生值）
	VisibleAt      string // 可见权限级别：system | org | team | user | public
	ConfigurableAt string // 最大可配置级别：system | org | team | user
	CategoryID     uint   // 外键关联 SettingCategory.ID
	Group          string // 分组显示标签：基本设置, 本地化 等（直接存 label，空字符串表示无分组）
	ValueType      string // 值类型：string, number, boolean, json（用于类型校验）
	Label          string // 显示标签
	Order          int    // 排序权重（小的在前）

	// UI 配置
	InputType  string // 控件类型：text, email, url, password, select 等（决定自动校验规则）
	Validation string // 自定义校验规则（JSON Logic 格式）
	UIConfig   string // 前端展示配置：hint、options、depends_on（JSONB 字符串）

	CreatedAt time.Time
	UpdatedAt time.Time
}

// =============================================================================
// 验证方法
// =============================================================================

// Validate 验证实体完整性。
//
// 检查：
//   - Key 非空且格式有效
//   - CategoryID 非零（由数据库外键保证引用完整性）
//   - ValueType 有效
//   - InputType 有效
//   - VisibleAt 有效
//   - ConfigurableAt 有效
//   - DefaultValue 与 ValueType 匹配
//   - DefaultValue 通过 InputType 格式校验
func (s *Setting) Validate() error {
	if s.Key == "" {
		return ErrInvalidValue
	}
	if _, err := NewSettingKey(s.Key); err != nil {
		return err
	}
	if s.CategoryID == 0 {
		return ErrCategoryNotFound
	}
	if !s.IsValidValueType() {
		return ErrInvalidValueType
	}
	if !s.IsValidInputType() {
		return ErrInvalidInputType
	}
	if !s.IsValidVisibleAt() {
		return ErrInvalidVisibleAt
	}
	if !s.IsValidConfigurableAt() {
		return ErrInvalidConfigurableAt
	}
	if err := s.ValidateValue(s.DefaultValue); err != nil {
		return err
	}
	if err := s.ValidateByInputType(s.DefaultValue); err != nil {
		return err
	}
	return nil
}

// ValidateValue 验证给定值是否符合配置定义的类型。
func (s *Setting) ValidateValue(value any) error {
	if value == nil {
		return nil // nil 值总是允许的
	}

	switch s.ValueType {
	case ValueTypeString:
		if _, ok := value.(string); !ok {
			return ErrInvalidValueType
		}
	case ValueTypeNumber:
		switch value.(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64, json.Number:
			// 有效的数值类型
		default:
			return ErrInvalidValueType
		}
	case ValueTypeBoolean:
		if _, ok := value.(bool); !ok {
			return ErrInvalidValueType
		}
	case ValueTypeJSON:
		// JSON 类型接受 map 或 slice
		switch value.(type) {
		case map[string]any, []any:
			// 有效的 JSON 复合类型
		default:
			return ErrInvalidValueType
		}
	default:
		return ErrInvalidValueType
	}
	return nil
}

// IsValidValueType 报告 ValueType 是否有效。
func (s *Setting) IsValidValueType() bool {
	switch s.ValueType {
	case ValueTypeString, ValueTypeNumber, ValueTypeBoolean, ValueTypeJSON:
		return true
	default:
		return false
	}
}

// IsValidVisibleAt 报告 VisibleAt 是否有效。
func (s *Setting) IsValidVisibleAt() bool {
	switch s.VisibleAt {
	case string(ScopeLevelSystem), string(ScopeLevelOrg), string(ScopeLevelTeam),
		string(ScopeLevelUser), string(ScopeLevelPublic):
		return true
	default:
		return false
	}
}

// IsValidConfigurableAt 报告 ConfigurableAt 是否有效。
// ConfigurableAt 不允许 public（public 用户无法配置）。
func (s *Setting) IsValidConfigurableAt() bool {
	switch s.ConfigurableAt {
	case string(ScopeLevelSystem), string(ScopeLevelOrg), string(ScopeLevelTeam), string(ScopeLevelUser):
		return true
	default:
		return false
	}
}

// =============================================================================
// 可见性和可配置性方法
// =============================================================================

// IsVisibleAtScope 报告设置在指定权限级别是否可见。
// 可见条件：查询级别的权限 <= VisibleAt 的权限级别
// （权限越高数值越小，所以用 <=）
//
// 例如：
//   - VisibleAt=system (0)：只有 system (0) 可见
//   - VisibleAt=user (3)：system/org/team/user 都可见
//   - VisibleAt=public (4)：所有人可见（包括未登录）
func (s *Setting) IsVisibleAtScope(scope ScopeLevel) bool {
	return compareScopeLevel(scope, ScopeLevel(s.VisibleAt)) <= 0
}

// IsConfigurableAtScope 报告设置在指定级别是否可配置。
// 可配置条件：查询级别的层级 <= ConfigurableAt 的层级
//
// 例如：
//   - ConfigurableAt=system (0)：只有 system (0) 可配置
//   - ConfigurableAt=team (2)：system (0)、org (1)、team (2) 都可配置
func (s *Setting) IsConfigurableAtScope(scope ScopeLevel) bool {
	return compareScopeLevel(scope, ScopeLevel(s.ConfigurableAt)) <= 0
}

// IsSystemLevel 报告是否为系统级别设置（仅管理员可见）。
func (s *Setting) IsSystemLevel() bool {
	return s.VisibleAt == string(ScopeLevelSystem)
}

// IsOrgLevel 报告是否为组织级别设置（组织管理者及以上可见）。
func (s *Setting) IsOrgLevel() bool {
	return s.VisibleAt == string(ScopeLevelOrg)
}

// IsTeamLevel 报告是否为团队级别设置（团队管理者及以上可见）。
func (s *Setting) IsTeamLevel() bool {
	return s.VisibleAt == string(ScopeLevelTeam)
}

// IsUserLevel 报告是否为用户级别设置（所有登录用户可见）。
func (s *Setting) IsUserLevel() bool {
	return s.VisibleAt == string(ScopeLevelUser)
}

// IsPublicLevel 报告是否为公开级别设置（任何人可见，包括未登录）。
func (s *Setting) IsPublicLevel() bool {
	return s.VisibleAt == string(ScopeLevelPublic)
}

// IsPublic 报告是否为公开设置。
// 与 IsPublicLevel 相同，提供更直观的方法名。
func (s *Setting) IsPublic() bool {
	return s.IsPublicLevel()
}

// IsVisibleToUser 报告普通用户是否可见此配置。
func (s *Setting) IsVisibleToUser() bool {
	return s.IsVisibleAtScope(ScopeLevelUser)
}

// CanOrgConfigure 报告组织是否可以配置此设置。
func (s *Setting) CanOrgConfigure() bool {
	return s.IsConfigurableAtScope(ScopeLevelOrg)
}

// CanTeamConfigure 报告团队是否可以配置此设置。
func (s *Setting) CanTeamConfigure() bool {
	return s.IsConfigurableAtScope(ScopeLevelTeam)
}

// CanUserConfigure 报告用户是否可以配置此设置。
func (s *Setting) CanUserConfigure() bool {
	return s.IsConfigurableAtScope(ScopeLevelUser)
}

// IsOrgOnly 报告是否为组织专用设置（Org 可见可配置，Team 不可配置）。
func (s *Setting) IsOrgOnly() bool {
	return s.VisibleAt == string(ScopeLevelOrg) && s.ConfigurableAt == string(ScopeLevelOrg)
}

// IsTeamDefaultForUser 报告是否为团队可为用户设置的默认值（User 可见，Team 可配置）。
func (s *Setting) IsTeamDefaultForUser() bool {
	return s.VisibleAt == string(ScopeLevelUser) && s.ConfigurableAt == string(ScopeLevelTeam)
}

// =============================================================================
// 向后兼容方法（废弃）
// =============================================================================

// Deprecated: 使用 IsVisibleAtScope 和 IsConfigurableAtScope 代替

// IsSystemScope 报告是否为系统级配置。
//
// Deprecated: 使用 IsSystemLevel() 代替
func (s *Setting) IsSystemScope() bool {
	return s.IsSystemLevel()
}

// IsUserScope 报告是否为用户级配置。
//
// Deprecated: 使用 IsUserLevel() 代替
func (s *Setting) IsUserScope() bool {
	return s.IsUserLevel()
}

// IsOrgScope 报告是否为组织级配置。
//
// Deprecated: 使用 IsOrgLevel() 代替
func (s *Setting) IsOrgScope() bool {
	return s.IsOrgLevel()
}

// IsTeamScope 报告是否为团队级配置。
//
// Deprecated: 使用 IsTeamLevel() 代替
func (s *Setting) IsTeamScope() bool {
	return s.IsTeamLevel()
}

// =============================================================================
// 查询方法
// =============================================================================

// BelongsToCategoryID 报告是否属于指定分类 ID。
func (s *Setting) BelongsToCategoryID(categoryID uint) bool {
	return s.CategoryID == categoryID
}

// HasValidationRule 报告是否配置了自定义验证规则。
//
// 检查 Validation 字段是否非空。
func (s *Setting) HasValidationRule() bool {
	return s.Validation != ""
}

// IsRequired 报告是否为必填配置。
//
// 通过检查 Validation 中的 required 字段判断。
func (s *Setting) IsRequired() bool {
	if s.Validation == "" {
		return false
	}
	// 简单检查是否包含 required: true
	return strings.Contains(s.Validation, `"required":true`) ||
		strings.Contains(s.Validation, `"required": true`)
}

// GetKeyCategory 从 Key 提取 category 部分。
//
// 例如 "general.site_name" 返回 "general"。
func (s *Setting) GetKeyCategory() string {
	key, err := NewSettingKey(s.Key)
	if err != nil {
		return ""
	}
	return key.Category()
}

// GetKeyName 从 Key 提取 name 部分。
//
// 例如 "general.site_name" 返回 "site_name"。
func (s *Setting) GetKeyName() string {
	key, err := NewSettingKey(s.Key)
	if err != nil {
		return ""
	}
	return key.Name()
}

// =============================================================================
// 值处理方法
// =============================================================================

// CoerceValue 尝试将任意值转换为正确类型。
//
// 支持的转换：
//   - string -> number（解析为 float64）
//   - string -> boolean（"true"/"false"）
//   - number -> string（格式化）
func (s *Setting) CoerceValue(raw any) (any, error) {
	if raw == nil {
		return nil, nil //nolint:nilnil // nil input is valid and returns nil output
	}

	// 如果类型已匹配，直接返回
	if err := s.ValidateValue(raw); err == nil {
		return raw, nil
	}

	// 尝试类型转换
	switch s.ValueType {
	case ValueTypeString:
		return coerceToString(raw)
	case ValueTypeNumber:
		return coerceToNumber(raw)
	case ValueTypeBoolean:
		return coerceToBool(raw)
	case ValueTypeJSON:
		return coerceToJSON(raw)
	default:
		return nil, ErrInvalidValueType
	}
}

// GetDefaultValue 返回默认值。
func (s *Setting) GetDefaultValue() any {
	return s.DefaultValue
}

// =============================================================================
// 状态变更方法
// =============================================================================

// UpdateDefault 更新默认值。
//
// 验证新值是否与 ValueType 匹配，不匹配则返回错误。
func (s *Setting) UpdateDefault(value any) error {
	if err := s.ValidateValue(value); err != nil {
		return err
	}
	s.DefaultValue = value
	return nil
}

// UpdateLabel 更新显示标签。
func (s *Setting) UpdateLabel(label string) {
	s.Label = label
}

// UpdateOrder 更新排序权重。
func (s *Setting) UpdateOrder(order int) {
	s.Order = order
}

// =============================================================================
// 辅助函数
// =============================================================================

func coerceToString(v any) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val), nil
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val), nil
	case float32, float64:
		return fmt.Sprintf("%v", val), nil
	case bool:
		return strconv.FormatBool(val), nil
	default:
		return "", ErrInvalidValue
	}
}

func coerceToNumber(v any) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	case json.Number:
		return val.Float64()
	default:
		return 0, ErrInvalidValue
	}
}

func coerceToBool(v any) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case string:
		return strconv.ParseBool(val)
	case int, int64:
		return val != 0, nil
	case float64:
		return val != 0, nil
	default:
		return false, ErrInvalidValue
	}
}

func coerceToJSON(v any) (any, error) {
	switch val := v.(type) {
	case map[string]any:
		return val, nil
	case []any:
		return val, nil
	case string:
		// 尝试解析 JSON 字符串
		var result any
		if err := json.Unmarshal([]byte(val), &result); err != nil {
			return nil, ErrInvalidValue
		}
		return result, nil
	default:
		return nil, ErrInvalidValue
	}
}
